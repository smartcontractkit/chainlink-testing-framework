package framework

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

/* Templates */

const (
	ProductSoakConfigTmpl = `# This file describes how many instances of your product we should deploy for soak test
# you can also override keys from other configs here, for example your [[{{ .ProductName }}]] or [[blockchains]] / [[nodesets]]
[[products]]
name = "ocr2"
instances = 1
`

	ProductBasicConfigTmpl = `[[{{ .ProductName}}]]
# TODO: define your product configuration here, see configurator.go ProductConfig`

	ProductsImplTmpl = `package {{ .ProductName }}

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/{{ .ProductName }}/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "{{ .ProductName }}"}).Logger()

type ProductConfig struct{}

type Configurator struct {
	Config []*ProductConfig ` + "`" + `toml:"productone"` + "`" + `
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, idx int) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateNodesConfig(
	ctx context.Context,
	fs *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) (string, error) {
	L.Info().Msg("Generating Chainlink node config")
	// node
	_ = bc.Out.Nodes[0]
	// chain ID
	_ = bc.Out.ChainID
	return "", nil
}

func (m *Configurator) GenerateNodesSecrets(
	ctx context.Context,
	fs *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) (string, error) {
	L.Info().Msg("Generating Chainlink node secrets")
	// node
	_ = bc.Out.Nodes[0]
	// chain ID
	_ = bc.Out.ChainID
	return "", nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	fake *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) error {
	L.Info().Msg("Configuring product: productone")
	return nil
}
`

	ProductsConfigTmpl = `package products

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvVarTestConfigs = "CTF_CONFIGS"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "product_config"}).Logger()

func Load[T any]() (*T, error) {
	var config T
	paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read product config file path %s: %w", path, err)
		}
		L.Trace().Str("ProductConfig", string(data)).Send()

		decoder := toml.NewDecoder(strings.NewReader(string(data)))

		if err := decoder.Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to decode TOML config, strict mode: %w", err)
		}
	}
	return &config, nil
}

// Store writes config to a file, adds -cache.toml suffix if it's an initial configuration.
func Store[T any](path string, cfg *T) error {
	baseConfigPath, err := BaseConfigPath(EnvVarTestConfigs)
	if err != nil {
		return err
	}
	newCacheName := strings.ReplaceAll(baseConfigPath, ".toml", "")
	var outCacheName string
	if strings.Contains(newCacheName, "cache") {
		L.Info().Str("Cache", baseConfigPath).Msg("Cache file already exists, overriding")
		outCacheName = baseConfigPath
	} else {
		outCacheName = strings.ReplaceAll(baseConfigPath, ".toml", "") + "-out.toml"
	}
	L.Info().Str("OutputFile", outCacheName).Msg("Storing configuration output")
	d, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(path, outCacheName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(d); err != nil {
		return err
	}
	return nil
}

// LoadOutput loads config output file from path.
func LoadOutput[T any](path string) (*T, error) {
	_ = os.Setenv(EnvVarTestConfigs, path)
	return Load[T]()
}

// BaseConfigPath returns base config path, ex. env.toml,overrides.toml -> env.toml.
func BaseConfigPath(envVar string) (string, error) {
	configs := os.Getenv(envVar)
	if configs == "" {
		return "", fmt.Errorf("no %s env var is provided, you should provide at least one test config in TOML", envVar)
	}
	L.Debug().Str("Configs", configs).Msg("Getting base config path")
	return strings.Split(configs, ",")[0], nil
}
`
	ProductsCommonTmpl = `package products

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// WaitMinedFast is a method for Anvil's instant blocks mode to ovecrome bind.WaitMined ticker hardcode.
func WaitMinedFast(ctx context.Context, b bind.DeployBackend, txHash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(5 * time.Millisecond)
	defer queryTicker.Stop()
	for {
		receipt, err := b.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}
	`
)

type ProductCommonParams struct{}

func (g *EnvCodegen) GenerateProductCommon() (string, error) {
	log.Info().Msg("Generating products common")
	p := ProductCommonParams{}
	return render(ProductsCommonTmpl, p)
}

type ProductConfigParams struct{}

func (g *EnvCodegen) GenerateProductsConfig() (string, error) {
	log.Info().Msg("Generating products config")
	p := ProductCommonParams{}
	return render(ProductsConfigTmpl, p)
}

type ProductImplParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateProductImpl() (string, error) {
	log.Info().Msg("Generating product implementation")
	p := ProductImplParams{
		ProductName: g.cfg.productName,
	}
	return render(ProductsImplTmpl, p)
}

type ProductBasicConfigParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateProductBasicConfigParams() (string, error) {
	log.Info().Msg("Generating product basic config")
	p := ProductBasicConfigParams{
		ProductName: g.cfg.productName,
	}
	return render(ProductBasicConfigTmpl, p)
}

type ProductSoakConfigParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateProductSoakConfigParams() (string, error) {
	log.Info().Msg("Generating product soak config")
	p := ProductSoakConfigParams{
		ProductName: g.cfg.productName,
	}
	return render(ProductSoakConfigTmpl, p)
}

// WriteProducts generates a complete products boilerplate
func (g *EnvCodegen) WriteProducts() error {
	productsRoot := filepath.Join(g.cfg.outputDir, "products")
	productRoot := filepath.Join(productsRoot, g.cfg.productName)
	// Create products directory with one product
	if err := os.MkdirAll( //nolint:gosec
		productRoot,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create products directory: %w", err)
	}

	// generate common.go
	commonContents, err := g.GenerateProductCommon()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(productsRoot, "common.go"),
		[]byte(commonContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write products common file: %w", err)
	}

	// generate config.go
	cfgContents, err := g.GenerateProductsConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(productsRoot, "config.go"),
		[]byte(cfgContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write product config file: %w", err)
	}

	// generate configurator.go (product implementation)
	productImplContents, err := g.GenerateProductImpl()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(productRoot, "configurator.go"),
		[]byte(productImplContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write product implementation file: %w", err)
	}

	// generate basic TOML config for product
	basicCfgContents, err := g.GenerateProductBasicConfigParams()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(productRoot, "basic.toml"),
		[]byte(basicCfgContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write product basic config file: %w", err)
	}

	// generate soak TOML config for product
	soakCfgContents, err := g.GenerateProductSoakConfigParams()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(productRoot, "soak.toml"),
		[]byte(soakCfgContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write product soak config file: %w", err)
	}

	// tidy and finalize
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// nolint
	defer os.Chdir(currentDir)
	if err := os.Chdir(g.cfg.outputDir); err != nil {
		return err
	}
	log.Info().Msg("Downloading dependencies and running 'go mod tidy' ..")
	_, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy generated module: %w", err)
	}
	log.Info().
		Str("OutputDir", g.cfg.outputDir).
		Str("Module", g.cfg.moduleName).
		Msg("Developer environment generated")
	return nil
}
