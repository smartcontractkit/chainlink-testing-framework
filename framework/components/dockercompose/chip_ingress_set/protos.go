package chipingressset

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v72/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type SchemaSet struct {
	URI        string `toml:"uri"`
	Ref        string `toml:"ref"`         // ref or tag or commit SHA
	SchemaDir  string `toml:"schema_dir"`  // optional sub-directory in the repo where protos are located
	ConfigFile string `toml:"config_file"` // optional path to config file in the repo (default: <schemaDir>/chip.json)
}

func (s *SchemaSet) ConfigFileName() string {
	if s.ConfigFile != "" {
		return s.ConfigFile
	}
	return "chip.json"
}

func FetchAndRegisterProtos(ctx context.Context, client *github.Client, chipIngressOutput *ChipIngressOutput, schemaSet []SchemaSet) error {
	framework.L.Info().Msgf("Registering and fetching schemas from %d repositories", len(schemaSet))

	for _, set := range schemaSet {
		framework.L.Debug().Msgf("Processing schema set: %s", set.URI)
		if valErr := validateSchemaSet(set); valErr != nil {
			return errors.Wrapf(valErr, "invalid repo configuration for schema set: %v", set)
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

	for _, set := range schemaSet {
		repoPath, repoErr := getSchemaRepository(ctx, ghClientFn, set.URI, set.Ref)
		if repoErr != nil {
			return errors.Wrapf(repoErr, "failed to get repository %s", set.URI)
		}

		schemaDir := filepath.Join(repoPath, set.SchemaDir)
		configFilePath := filepath.Join(schemaDir, set.ConfigFileName())

		registerErr := registerWithChipConfigService(ctx, chipIngressOutput, schemaDir, configFilePath)
		if registerErr != nil {
			return errors.Wrapf(registerErr, "failed to register schemas from '%s' using Chip Config", set.URI)
		}
	}

	return nil
}

func registerWithChipConfigService(ctx context.Context, chipIngressOutput *ChipIngressOutput, schemaDir, configFilePath string) error {
	registrationConfig, schemas, rErr := parseSchemaConfig(configFilePath, schemaDir)
	if rErr != nil {
		return fmt.Errorf("failed to parse schema config: %w", rErr)
	}

	fmt.Printf("📋 Parsed %d schema(s) from \033[1m%s\033[0m\n", len(schemas), configFilePath)
	fmt.Printf("✅ All entity names validated successfully\n\n")

	pbSchemas := convertToPbSchemas(schemas, registrationConfig.Domain)

	client, cErr := chipingress.NewClient(chipIngressOutput.GRPCExternalURL)
	if cErr != nil {
		return fmt.Errorf("failed to create Chip client: %w", cErr)
	}

	_, regErr := client.RegisterSchemas(ctx, pbSchemas...)
	if regErr != nil {
		return fmt.Errorf("failed to register schemas: %w", regErr)
	}

	fmt.Printf("✅ Registered %d schema(s)\n", len(pbSchemas))
	for _, schema := range schemas {
		subject := fmt.Sprintf("%s-%s", registrationConfig.Domain, schema.Entity)
		fmt.Printf("   └── \033[1m%s\033[0m (\033[34m%s\033[0m)\n", subject, schema.Path)
	}
	fmt.Printf("\n")

	return nil
}

func validateSchemaSet(set SchemaSet) error {
	if set.URI == "" {
		return errors.New("uri is required")
	}

	if !strings.HasPrefix(set.URI, "https://") && !strings.HasPrefix(set.URI, "file://") {
		return errors.New("uri has to start with either 'file://' or 'https://'")
	}

	if strings.HasPrefix(set.URI, "file://") {
		if set.Ref != "" {
			return errors.New("ref is not supported with local protos with 'file://' prefix")
		}
		return nil
	}

	trimmedURI := strings.TrimPrefix(set.URI, "https://")
	if !strings.HasPrefix(trimmedURI, "github.com") {
		return fmt.Errorf("only repositories hosted at github.com are supported, but %s was found", set.URI)
	}

	parts := strings.Split(trimmedURI, "/")
	if len(parts) < 3 {
		return fmt.Errorf("URI should have following format: 'https://github.com/<OWNER>/<REPOSITORY>', but %s was found", set.URI)
	}

	if set.Ref == "" {
		return errors.New("ref is required, when fetching protos from Github repository")
	}

	return nil
}

func getSchemaRepository(ctx context.Context, clientFn func() *github.Client, uri, ref string) (string, error) {
	uriParts := strings.Split(strings.TrimPrefix(uri, "https://"), "/")
	if pathClean, ok := strings.CutPrefix(uri, "file://"); ok {
		if _, err := os.Stat(pathClean); err == nil {
			if hasFiles, hasErr := hasProtoFiles(pathClean); hasErr != nil {
				return "", fmt.Errorf("failed to check for proto files in cache at %s: %w", pathClean, hasErr)
			} else if !hasFiles {
				return fetchFromGithub(ctx, clientFn, uriParts[1], uriParts[2], ref) // cache is invalid, download from GitHub
			}

			abs, absErr := filepath.Abs(pathClean)
			if absErr != nil {
				return "", errors.Wrapf(absErr, "failed to get absolute path for %s", pathClean)
			}
			return abs, nil
		}
	}

	return fetchFromGithub(ctx, clientFn, uriParts[1], uriParts[2], ref)
}

type githubFile struct {
	Name    string
	Path    string
	Content string
}

func cacheFilePath(owner, repository, ref string) (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", errors.Wrap(homeErr, "failed to get user home directory")
	}
	return filepath.Join(homeDir, ".local", "share", "beholder", "protobufs", owner, repository, ref), nil
}

func resolveRefSHA(ctx context.Context, client *github.Client, owner, repository, ref string) (string, error) {
	// These lookups are expected to 404 when ref is a commit SHA; a miss only
	// means "try the next form of ref".
	for _, prefix := range []string{"refs/tags/", "refs/heads/"} {
		refObj, _, err := client.Git.GetRef(ctx, owner, repository, prefix+ref)
		if err == nil {
			return refObj.GetObject().GetSHA(), nil
		}
	}

	commit, err := getCommitWithRetry(ctx, client, owner, repository, ref)
	if err != nil {
		if isNotFound(err) {
			return "", fmt.Errorf("ref %q not found", ref)
		}
		// Surface the underlying failure: a rate limit or server error is not
		// a missing ref, and must not be reported as one.
		return "", fmt.Errorf("resolve commit %q: %w", ref, err)
	}
	return commit.GetSHA(), nil
}

const (
	commitLookupAttempts = 4
	commitLookupBaseWait = 500 * time.Millisecond
)

// retryWait is a variable so tests can avoid sleeping.
var retryWait = func(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func getCommitWithRetry(ctx context.Context, client *github.Client, owner, repository, ref string) (*github.RepositoryCommit, error) {
	var lastErr error
	for attempt := 0; attempt < commitLookupAttempts; attempt++ {
		commit, _, err := client.Repositories.GetCommit(ctx, owner, repository, ref, nil)
		if err == nil {
			return commit, nil
		}
		lastErr = err
		if attempt == commitLookupAttempts-1 || !isRetryable(err) {
			break
		}
		if waitErr := retryWait(ctx, time.Duration(attempt+1)*commitLookupBaseWait); waitErr != nil {
			return nil, waitErr
		}
	}
	return nil, lastErr
}

func isNotFound(err error) bool {
	var errResponse *github.ErrorResponse
	return errors.As(err, &errResponse) &&
		errResponse.Response != nil &&
		errResponse.Response.StatusCode == http.StatusNotFound
}

func isRetryable(err error) bool {
	var errResponse *github.ErrorResponse
	if errors.As(err, &errResponse) {
		if errResponse.Response == nil {
			return false
		}
		switch code := errResponse.Response.StatusCode; {
		case code >= 500, code == http.StatusTooManyRequests:
			return true
		case code == http.StatusForbidden:
			// Secondary rate limits come back as 403s.
			return strings.Contains(strings.ToLower(errResponse.Message), "rate limit")
		}
		return false
	}
	// Transport failures (DNS, TLS, timeouts) are worth another attempt.
	var netErr net.Error
	return errors.As(err, &netErr)
}

func fetchFromGithub(ctx context.Context, clientFn func() *github.Client, owner, repository, ref string) (string, error) {
	cachePath, found, cacheErr := getCachedProtoFiles(owner, repository, ref)
	if cacheErr == nil && found {
		framework.L.Debug().Msgf("Using cached proto files for %s/%s at ref %s", owner, repository, ref)
		return cachePath, nil
	}
	if cacheErr != nil {
		framework.L.Warn().Msgf("Failed to load cached proto files for %s/%s at ref %s: %v", owner, repository, ref, cacheErr)
	}

	client := clientFn()
	var files []githubFile

	sha, shaErr := resolveRefSHA(ctx, client, owner, repository, ref)
	if shaErr != nil {
		return "", errors.Wrapf(shaErr, "cannot resolve ref %q", ref)
	}

	tree, _, treeErr := client.Git.GetTree(ctx, owner, repository, sha, true)
	if treeErr != nil {
		return "", errors.Wrap(treeErr, "failed to fetch tree")
	}

	for _, entry := range tree.Entries {
		// skip non-blob entries and non-proto or json files [JSON describes the schemas to be registered]
		if entry.GetType() != "blob" || entry.Path == nil || (!strings.HasSuffix(*entry.Path, ".proto") && !strings.HasSuffix(*entry.Path, ".json")) {
			continue
		}

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repository, sha, *entry.Path)
		resp, respErr := http.Get(rawURL)
		if respErr != nil {
			return "", errors.Wrapf(respErr, "failed to fetch %s", *entry.Path)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return "", errors.Errorf("bad status from GitHub for %s: %d", *entry.Path, resp.StatusCode)
		}

		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return "", errors.Wrapf(bodyErr, "failed to read body for %s", *entry.Path)
		}

		files = append(files, githubFile{
			Name:    filepath.Base(*entry.Path),
			Path:    *entry.Path,
			Content: string(body),
		})
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no proto files found in %s/%s", owner, repository)
	}

	framework.L.Debug().Msgf("Fetched %d files from %s/%s", len(files), owner, repository)

	savedPath, saveErr := saveFilesToCache(owner, repository, ref, files)
	if saveErr != nil {
		framework.L.Warn().Msgf("Failed to save files to cache for %s/%s at ref %s: %v", owner, repository, ref, saveErr)
	}

	return savedPath, nil
}

func getCachedProtoFiles(owner, repository, ref string) (string, bool, error) {
	cachePath, cacheErr := cacheFilePath(owner, repository, ref)
	if cacheErr != nil {
		return "", false, errors.Wrapf(cacheErr, "failed to get cache file path for %s/%s at ref %s", owner, repository, ref)
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return "", false, nil // cache not found
	}

	if hasFiles, hasErr := hasProtoFiles(cachePath); hasErr != nil {
		return "", false, fmt.Errorf("failed to check for proto files in cache at %s: %w", cachePath, hasErr)
	} else if !hasFiles {
		return "", false, nil // cache is invalid
	}

	return cachePath, true, nil
}

func hasProtoFiles(dir string) (bool, error) {
	found := false
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".proto") {
			found = true
			return filepath.SkipAll // stop walking immediately
		}
		return nil
	})
	return found, err
}

func saveFilesToCache(owner, repository, ref string, files []githubFile) (string, error) {
	cachePath, cacheErr := cacheFilePath(owner, repository, ref)
	if cacheErr != nil {
		return "", errors.Wrapf(cacheErr, "failed to get cache file path for %s/%s at ref %s", owner, repository, ref)
	}

	for _, file := range files {
		path := filepath.Join(cachePath, file.Path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory for cache file %s", path)
		}
		if writeErr := os.WriteFile(path, []byte(file.Content), 0755); writeErr != nil {
			return "", errors.Wrapf(writeErr, "failed to write cached proto file to %s", path)
		}
	}

	framework.L.Debug().Msgf("Saved %d proto files to cache at %s", len(files), cachePath)
	return cachePath, nil
}
