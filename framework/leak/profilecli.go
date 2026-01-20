package leak

/*
 * This is a simple utility to download last 12h pprof for memory using Pyroscope's profilecli.
 */

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultPyroscopeBinaryVersion = "1.18.0"
	DefaultOutputPath             = "alloc.pprof"
	DefaultProfileType            = "memory:alloc_space:bytes:space:bytes"
	DefaultFrom                   = "now-12h"
	DefaultTo                     = "now"
)

// ProfileDumperConfig describes profile dump configuration
type ProfileDumperConfig struct {
	PyroscopeURL string
	ServiceName  string
	ProfileType  string
	From         string // Relative time like "now-30m" or unix timestamp (seconds)
	To           string // Relative time like "now" or unix timestamp (seconds)
	OutputPath   string // Output file path
}

// ProfileDumper is responsible for downloading profiles from Pyroscope
type ProfileDumper struct {
	pyroscopeURL   string
	profileCLIPath string
}

// NewProfileDumper creates a new profile dumper instance
func NewProfileDumper(pyroscopeURL string, opts ...func(*ProfileDumper)) *ProfileDumper {
	dumper := &ProfileDumper{
		pyroscopeURL:   pyroscopeURL,
		profileCLIPath: "",
	}
	for _, opt := range opts {
		opt(dumper)
	}
	return dumper
}

// WithProfileCLIPath sets a custom path to profilecli binary
func WithProfileCLIPath(path string) func(*ProfileDumper) {
	return func(d *ProfileDumper) {
		d.profileCLIPath = path
	}
}

// InstallProfileCLI downloads and installs profilecli if not already available
func (d *ProfileDumper) InstallProfileCLI() (string, error) {
	// Check if profilecli is already available
	if d.profileCLIPath != "" {
		if _, err := os.Stat(d.profileCLIPath); err == nil {
			f.L.Info().Str("path", d.profileCLIPath).Msg("Using existing profilecli")
			return d.profileCLIPath, nil
		}
	}
	if path, err := exec.LookPath("profilecli"); err == nil {
		f.L.Info().Str("path", path).Msg("profilecli found in PATH")
		d.profileCLIPath = path
		return path, nil
	}

	// detect OS
	osName := runtime.GOOS
	if osName != "darwin" && osName != "linux" {
		return "", fmt.Errorf("unsupported OS: %s", osName)
	}
	arch := runtime.GOARCH
	if arch != "amd64" && arch != "arm64" {
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/grafana/pyroscope/releases/download/v%s/profilecli_%s_%s_%s.tar.gz",
		DefaultPyroscopeBinaryVersion, DefaultPyroscopeBinaryVersion, osName, arch,
	)

	f.L.Info().Str("url", downloadURL).Msg("Downloading profilecli")

	// Download and extract in current directory
	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -fL %s | tar xvz", downloadURL))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to download and extract profilecli: %w", err)
	}
	binaryPath := "./profilecli"
	if err := os.Chmod(binaryPath, 0o755); err != nil {
		return "", fmt.Errorf("failed to make profilecli executable: %w", err)
	}

	d.profileCLIPath = binaryPath
	f.L.Info().Str("path", binaryPath).Msg("profilecli ready to use")
	return binaryPath, nil
}

func defaults(config *ProfileDumperConfig) {
	if config.OutputPath == "" {
		config.OutputPath = DefaultOutputPath
	}
	if config.ProfileType == "" {
		config.ProfileType = DefaultProfileType
	}
	if config.From == "" {
		config.From = DefaultFrom
	}
	if config.To == "" {
		config.To = DefaultTo
	}
}

// DownloadProfile runs profilecli to download a profile
func (d *ProfileDumper) DownloadProfile(config *ProfileDumperConfig) (string, error) {
	defaults(config)

	cmdArgs := []string{
		"query", "merge",
		fmt.Sprintf("--query={service_name=\"%s\"}", config.ServiceName),
		fmt.Sprintf("--profile-type=%s", config.ProfileType),
		fmt.Sprintf("--from=%s", config.From),
		fmt.Sprintf("--to=%s", config.To),
		fmt.Sprintf("--output=pprof=%s", config.OutputPath),
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("PROFILECLI_URL=%s", d.pyroscopeURL))
	cmd := exec.Command(d.profileCLIPath, cmdArgs...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	f.L.Info().
		Str("command", cmd.String()).
		Str("output", config.OutputPath).
		Msg("Downloading profile from Pyroscope")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to download profile: %w", err)
	}
	if _, err := os.Stat(config.OutputPath); err != nil {
		return "", fmt.Errorf("profile file not created: %w", err)
	}

	fileInfo, _ := os.Stat(config.OutputPath)
	f.L.Info().
		Str("path", config.OutputPath).
		Str("size", fmt.Sprintf("%d bytes", fileInfo.Size())).
		Msg("Profile downloaded successfully")
	return config.OutputPath, nil
}

// MemoryProfile downloads memory profile and saves it
func (d *ProfileDumper) MemoryProfile(config *ProfileDumperConfig) (string, error) {
	if _, err := d.InstallProfileCLI(); err != nil {
		return "", fmt.Errorf("installation failed: %w", err)
	}
	return d.DownloadProfile(config)
}
