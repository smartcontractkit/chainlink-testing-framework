package chipingressset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v72/github"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type protoFile struct {
	Name    string
	Path    string
	Content string
}

type ProtoSchemaSet struct {
	URI           string   `toml:"uri"`
	Ref           string   `toml:"ref"`            // ref or tag or commit SHA
	Folders       []string `toml:"folders"`        // if not provided, all protos will be fetched, otherwise only protos in these folders will be fetched
	SubjectPrefix string   `toml:"subject_prefix"` // optional prefix for subjects
}

// SubjectNamingStrategyFn is a function that is used to determine the subject name for a given proto file in a given repo
type SubjectNamingStrategyFn func(subjectPrefix string, protoFile protoFile, repoConfig ProtoSchemaSet) (string, error)

// RepositoryToSubjectNamingStrategyFn is a map of repository names to SubjectNamingStrategyFn functions
type RepositoryToSubjectNamingStrategyFn map[string]SubjectNamingStrategyFn

func validateRepoConfiguration(repoConfig ProtoSchemaSet) error {
	if repoConfig.URI == "" {
		return errors.New("uri is required")
	}

	if !strings.HasPrefix(repoConfig.URI, "https://") && !strings.HasPrefix(repoConfig.URI, "file://") {
		return errors.New("uri has to start with either 'file://' or 'https://'")
	}

	if strings.HasPrefix(repoConfig.URI, "file://") {
		if repoConfig.Ref != "" {
			return errors.New("ref is not supported with local protos with 'file://' prefix")
		}

		return nil
	}

	trimmedURI := strings.TrimPrefix(repoConfig.URI, "https://")
	if !strings.HasPrefix(trimmedURI, "github.com") {
		return fmt.Errorf("only repositories hosted at github.com are supported, but %s was found", repoConfig.URI)
	}

	parts := strings.Split(trimmedURI, "/")
	if len(parts) < 3 {
		return fmt.Errorf("URI should have following format: 'https://github.com/<OWNER>/<REPOSITORY>', but %s was found", repoConfig.URI)
	}

	if repoConfig.Ref == "" {
		return errors.New("ref is required, when fetching protos from Github repository")
	}

	return nil
}

func DefaultRegisterAndFetchProtos(ctx context.Context, client *github.Client, protoSchemaSets []ProtoSchemaSet, schemaRegistryURL string) error {
	return RegisterAndFetchProtos(ctx, client, protoSchemaSets, schemaRegistryURL, map[string]SubjectNamingStrategyFn{})
}

func RegisterAndFetchProtos(ctx context.Context, client *github.Client, protoSchemaSets []ProtoSchemaSet, schemaRegistryURL string, repoToSubjectNamingStrategy RepositoryToSubjectNamingStrategyFn) error {
	framework.L.Debug().Msgf("Registering and fetching protos from %d repositories", len(protoSchemaSets))

	for _, protoSchemaSet := range protoSchemaSets {
		if valErr := validateRepoConfiguration(protoSchemaSet); valErr != nil {
			return errors.Wrapf(valErr, "invalid repo configuration for schema set: %v", protoSchemaSet)
		}
	}

	for _, protoSchemaSet := range protoSchemaSets {
		protos, protosErr := fetchProtoFilesInFolders(ctx, client, protoSchemaSet.URI, protoSchemaSet.Ref, protoSchemaSet.Folders)
		if protosErr != nil {
			return errors.Wrapf(protosErr, "failed to fetch protos from %s", protoSchemaSet.URI)
		}

		protoMap := make(map[string]string)
		subjects := make(map[string]string)

		for _, proto := range protos {
			protoMap[proto.Path] = proto.Content

			var subjectStrategy SubjectNamingStrategyFn
			if strategy, ok := repoToSubjectNamingStrategy[protoSchemaSet.URI]; ok {
				subjectStrategy = strategy
			} else {
				subjectStrategy = DefaultSubjectNamingStrategy
			}

			subjectMessage, nameErr := subjectStrategy(protoSchemaSet.SubjectPrefix, proto, protoSchemaSet)
			if nameErr != nil {
				return errors.Wrapf(nameErr, "failed to extract message name from %s", proto.Path)
			}
			subjects[proto.Path] = subjectMessage
		}

		registerErr := registerAllWithTopologicalSortingByTrial(schemaRegistryURL, protoMap, subjects)
		if registerErr != nil {
			return errors.Wrapf(registerErr, "failed to register protos from %s", protoSchemaSet.URI)
		}
	}

	return nil
}

func DefaultSubjectNamingStrategy(subjectPrefix string, proto protoFile, protoSchemaSet ProtoSchemaSet) (string, error) {
	packageName, packageErr := extractPackageNameWithRegex(proto.Content)
	if packageErr != nil {
		return "", errors.Wrapf(packageErr, "failed to extract package name from %s", proto.Path)
	}

	messageNames, nameErr := extractTopLevelMessageNamesWithRegex(proto.Content)
	if nameErr != nil {
		return "", errors.Wrapf(nameErr, "failed to extract message name from %s", proto.Path)
	}
	messageName := messageNames[0]

	return subjectPrefix + packageName + "." + messageName, nil
}

// extractPackageNameWithRegex extracts the package name from a proto source file using regex.
// It returns an error if no package name is found.
func extractPackageNameWithRegex(protoSrc string) (string, error) {
	matches := regexp.MustCompile(`(?m)^\s*package\s+([a-zA-Z0-9.]+)\s*;`).FindStringSubmatch(protoSrc)
	if len(matches) < 2 {
		return "", fmt.Errorf("no package name found in proto source")
	}

	if matches[1] == "" {
		return "", fmt.Errorf("empty package name found in proto source")
	}

	return matches[1], nil
}

// we use simple regex to extract top-level message names from a proto file
// so that we don't need to parse the proto file with a parser (which would require a lot of dependencies)
func extractTopLevelMessageNamesWithRegex(protoSrc string) ([]string, error) {
	matches := regexp.MustCompile(`(?m)^\s*message\s+(\w+)\s*{`).FindAllStringSubmatch(protoSrc, -1)
	var names []string
	for _, match := range matches {
		if len(match) >= 2 {
			names = append(names, match[1])
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no message names found in %s", protoSrc)
	}

	return names, nil
}

// Fetches .proto files from a GitHub repo optionally scoped to specific folders. It is recommended to use `*github.Client` with auth token to avoid rate limiting.
func fetchProtoFilesInFolders(ctx context.Context, client *github.Client, uri, ref string, folders []string) ([]protoFile, error) {
	framework.L.Debug().Msgf("Fetching proto files from %s in folders: %s", uri, strings.Join(folders, ", "))

	if strings.HasPrefix(uri, "file://") {
		return fetchProtosFromFilesystem(uri, folders)
	}

	parts := strings.Split(strings.TrimPrefix(uri, "https://"), "/")

	return fetchProtosFromGithub(ctx, client, parts[1], parts[2], ref, folders)
}

func fetchProtosFromGithub(ctx context.Context, client *github.Client, owner, repository, ref string, folders []string) ([]protoFile, error) {
	if client == nil {
		return nil, errors.New("github client cannot be nil")
	}

	var files []protoFile

	sha, shaErr := resolveRefSHA(ctx, client, owner, repository, ref)
	if shaErr != nil {
		return nil, errors.Wrapf(shaErr, "cannot resolve ref %q", ref)
	}

	tree, _, treeErr := client.Git.GetTree(ctx, owner, repository, sha, true)
	if treeErr != nil {
		return nil, errors.Wrap(treeErr, "failed to fetch tree")
	}

searchLoop:
	for _, entry := range tree.Entries {
		// skip non-blob entries and non-proto files
		if entry.GetType() != "blob" || entry.Path == nil || !strings.HasSuffix(*entry.Path, ".proto") {
			continue
		}

		// if folders are specified, check prefix match
		var folderFound string
		if len(folders) > 0 {
			matched := false
			for _, folder := range folders {
				if strings.HasPrefix(*entry.Path, strings.TrimSuffix(folder, "/")+"/") {
					matched = true
					folderFound = folder
					break
				}
			}
			if !matched {
				continue searchLoop
			}
		}

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repository, sha, *entry.Path)
		resp, respErr := http.Get(rawURL)
		if respErr != nil {
			return nil, errors.Wrapf(respErr, "failed tofetch %s", *entry.Path)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, errors.Errorf("bad status from GitHub for %s: %d", *entry.Path, resp.StatusCode)
		}

		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, errors.Wrapf(bodyErr, "failed to read body for %s", *entry.Path)
		}

		// subtract the folder from the path if it was provided, because if it is imported by some other protos
		// most probably it will be imported as a relative path, so we need to remove the folder from the path
		protoPath := *entry.Path
		if folderFound != "" {
			protoPath = strings.TrimPrefix(protoPath, strings.TrimSuffix(folderFound, "/")+"/")
		}

		files = append(files, protoFile{
			Name:    filepath.Base(*entry.Path),
			Path:    protoPath,
			Content: string(body),
		})
	}

	framework.L.Debug().Msgf("Fetched %d proto files from Github's %s/%s", len(files), owner, repository)

	if len(files) == 0 {
		return nil, fmt.Errorf("no proto files found in %s/%s in folders %s", owner, repository, strings.Join(folders, ", "))
	}

	return files, nil
}

func fetchProtosFromFilesystem(uri string, folders []string) ([]protoFile, error) {
	var files []protoFile

	protoDirPath := strings.TrimPrefix(uri, "file://")
	walkErr := filepath.Walk(protoDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var folderFound string
		if len(folders) > 0 {
			matched := false
			for _, folder := range folders {
				if strings.HasPrefix(strings.TrimPrefix(strings.TrimPrefix(path, protoDirPath), "/"), folder) {
					matched = true
					folderFound = folder
					break
				}
			}

			if !matched {
				return nil
			}
		}

		if !strings.HasSuffix(path, ".proto") {
			return nil
		}

		content, contentErr := os.ReadFile(path)
		if contentErr != nil {
			return errors.Wrapf(contentErr, "failed to read file at %s", path)
		}

		// subtract the folder from the path if it was provided, because if it is imported by some other protos
		// most probably it will be imported as a relative path, so we need to remove the folder from the path
		protoPath := strings.TrimPrefix(strings.TrimPrefix(path, protoDirPath), "/")
		if folderFound != "" {
			protoPath = strings.TrimPrefix(strings.TrimPrefix(protoPath, folderFound), strings.TrimSuffix(folderFound, "/"))
			protoPath = strings.TrimPrefix(protoPath, "/")
		}

		files = append(files, protoFile{
			Name:    filepath.Base(path),
			Path:    protoPath,
			Content: string(content),
		})

		return nil
	})
	if walkErr != nil {
		return nil, errors.Wrapf(walkErr, "failed to walk through directory %s", protoDirPath)
	}

	framework.L.Debug().Msgf("Fetched %d proto files from local %s", len(files), protoDirPath)

	if len(files) == 0 {
		return nil, fmt.Errorf("no proto files found in '%s' in folders %s", protoDirPath, strings.Join(folders, ", "))
	}

	return files, nil
}

func resolveRefSHA(ctx context.Context, client *github.Client, owner, repository, ref string) (string, error) {
	if refObj, _, err := client.Git.GetRef(ctx, owner, repository, "refs/tags/"+ref); err == nil {
		return refObj.GetObject().GetSHA(), nil
	}
	if refObj, _, err := client.Git.GetRef(ctx, owner, repository, "refs/heads/"+ref); err == nil {
		return refObj.GetObject().GetSHA(), nil
	}
	if commit, _, err := client.Repositories.GetCommit(ctx, owner, repository, ref, nil); err == nil {
		return commit.GetSHA(), nil
	}
	return "", fmt.Errorf("ref %q not found", ref)
}

type schemaStatus struct {
	Source     string
	Registered bool
	Version    int
}

// registerAllWithTopologicalSortingByTrial tries to register protos that have not been registered yet, and if it fails, it tries again with a different order
// it keeps doing this until all protos are registered or it fails to register any more protos
func registerAllWithTopologicalSortingByTrial(
	schemaRegistryURL string,
	protoMap map[string]string, // path -> proto source
	subjectMap map[string]string, // path -> subject
) error {
	framework.L.Info().Msgf("🔄 registering %d protobuf schemas", len(protoMap))
	schemas := map[string]*schemaStatus{}
	for path, src := range protoMap {
		schemas[path] = &schemaStatus{Source: src}
	}

	refs := []map[string]any{}

	for {
		progress := false
		failures := []string{}

		for path, schema := range schemas {
			if schema.Registered {
				continue
			}

			subject, ok := subjectMap[path]
			if !ok {
				failures = append(failures, fmt.Sprintf("%s: no subject found", path))
				continue
			}

			singleProtoFailures := []error{}
			framework.L.Debug().Msgf("🔄 registering %s as %s", path, subject)
			_, registerErr := registerSingleProto(schemaRegistryURL, subject, schema.Source, refs)
			if registerErr != nil {
				failures = append(failures, fmt.Sprintf("%s: %v", path, registerErr))
				singleProtoFailures = append(singleProtoFailures, registerErr)
				continue
			}

			schema.Registered = true
			schema.Version = 1
			refs = append(refs, map[string]any{
				"name":    path,
				"subject": subject,
				"version": 1,
			})

			framework.L.Info().Msgf("✔ registered: %s as %s", path, subject)

			progress = true
		}

		if !progress {
			if len(failures) > 0 {
				framework.L.Error().Msg("❌ Failed to register remaining schemas:")
				for _, msg := range failures {
					framework.L.Error().Msg("  " + msg)
				}
				return fmt.Errorf("unable to register %d schemas", len(failures))
			}
			break
		}
	}

	framework.L.Info().Msgf("✅ Successfully registered %d schemas", len(protoMap))
	return nil
}

func registerSingleProto(
	registryURL, subject, schemaSrc string,
	references []map[string]any,
) (int, error) {
	framework.L.Trace().Msgf("Registering schema %s", subject)

	body := map[string]any{
		"schemaType": "PROTOBUF",
		"schema":     schemaSrc,
	}
	if references != nil {
		body["references"] = references
	}

	payload, payloadErr := json.Marshal(body)
	if payloadErr != nil {
		return 0, errors.Wrap(payloadErr, "failed to marshal payload")
	}

	url := fmt.Sprintf("%s/subjects/%s/versions", registryURL, subject)

	resp, respErr := http.Post(url, "application/vnd.schemaregistry.v1+json", bytes.NewReader(payload))
	if respErr != nil {
		return 0, errors.Wrap(respErr, "failed to post to schema registry")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		data, dataErr := io.ReadAll(resp.Body)
		if dataErr != nil {
			return 0, errors.Wrap(dataErr, "failed to read response body")
		}
		return 0, fmt.Errorf("schema registry error (%d): %s", resp.StatusCode, data)
	}

	var result struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, errors.Wrap(err, "failed to decode response")
	}

	framework.L.Debug().Msgf("Registered schema %s with ID %d", subject, result.ID)

	return result.ID, nil
}
