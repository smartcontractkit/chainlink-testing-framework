package framework

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"strings"
)

func NoCache() bool {
	return os.Getenv("CTF_USE_CACHED_OUTPUTS") == "true"
}

func getBaseConfigPath() (string, error) {
	configs := os.Getenv("CTF_CONFIGS")
	if configs == "" {
		return "", fmt.Errorf("no %s env var is provided, you should provide at least one test promtailConfig in TOML", EnvVarTestConfigs)
	}
	return strings.Split(configs, ",")[0], nil
}

func Store[T any](cfg *T) error {
	baseConfigPath, err := getBaseConfigPath()
	if err != nil {
		return err
	}
	cachedOutName := fmt.Sprintf("%s-cache.toml", strings.Replace(baseConfigPath, ".toml", "", -1))
	L.Info().Str("OutputFile", cachedOutName).Msg("Storing configuration output")
	d, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(DefaultConfigDir, cachedOutName), d, os.ModePerm)
}
