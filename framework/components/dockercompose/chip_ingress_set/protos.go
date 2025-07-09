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
	Name          string
	Path          string
	Content       string
	IsImportOnly  bool   // true if this is an import-only schema
	TargetMessage string // for import-only schemas, the message name this schema is for
}

type ProtoSchemaSet struct {
	Owner         string   `toml:"owner"`
	Repository    string   `toml:"repository"`
	Ref           string   `toml:"ref"`            // ref or tag or commit SHA
	Folders       []string `toml:"folders"`        // if not provided, all protos will be fetched, otherwise only protos in these folders will be fetched
	SubjectPrefix string   `toml:"subject_prefix"` // optional prefix for subjects
}

// SubjectNamingStrategyFn is a function that is used to determine the subject name for a given proto file in a given repo
type SubjectNamingStrategyFn func(subjectPrefix string, protoFile protoFile, repoConfig ProtoSchemaSet) (string, error)

// RepositoryToSubjectNamingStrategyFn is a map of repository names to SubjectNamingStrategyFn functions
type RepositoryToSubjectNamingStrategyFn map[string]SubjectNamingStrategyFn

func ValidateRepoConfiguration(repoConfig ProtoSchemaSet) error {
	if repoConfig.Owner == "" {
		return errors.New("owner is required")
	}
	if repoConfig.Repository == "" {
		return errors.New("repo is required")
	}

	if repoConfig.Ref == "" {
		return errors.New("ref is required")
	}

	return nil
}

func DefaultRegisterAndFetchProtos(ctx context.Context, client *github.Client, protoSchemaSets []ProtoSchemaSet, schemaRegistryURL string) error {
	return RegisterAndFetchProtos(ctx, client, protoSchemaSets, schemaRegistryURL, map[string]SubjectNamingStrategyFn{})
}

func RegisterAndFetchProtos(ctx context.Context, client *github.Client, protoSchemaSets []ProtoSchemaSet, schemaRegistryURL string, repoToSubjectNamingStrategy RepositoryToSubjectNamingStrategyFn) error {
	framework.L.Debug().Msgf("Registering and fetching protos from %d repositories", len(protoSchemaSets))

	for _, protoSchemaSet := range protoSchemaSets {
		protos, protosErr := fetchProtoFilesInFolders(ctx, client, protoSchemaSet.Owner, protoSchemaSet.Repository, protoSchemaSet.Ref, protoSchemaSet.Folders)
		if protosErr != nil {
			return errors.Wrapf(protosErr, "failed to fetch protos from %s/%s", protoSchemaSet.Owner, protoSchemaSet.Repository)
		}

		protoMap := make(map[string]string)
		subjects := make(map[string]string)

		for _, proto := range protos {
			protoMap[proto.Path] = proto.Content

			var subjectStrategy SubjectNamingStrategyFn
			if strategy, ok := repoToSubjectNamingStrategy[protoSchemaSet.Owner+"/"+protoSchemaSet.Repository]; ok {
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
			return errors.Wrapf(registerErr, "failed to register protos from %s/%s", protoSchemaSet.Owner, protoSchemaSet.Repository)
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

// createSchemasForProto creates schemas for a proto file:
// 1. The original proto file (registered first, Red Panda uses first message by default)
// 2. Import-only schemas for each additional message (2nd, 3rd, etc., skipping the first)
func createSchemasForProto(proto protoFile) ([]protoFile, error) {
	messageNames, err := extractTopLevelMessageNamesWithRegex(proto.Content)
	if err != nil {
		return nil, err
	}

	var schemas []protoFile

	// First, add the original proto file
	schemas = append(schemas, proto)

	// If there are multiple messages, create import-only schemas for messages after the first
	if len(messageNames) > 1 {
		framework.L.Debug().Msgf("Creating import-only schemas for %d additional messages in %s: %v", len(messageNames)-1, proto.Path, messageNames[1:])

		// Skip the first message (index 0) since it's covered by the original proto
		for _, messageName := range messageNames[1:] {
			importSchema, err := createImportOnlySchema(proto, messageName)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create import-only schema for message %s", messageName)
			}
			schemas = append(schemas, importSchema)
		}
	}

	return schemas, nil
}

// createImportOnlySchema creates an import-only schema for a specific message
func createImportOnlySchema(originalProto protoFile, messageName string) (protoFile, error) {
	// Extract syntax and package from original proto
	syntax, err := extractSyntaxDeclaration(originalProto.Content)
	if err != nil {
		return protoFile{}, err
	}

	packageName, err := extractPackageNameWithRegex(originalProto.Content)
	if err != nil {
		return protoFile{}, err
	}

	// Create import-only schema content
	var content strings.Builder
	content.WriteString(syntax)
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("package %s;\n\n", packageName))
	content.WriteString(fmt.Sprintf("import \"%s\";\n", originalProto.Path))

	// Create new schema with clear prefix: import_MessageName_originalname.proto
	baseNameWithoutExt := strings.TrimSuffix(originalProto.Name, ".proto")
	newPath := fmt.Sprintf("import_only_%s_%s.proto", messageName, baseNameWithoutExt)
	newName := fmt.Sprintf("import_only_%s_%s.proto", messageName, baseNameWithoutExt)

	return protoFile{
		Name:          newName,
		Path:          newPath,
		Content:       content.String(),
		IsImportOnly:  true,
		TargetMessage: messageName,
	}, nil
}

// extractSyntaxDeclaration extracts the syntax declaration from a proto file
func extractSyntaxDeclaration(protoContent string) (string, error) {
	syntaxMatch := regexp.MustCompile(`(?m)^syntax\s*=\s*"[^"]+"\s*;`).FindString(protoContent)
	if syntaxMatch == "" {
		// Default to proto3 if no syntax specified
		return `syntax = "proto3";`, nil
	}
	return syntaxMatch, nil
}

// Fetches .proto files from a GitHub repo optionally scoped to specific folders. It is recommended to use `*github.Client` with auth token to avoid rate limiting.
func fetchProtoFilesInFolders(ctx context.Context, client *github.Client, owner, repository, ref string, folders []string) ([]protoFile, error) {
	framework.L.Debug().Msgf("Fetching proto files from %s/%s in folders: %s", owner, repository, strings.Join(folders, ", "))

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

	framework.L.Debug().Msgf("Fetched %d proto files from %s/%s", len(files), owner, repository)

	if len(files) == 0 {
		return nil, fmt.Errorf("no proto files found in %s/%s in folders %s", owner, repository, strings.Join(folders, ", "))
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
