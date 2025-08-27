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
	"golang.org/x/oauth2"
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
	ExcludeFiles  []string `toml:"exclude_files"`  // files to exclude from registration (e.g., ['workflows/v2/execution_status.proto'])
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
	framework.L.Info().Msgf("Registering and fetching protos from %d repositories", len(protoSchemaSets))

	for _, protoSchemaSet := range protoSchemaSets {
		framework.L.Debug().Msgf("Processing proto schema set: %s", protoSchemaSet.URI)
		if len(protoSchemaSet.ExcludeFiles) > 0 {
			framework.L.Debug().Msgf("Excluding files: %s", strings.Join(protoSchemaSet.ExcludeFiles, ", "))
		}
		if valErr := validateRepoConfiguration(protoSchemaSet); valErr != nil {
			return errors.Wrapf(valErr, "invalid repo configuration for schema set: %v", protoSchemaSet)
		}
	}

	ghClientFn := func() *github.Client {
		if client != nil {
			return client
		}

		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			tc := oauth2.NewClient(ctx, ts)
			return github.NewClient(tc)
		}

		framework.L.Warn().Msg("GITHUB_TOKEN is not set, using unauthenticated GitHub client. This may cause rate limiting issues when downloading proto files")
		return github.NewClient(nil)
	}

	for _, protoSchemaSet := range protoSchemaSets {
		protos, protosErr := fetchProtoFilesInFolders(ctx, ghClientFn, protoSchemaSet.URI, protoSchemaSet.Ref, protoSchemaSet.Folders, protoSchemaSet.ExcludeFiles)
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

		registerErr := registerAllWithTopologicalSorting(schemaRegistryURL, protoMap, subjects, protoSchemaSet.Folders)
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

// extractTopLevelMessageNamesWithRegex extracts top-level message and enum names from a proto file using regex.
func extractTopLevelMessageNamesWithRegex(protoSrc string) ([]string, error) {
	// Extract message names
	messageMatches := regexp.MustCompile(`(?m)^\s*message\s+(\w+)\s*{`).FindAllStringSubmatch(protoSrc, -1)
	var names []string
	for _, match := range messageMatches {
		if len(match) >= 2 {
			names = append(names, match[1])
		}
	}

	// Extract enum names
	enumMatches := regexp.MustCompile(`(?m)^\s*enum\s+(\w+)\s*{`).FindAllStringSubmatch(protoSrc, -1)
	for _, match := range enumMatches {
		if len(match) >= 2 {
			names = append(names, match[1])
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no message or enum names found in proto source")
	}

	return names, nil
}

// extractImportStatements extracts import statements from a proto source file using regex.
func extractImportStatements(protoSrc string) []string {
	matches := regexp.MustCompile(`(?m)^\s*import\s+"([^"]+)"\s*;`).FindAllStringSubmatch(protoSrc, -1)
	var imports []string
	for _, match := range matches {
		if len(match) >= 2 {
			imports = append(imports, match[1])
		}
	}
	return imports
}

// fetchProtoFilesInFolders fetches .proto files from a GitHub repo optionally scoped to specific folders.
// It is recommended to use `*github.Client` with auth token to avoid rate limiting.
func fetchProtoFilesInFolders(ctx context.Context, clientFn func() *github.Client, uri, ref string, folders []string, excludeFiles []string) ([]protoFile, error) {
	if strings.HasPrefix(uri, "file://") {
		return fetchProtosFromFilesystem(uri, folders, excludeFiles)
	}

	parts := strings.Split(strings.TrimPrefix(uri, "https://"), "/")
	return fetchProtosFromGithub(ctx, clientFn, parts[1], parts[2], ref, folders, excludeFiles)
}

func fetchProtosFromGithub(ctx context.Context, clientFn func() *github.Client, owner, repository, ref string, folders []string, excludeFiles []string) ([]protoFile, error) {
	cachedFiles, found, cacheErr := loadCachedProtoFiles(owner, repository, ref, folders, excludeFiles)
	if cacheErr == nil && found {
		framework.L.Debug().Msgf("Using cached proto files for %s/%s at ref %s", owner, repository, ref)
		return cachedFiles, nil
	}
	if cacheErr != nil {
		framework.L.Warn().Msgf("Failed to load cached proto files for %s/%s at ref %s: %v", owner, repository, ref, cacheErr)
	}

	client := clientFn()
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
		if len(folders) > 0 {
			matched := false
			for _, folder := range folders {
				if strings.HasPrefix(*entry.Path, strings.TrimSuffix(folder, "/")+"/") {
					matched = true
					break
				}
			}
			if !matched {
				continue searchLoop
			}
		}

		// if excludeFiles are specified, check if the file should be excluded
		if len(excludeFiles) > 0 {
			excluded := false
			for _, exclude := range excludeFiles {
				if strings.HasPrefix(*entry.Path, exclude) {
					framework.L.Debug().Msgf("Excluding proto file %s (matches exclude pattern: %s)", *entry.Path, exclude)
					excluded = true
					break
				}
			}
			if excluded {
				continue searchLoop
			}
		}

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repository, sha, *entry.Path)
		resp, respErr := http.Get(rawURL)
		if respErr != nil {
			return nil, errors.Wrapf(respErr, "failed to fetch %s", *entry.Path)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, errors.Errorf("bad status from GitHub for %s: %d", *entry.Path, resp.StatusCode)
		}

		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, errors.Wrapf(bodyErr, "failed to read body for %s", *entry.Path)
		}

		files = append(files, protoFile{
			Name:    filepath.Base(*entry.Path),
			Path:    *entry.Path,
			Content: string(body),
		})
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no proto files found in %s/%s in folders %s", owner, repository, strings.Join(folders, ", "))
	}

	framework.L.Debug().Msgf("Fetched %d proto files from %s/%s", len(files), owner, repository)

	saveErr := saveProtoFilesToCache(owner, repository, ref, files)
	if saveErr != nil {
		framework.L.Warn().Msgf("Failed to save proto files to cache for %s/%s at ref %s: %v", owner, repository, ref, saveErr)
	}

	return files, nil
}

func loadCachedProtoFiles(owner, repository, ref string, folders []string, excludeFiles []string) ([]protoFile, bool, error) {
	cachePath, cacheErr := cacheFilePath(owner, repository, ref)
	if cacheErr != nil {
		return nil, false, errors.Wrapf(cacheErr, "failed to get cache file path for %s/%s at ref %s", owner, repository, ref)
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, false, nil // cache not found
	}

	cachedFiles, cachedErr := fetchProtosFromFilesystem("file://"+cachePath, folders, excludeFiles)
	if cachedErr != nil {
		return nil, false, errors.Wrapf(cachedErr, "failed to load cached proto files from %s", cachePath)
	}

	return cachedFiles, true, nil
}

func saveProtoFilesToCache(owner, repository, ref string, files []protoFile) error {
	cachePath, cacheErr := cacheFilePath(owner, repository, ref)
	if cacheErr != nil {
		return errors.Wrapf(cacheErr, "failed to get cache file path for %s/%s at ref %s", owner, repository, ref)
	}

	for _, file := range files {
		path := filepath.Join(cachePath, file.Path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory for cache file %s", path)
		}
		if writeErr := os.WriteFile(path, []byte(file.Content), 0755); writeErr != nil {
			return errors.Wrapf(writeErr, "failed to write cached proto file to %s", path)
		}
	}

	framework.L.Debug().Msgf("Saved %d proto files to cache at %s", len(files), cachePath)
	return nil
}

func cacheFilePath(owner, repository, ref string) (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", errors.Wrap(homeErr, "failed to get user home directory")
	}
	return filepath.Join(homeDir, ".local", "share", "beholder", "protobufs", owner, repository, ref), nil
}

func fetchProtosFromFilesystem(uri string, folders []string, excludeFiles []string) ([]protoFile, error) {
	var files []protoFile
	protoDirPath := strings.TrimPrefix(uri, "file://")

	walkErr := filepath.Walk(protoDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".proto") {
			return nil
		}

		relativePath := strings.TrimPrefix(strings.TrimPrefix(path, protoDirPath), "/")

		// if folders are specified, check prefix match
		if len(folders) > 0 {
			matched := false
			for _, folder := range folders {
				if strings.HasPrefix(relativePath, folder) {
					matched = true
					break
				}
			}
			if !matched {
				return nil
			}
		}

		// if excludeFiles are specified, check if the file should be excluded
		if len(excludeFiles) > 0 {
			excluded := false
			for _, exclude := range excludeFiles {
				if strings.HasPrefix(relativePath, exclude) {
					framework.L.Debug().Msgf("Excluding proto file %s (matches exclude pattern: %s)", relativePath, exclude)
					excluded = true
					break
				}
			}
			if excluded {
				return nil
			}
		}

		content, contentErr := os.ReadFile(path)
		if contentErr != nil {
			return errors.Wrapf(contentErr, "failed to read file at %s", path)
		}

		files = append(files, protoFile{
			Name:    filepath.Base(path),
			Path:    relativePath,
			Content: string(content),
		})

		return nil
	})

	if walkErr != nil {
		return nil, errors.Wrapf(walkErr, "failed to walk through directory %s", protoDirPath)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no proto files found in '%s' in folders %s", protoDirPath, strings.Join(folders, ", "))
	}

	framework.L.Debug().Msgf("Fetched %d proto files from local %s", len(files), protoDirPath)
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

// registerAllWithTopologicalSorting registers protos in dependency order using topological sorting
func registerAllWithTopologicalSorting(
	schemaRegistryURL string,
	protoMap map[string]string, // path -> proto source
	subjectMap map[string]string, // path -> subject
	folders []string, // folders configuration used to determine import prefix transformations
) error {
	framework.L.Info().Msgf("Registering %d protobuf schemas", len(protoMap))

	// Build dependency graph and sort topologically
	dependencies, depErr := buildDependencyGraph(protoMap)
	if depErr != nil {
		return errors.Wrap(depErr, "failed to build dependency graph")
	}

	sortedFiles, sortErr := topologicalSort(dependencies)
	if sortErr != nil {
		return errors.Wrap(sortErr, "failed to sort files topologically")
	}

	framework.L.Debug().Msgf("Registration order (topologically sorted): %v", sortedFiles)

	schemas := map[string]*schemaStatus{}
	for path, src := range protoMap {
		schemas[path] = &schemaStatus{Source: src}
	}

	// Register files in topological order
	for _, path := range sortedFiles {
		schema, exists := schemas[path]
		if !exists {
			framework.L.Warn().Msgf("File %s not found in schemas map", path)
			continue
		}

		if schema.Registered {
			continue
		}

		subject, ok := subjectMap[path]
		if !ok {
			return fmt.Errorf("no subject found for %s", path)
		}

		// Determine which folder prefixes should be stripped based on configuration
		prefixesToStrip := determineFolderPrefixesToStrip(folders)

		// Build references only for files that have dependencies
		var fileRefs []map[string]any
		if deps, hasDeps := dependencies[path]; hasDeps && len(deps) > 0 {
			for _, dep := range deps {
				if depSubject, depExists := subjectMap[dep]; depExists {
					// The schema registry expects import names without the configured folder prefixes
					// So if folders=["workflows"] and the import is "workflows/v1/metadata.proto",
					// the name should be "v1/metadata.proto"
					importName := stripFolderPrefix(dep, prefixesToStrip)

					fileRefs = append(fileRefs, map[string]any{
						"name":    importName,
						"subject": depSubject,
						"version": 1,
					})
				}
			}
		}

		// Check if schema is already registered
		if existingID, exists := checkSchemaExists(schemaRegistryURL, subject); exists {
			framework.L.Debug().Msgf("Schema %s already exists with ID %d, skipping registration", subject, existingID)
			schema.Registered = true
			schema.Version = existingID
			continue
		}

		// The schema registry expects import statements without the configured folder prefixes
		// Transform the schema content to remove these prefixes from import statements
		modifiedSchema := transformSchemaContent(schema.Source, prefixesToStrip)

		_, registerErr := registerSingleProto(schemaRegistryURL, subject, modifiedSchema, fileRefs)
		if registerErr != nil {
			return errors.Wrapf(registerErr, "failed to register %s as %s", path, subject)
		}

		schema.Registered = true
		schema.Version = 1

		framework.L.Info().Msgf("✔ Registered: %s as %s", path, subject)
	}

	framework.L.Info().Msgf("✅ Successfully registered %d schemas", len(protoMap))
	return nil
}

// checkSchemaExists checks if a schema already exists in the registry
func checkSchemaExists(registryURL, subject string) (int, bool) {
	url := fmt.Sprintf("%s/subjects/%s/versions", registryURL, subject)

	resp, err := http.Get(url)
	if err != nil {
		framework.L.Debug().Msgf("Failed to check schema existence for %s: %v", subject, err)
		return 0, false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var versions []struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
			framework.L.Debug().Msgf("Failed to decode versions for %s: %v", subject, err)
			return 0, false
		}
		if len(versions) > 0 {
			return versions[len(versions)-1].ID, true
		}
	}

	return 0, false
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

// determineFolderPrefixesToStrip determines which folder prefixes should be stripped from import paths
// based on the folders configuration. The schema registry expects import names to be relative to the
// configured folders, so we strip these prefixes to make imports work correctly.
func determineFolderPrefixesToStrip(folders []string) []string {
	var prefixes []string
	for _, folder := range folders {
		// Ensure folder ends with / for prefix matching
		prefix := strings.TrimSuffix(folder, "/") + "/"
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}

// stripFolderPrefix removes any configured folder prefixes from the given path
func stripFolderPrefix(path string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return strings.TrimPrefix(path, prefix)
		}
	}
	return path
}

// transformSchemaContent removes folder prefixes from import statements in protobuf source
func transformSchemaContent(content string, prefixes []string) string {
	modified := content
	for _, prefix := range prefixes {
		// Transform import statements like "workflows/v1/" to "v1/"
		modified = strings.ReplaceAll(modified, `"`+prefix, `"`)
	}
	return modified
}

// buildDependencyGraph builds a dependency graph from protobuf files
func buildDependencyGraph(protoMap map[string]string) (map[string][]string, error) {
	dependencies := make(map[string][]string)

	framework.L.Debug().Msgf("Building dependency graph for %d proto files", len(protoMap))

	// Initialize dependencies map
	for path := range protoMap {
		dependencies[path] = []string{}
	}

	// Parse imports and build dependency graph
	for path, content := range protoMap {
		imports := extractImportStatements(content)

		for _, importPath := range imports {
			if strings.HasPrefix(importPath, "google/protobuf/") {
				// Skip Google protobuf imports as they're not in our protoMap
				continue
			}

			// Check if this import exists in our protoMap
			if _, exists := protoMap[importPath]; exists {
				// Check for self-reference - this indicates either an invalid proto file
				// or a potential bug in our import/path handling
				if importPath == path {
					framework.L.Warn().Msgf("Self-reference detected: file %s imports itself (import: %s). This suggests either an invalid proto file or a path normalization issue. Skipping this dependency to avoid cycles.", path, importPath)
					// Continue without adding the dependency to avoid cycles, but don't fail registration
					// as this might be a recoverable issue or edge case
					continue
				}

				dependencies[path] = append(dependencies[path], importPath)
			} else {
				framework.L.Warn().Msgf("Import %s in %s not found in protoMap", importPath, path)
			}
		}
	}

	return dependencies, nil
}

// topologicalSort performs topological sorting using Kahn's algorithm
func topologicalSort(dependencies map[string][]string) ([]string, error) {
	// Calculate in-degrees (how many files each file depends on)
	inDegree := make(map[string]int)
	for file := range dependencies {
		inDegree[file] = 0
	}

	// Count dependencies for each file
	for file, deps := range dependencies {
		inDegree[file] = len(deps)
	}

	// Find files with no dependencies (in-degree = 0)
	var queue []string
	for file, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, file)
		}
	}

	var result []string
	for len(queue) > 0 {
		file := queue[0]
		queue = queue[1:]
		result = append(result, file)

		// Reduce in-degree for files that depend on the current file
		for dependent, deps := range dependencies {
			for _, dep := range deps {
				if dep == file {
					inDegree[dependent]--
					if inDegree[dependent] == 0 {
						queue = append(queue, dependent)
					}
				}
			}
		}
	}

	// Check for cycles
	if len(result) != len(dependencies) {
		return nil, fmt.Errorf("circular dependency detected in protobuf files")
	}

	return result, nil
}
