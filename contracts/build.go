package contracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ArtifactSourceType artifact source type, local or external
type ArtifactSourceType int

const (
	DefaultBindPkg = "ethereum"

	LocalArtifact ArtifactSourceType = iota
	ExternalArtifact
)

// ContractBuilder is an interface to build bindings to a different chains contracts
type ContractBuilder interface {
	Targets() ([]*ContractTarget, error)
	UpdateExternalSources() error
	Compile(targets []*ContractTarget) error
	GenerateBindings(targets []*ContractTarget) error
}

// ContractTarget contract build target
type ContractTarget struct {
	Version         string
	Contracts       []string
	ExternalRepoDir string
	SourceType      ArtifactSourceType
}

// abigenArgs is the arguments to the abigen executable
type abigenArgs struct {
	ExecutablePath, Bin, ABI, Out, Type, Pkg string
}

// EthereumContractBuilder is a builder struct for Ethereum
type EthereumContractBuilder struct {
	Cfg *config.EthereumSources
	S3  *S3Downloader
}

// NewEthereumContractBuilder builds solidity contracts
func NewEthereumContractBuilder(cfg *config.EthereumSources) *EthereumContractBuilder {
	return &EthereumContractBuilder{
		Cfg: cfg,
		S3:  NewS3Downloader(cfg.Sources.External),
	}
}

// Compile compiles contracts
func (b *EthereumContractBuilder) Compile(_ []*ContractTarget) error {
	// step is not needed for Ethereum, it's more convenient to use hardhat artifacts
	return nil
}

// UpdateExternalSources update all external sources to commits from config
func (b *EthereumContractBuilder) UpdateExternalSources() error {
	err := b.S3.UpdateSources(b.Cfg.Sources.External)
	if err != nil {
		return err
	}
	return nil
}

// Targets what contracts we want to build
func (b *EthereumContractBuilder) Targets() ([]*ContractTarget, error) {
	return []*ContractTarget{
		{
			Contracts:  []string{"LinkToken"},
			Version:    "v0.4",
			SourceType: LocalArtifact,
		},
		{
			Contracts: []string{"APIConsumer", "Oracle", "FluxAggregator",
				"BlockhashStore", "VRF", "VRFConsumer", "VRFCoordinator"},
			Version:    "v0.6",
			SourceType: LocalArtifact,
		},
		{
			Contracts:  []string{"KeeperConsumer", "MockGASAggregator", "MockETHLINKAggregator"},
			Version:    "v0.7",
			SourceType: LocalArtifact,
		},
		{
			Contracts:       []string{"KeeperRegistry", "UpkeepRegistrationRequests"},
			ExternalRepoDir: "keeper",
			SourceType:      ExternalArtifact,
		},
		{
			Contracts:       []string{"OffchainAggregator"},
			ExternalRepoDir: "ocr",
			SourceType:      ExternalArtifact,
		},
	}, nil
}

// FindArtifact searches an artifact file recursively
func (b *EthereumContractBuilder) FindArtifact(rootPath string, name string) (string, error) {
	artifactPaths := make([]string, 0)
	err := filepath.Walk(
		rootPath,
		func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			_, filename := filepath.Split(path)
			if !info.IsDir() && filename == fmt.Sprintf("%s.json", name) {
				log.Debug().Str("File", path).Msg("Found json artifact")
				artifactPaths = append(artifactPaths, path)
			}
			return nil
		})
	if err != nil {
		return "", err
	}
	if len(artifactPaths) > 1 {
		return "", fmt.Errorf("ambiguous artifact paths: %s", artifactPaths)
	}
	if len(artifactPaths) == 0 {
		return "", fmt.Errorf("artifact %s not found", name)
	}
	return artifactPaths[0], nil
}

// ExtractArtifact extracts bytecode and abi from hardhat artifacts
func (b *EthereumContractBuilder) ExtractArtifact(path string) error {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	dir, file := filepath.Split(path)
	var obj map[string]interface{}
	if err := json.Unmarshal(d, &obj); err != nil {
		return err
	}
	log.Debug().Str("Contract", path).Msg("Extracted bytecode")
	fileBin := strings.Replace(file, ".json", ".bin", -1)
	err = ioutil.WriteFile(filepath.Join(dir, fileBin), []byte(obj["bytecode"].(string)), os.ModePerm)
	if err != nil {
		return err
	}
	abiBytes, err := json.Marshal(obj["abi"])
	if err != nil {
		return err
	}
	fileABI := strings.Replace(file, ".json", ".abi", -1)
	err = ioutil.WriteFile(filepath.Join(dir, fileABI), abiBytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// ArtifactSourcePath finds artifact by ContractTarget data
func (b *EthereumContractBuilder) ArtifactSourcePath(t *ContractTarget) (string, error) {
	var path string
	switch t.SourceType {
	case LocalArtifact:
		log.Debug().
			Str("Version", t.Version).
			Msg("Generating bindings")
		path = filepath.Join(b.Cfg.Sources.Local.Path, t.Version)
	case ExternalArtifact:
		extSource, ok := b.Cfg.Sources.External.Repositories[t.ExternalRepoDir]
		if !ok {
			return "", fmt.Errorf("external source %s missing configuration", t.ExternalRepoDir)
		}
		log.Debug().
			Str("Path", b.Cfg.Sources.External.RootPath).
			Str("Dir", extSource.Path).
			Str("S3", b.Cfg.Sources.External.S3URL).
			Msg("Generating bindings for external source")
		path = filepath.Join(b.Cfg.Sources.External.RootPath, extSource.Path)
	}
	return path, nil
}

// GenerateBindings generate go bindings for contracts
func (b *EthereumContractBuilder) GenerateBindings(targets []*ContractTarget) error {
	for _, t := range targets {
		for _, c := range t.Contracts {
			path, err := b.ArtifactSourcePath(t)
			if err != nil {
				return err
			}
			jsonPath, err := b.FindArtifact(path, c)
			if err != nil {
				return err
			}
			err = b.ExtractArtifact(jsonPath)
			if err != nil {
				return err
			}
			dir, _ := filepath.Split(jsonPath)
			abiPath := filepath.Join(dir, fmt.Sprintf("%s.abi", c))
			log.Debug().Str("Path", abiPath).Msg("ABI filepath")
			binPath := filepath.Join(dir, fmt.Sprintf("%s.bin", c))
			log.Debug().Str("Path", binPath).Msg("BIN filepath")
			err = abigen(abigenArgs{
				ExecutablePath: b.Cfg.ExecutablePath,
				ABI:            abiPath,
				Bin:            binPath,
				Out:            filepath.Join(b.Cfg.OutPath, fmt.Sprintf("%s.go", c)),
				Pkg:            c,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// commonBindPackage replaces individual pkg name with celoextended name for convenience
func commonBindPackage(filepath string, pkg string) error {
	read, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	newContents := strings.Replace(string(read), pkg, DefaultBindPkg, 1)
	err = ioutil.WriteFile(filepath, []byte(newContents), 0)
	if err != nil {
		return err
	}
	return nil
}

// abigen calls abigen binary to generate bindings
func abigen(a abigenArgs) error {
	buildCommand := exec.Command(
		a.ExecutablePath,
		"-bin", a.Bin,
		"-abi", a.ABI,
		"-out", a.Out,
		"-pkg", a.Pkg,
	)
	log.Debug().Str("Command", buildCommand.String()).Msg("Build command")
	var buildResponse bytes.Buffer
	buildCommand.Stderr = &buildResponse
	if err := buildCommand.Run(); err != nil {
		return errors.Wrapf(err, "Failed to build contract: %s", a.ABI)
	}
	if err := commonBindPackage(a.Out, a.Pkg); err != nil {
		return errors.Wrapf(err, "Failed to replace package name: %s", a.Pkg)
	}
	return nil
}
