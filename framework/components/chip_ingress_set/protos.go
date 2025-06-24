package chipingressset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type RepoConfiguration struct {
	Owner   string
	Repo    string
	Ref     string   // ref or tag or commit SHA
	Folders []string // if not provided, all protos will be fetched, otherwise only protos in these folders will be fetched
}

// SubjectNamingStrategyFn is a function that is used to determine the subject name for a given proto file in a given repo
type SubjectNamingStrategyFn func(path, source string, repoConfig RepoConfiguration) (string, error)

// RepoToSubjectNamingStrategyFn is a map of repo names to SubjectNamingStrategyFn functions
type RepoToSubjectNamingStrategyFn map[string]SubjectNamingStrategyFn

// DefaultRepoToSubjectNamingStrategy is a map of repo names to SubjectNamingStrategyFn functions
var DefaultRepoToSubjectNamingStrategy = RepoToSubjectNamingStrategyFn{
	"chainlink-protos": ChainlinkProtosSubjectNamingStrategy,
}

func DefaultRegisterAndFetchProtos(ctx context.Context, client *github.Client, reposConfig []RepoConfiguration, schemaRegistryURL string) error {
	return RegisterAndFetchProtos(ctx, client, reposConfig, schemaRegistryURL, DefaultRepoToSubjectNamingStrategy)
}

func RegisterAndFetchProtos(ctx context.Context, client *github.Client, reposConfig []RepoConfiguration, schemaRegistryURL string, repoToSubjectNamingStrategy RepoToSubjectNamingStrategyFn) error {
	for _, repoConfig := range reposConfig {
		protos, protosErr := fetchProtoFilesInFolders(ctx, client, repoConfig.Owner, repoConfig.Repo, repoConfig.Ref, repoConfig.Folders)
		if protosErr != nil {
			return errors.Wrapf(protosErr, "failed to fetch protos from %s/%s", repoConfig.Owner, repoConfig.Repo)
		}

		protoMap := make(map[string]string)
		subjectMap := make(map[string]string)

		for _, pf := range protos {
			protoMap[pf.Path] = pf.Content

			var subjectStrategy SubjectNamingStrategyFn
			if strategy, ok := repoToSubjectNamingStrategy[repoConfig.Repo]; ok {
				subjectStrategy = strategy
			} else {
				subjectStrategy = DefaultSubjectNamingStrategy
			}

			subject, nameErr := subjectStrategy(pf.Path, pf.Content, repoConfig)
			if nameErr != nil {
				return errors.Wrapf(nameErr, "failed to extract message name from %s", pf.Path)
			}
			subjectMap[pf.Path] = subject
		}

		registerErr := registerAllWithTopologicalSortingByTrial(schemaRegistryURL, protoMap, subjectMap)
		if registerErr != nil {
			return errors.Wrapf(registerErr, "failed to register protos from %s/%s", repoConfig.Owner, repoConfig.Repo)
		}
	}

	return nil
}

func DefaultSubjectNamingStrategy(path, source string, repoConfig RepoConfiguration) (string, error) {
	messageName, nameErr := extractTopLevelMessageNamesWithRegex(source)
	if nameErr != nil {
		return "", errors.Wrapf(nameErr, "failed to extract message name from %s", path)
	}
	return repoConfig.Repo + "." + messageName, nil
}

// TODO once we have single source of truth for the relationship between protos and subjects, we need to modify this function
func ChainlinkProtosSubjectNamingStrategy(path, source string, repoConfig RepoConfiguration) (string, error) {
	messageName, nameErr := extractTopLevelMessageNamesWithRegex(source)
	if nameErr != nil {
		return "", errors.Wrapf(nameErr, "failed to extract message name from %s", path)
	}

	// this only covers BaseMessage
	if strings.HasPrefix(path, "common") {
		return "cre-pb." + messageName, nil
	}

	// this covers all other protos we currently have in the chainlink-protos repo
	subject := "cre-workflows."
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) > 1 {
		for _, part := range pathSplit {
			matches := regexp.MustCompile(`v[0-9]+`).FindAllStringSubmatch(part, -1)
			if len(matches) > 0 {
				subject += matches[0][0]
			}
		}
	} else {
		return "", fmt.Errorf("no subject found for %s", path)
	}

	return subject + "." + messageName, nil
}

// we use simple regex to extract top-level message names from a proto file
// so that we don't need to parse the proto file with a parser (which would require a lot of dependencies)
func extractTopLevelMessageNamesWithRegex(protoSrc string) (string, error) {
	matches := regexp.MustCompile(`(?m)^\s*message\s+(\w+)\s*{`).FindAllStringSubmatch(protoSrc, -1)
	var names []string
	for _, match := range matches {
		if len(match) >= 2 {
			names = append(names, match[1])
		}
	}

	if len(names) == 0 {
		return "", fmt.Errorf("no message names found in %s", protoSrc)
	}

	// even though there could be more than 1 message in a single proto, we still need to register all of them under one subject
	return names[0], nil
}

// Fetches .proto files from a GitHub repo optionally scoped to specific folders. It is recommended to use `*github.Client` with auth token to avoid rate limiting.
func fetchProtoFilesInFolders(ctx context.Context, client *github.Client, owner, repo, ref string, folders []string) ([]protoFile, error) {
	var files []protoFile

	sha, shaErr := resolveRefSHA(ctx, client, owner, repo, ref)
	if shaErr != nil {
		return nil, errors.Wrapf(shaErr, "cannot resolve ref %q", ref)
	}

	tree, _, treeErr := client.Git.GetTree(ctx, owner, repo, sha, true)
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

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, sha, *entry.Path)
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

	return files, nil
}

func resolveRefSHA(ctx context.Context, client *github.Client, owner, repo, ref string) (string, error) {
	if refObj, _, err := client.Git.GetRef(ctx, owner, repo, "refs/tags/"+ref); err == nil {
		return refObj.GetObject().GetSHA(), nil
	}
	if refObj, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+ref); err == nil {
		return refObj.GetObject().GetSHA(), nil
	}
	if commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, ref, nil); err == nil {
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

			framework.L.Debug().Msgf("🔄 registering %s as %s\n", path, subject)
			_, registerErr := registerSingleProto(schemaRegistryURL, subject, schema.Source, refs)
			if registerErr != nil {
				failures = append(failures, fmt.Sprintf("%s: %v", path, registerErr))
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

	framework.L.Info().Msg("✅ All schemas successfully registered.")
	return nil
}

func registerSingleProto(
	registryURL, subject, schemaSrc string,
	references []map[string]any,
) (int, error) {
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

	return result.ID, nil
}
