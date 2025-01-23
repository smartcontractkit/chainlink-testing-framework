package seth

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

const (
	ErrNoNetwork = "no network specified, use -n flag. Ex.: 'seth -n Geth stats' or -u and -c flags. Ex.: 'seth -u http://localhost:8545 -c 1337 stats'"
)

var C *seth.Client

func RunCLI(args []string) error {
	app := &cli.App{
		Name:      "seth",
		Version:   "v1.0.0",
		Usage:     "seth CLI",
		UsageText: `utility to create and control Ethereum keys and give you more debug info about chains`,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "networkName", Aliases: []string{"n"}},
			&cli.StringFlag{Name: "url", Aliases: []string{"u"}},
		},
		Before: func(cCtx *cli.Context) error {
			networkName := cCtx.String("networkName")
			url := cCtx.String("url")
			if networkName == "" && url == "" {
				return errors.New(ErrNoNetwork)
			}
			if networkName != "" {
				_ = os.Setenv(seth.NETWORK_ENV_VAR, networkName)
			}
			if url != "" {
				_ = os.Setenv(seth.URL_ENV_VAR, url)
			}
			if cCtx.Args().Len() > 0 && cCtx.Args().First() != "trace" {
				var err error
				switch cCtx.Args().First() {
				case "gas", "stats":
					var cfg *seth.Config
					var pk string
					_, pk, err = seth.NewAddress()
					if err != nil {
						return err
					}

					err = os.Setenv(seth.ROOT_PRIVATE_KEY_ENV_VAR, pk)
					if err != nil {
						return err
					}

					cfg, err = seth.ReadConfig()
					if err != nil {
						return err
					}
					C, err = seth.NewClientWithConfig(cfg)
					if err != nil {
						return err
					}
				case "trace":
					return nil
				}
				if err != nil {
					return err
				}
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "stats",
				HelpName:    "stats",
				Aliases:     []string{"s"},
				Description: "get various network related stats",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "start_block", Aliases: []string{"s"}},
					&cli.Int64Flag{Name: "end_block", Aliases: []string{"e"}},
				},
				Action: func(cCtx *cli.Context) error {
					start := cCtx.Int64("start_block")
					end := cCtx.Int64("end_block")
					if start == 0 {
						return fmt.Errorf("at least start block should be defined, ex.: -s -10")
					}
					if start > 0 && end == 0 {
						return fmt.Errorf("invalid block params. Last N blocks example: -s -10, interval example: -s 10 -e 20")
					}
					cs, err := seth.NewBlockStats(C)
					if err != nil {
						return err
					}
					return cs.Stats(big.NewInt(start), big.NewInt(end))
				},
			},
			{
				Name:        "gas",
				HelpName:    "gas",
				Aliases:     []string{"g"},
				Description: "get various info about gas prices",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "blocks", Aliases: []string{"b"}},
					&cli.Float64Flag{Name: "tipPercentile", Aliases: []string{"tp"}},
				},
				Action: func(cCtx *cli.Context) error {
					ge := seth.NewGasEstimator(C)
					blocks := cCtx.Uint64("blocks")
					tipPerc := cCtx.Float64("tipPercentile")
					stats, err := ge.Stats(blocks, tipPerc)
					if err != nil {
						return err
					}
					seth.L.Info().
						Interface("Max", stats.GasPrice.Max).
						Interface("99", stats.GasPrice.Perc99).
						Interface("75", stats.GasPrice.Perc75).
						Interface("50", stats.GasPrice.Perc50).
						Interface("25", stats.GasPrice.Perc25).
						Msg("Base fee (Wei)")
					seth.L.Info().
						Interface("Max", stats.TipCap.Max).
						Interface("99", stats.TipCap.Perc99).
						Interface("75", stats.TipCap.Perc75).
						Interface("50", stats.TipCap.Perc50).
						Interface("25", stats.TipCap.Perc25).
						Msg("Priority fee (Wei)")
					seth.L.Info().
						Interface("GasPrice", stats.SuggestedGasPrice).
						Msg("Suggested gas price now")
					seth.L.Info().
						Interface("GasTipCap", stats.SuggestedGasTipCap).
						Msg("Suggested gas tip cap now")

					type asTomlCfg struct {
						GasPrice int64 `toml:"gas_price"`
						GasTip   int64 `toml:"gas_tip_cap"`
						GasFee   int64 `toml:"gas_fee_cap"`
					}

					tomlCfg := asTomlCfg{
						GasPrice: stats.SuggestedGasPrice.Int64(),
						GasTip:   stats.SuggestedGasTipCap.Int64(),
						GasFee:   stats.SuggestedGasPrice.Int64() + stats.SuggestedGasTipCap.Int64(),
					}

					marshalled, err := toml.Marshal(tomlCfg)
					if err != nil {
						return err
					}

					seth.L.Info().Msgf("Fallback prices for TOML config:\n%s", string(marshalled))

					return err
				},
			},
			{
				Name:        "trace",
				HelpName:    "trace",
				Aliases:     []string{"t"},
				Description: "trace transactions loaded from JSON file",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Aliases: []string{"f"}},
					&cli.StringFlag{Name: "txHash", Aliases: []string{"t"}},
				},
				Action: func(cCtx *cli.Context) error {
					file := cCtx.String("file")
					txHash := cCtx.String("txHash")

					if file == "" && txHash == "" {
						return fmt.Errorf("no file or transaction hash specified, use -f or -t flags")
					}

					if file != "" && txHash != "" {
						return fmt.Errorf("both file and transaction hash specified, use only one")
					}

					var transactions []string
					if file != "" {
						err := seth.OpenJsonFileAsStruct(file, &transactions)
						if err != nil {
							return err
						}
					} else {
						transactions = append(transactions, txHash)
					}

					_ = os.Setenv(seth.LogLevelEnvVar, "debug")

					cfgPath := os.Getenv(seth.CONFIG_FILE_ENV_VAR)
					if cfgPath == "" {
						return errors.New(seth.ErrEmptyConfigPath)
					}
					var cfg *seth.Config
					d, err := os.ReadFile(cfgPath)
					if err != nil {
						return errors.Wrap(err, seth.ErrReadSethConfig)
					}
					err = toml.Unmarshal(d, &cfg)
					if err != nil {
						return errors.Wrap(err, seth.ErrUnmarshalSethConfig)
					}
					absPath, err := filepath.Abs(cfgPath)
					if err != nil {
						return err
					}
					cfg.ConfigDir = filepath.Dir(absPath)

					selectedNetwork := os.Getenv(seth.NETWORK_ENV_VAR)
					if selectedNetwork != "" {
						for _, n := range cfg.Networks {
							if n.Name == selectedNetwork {
								cfg.Network = n
								break
							}
						}
						if cfg.Network == nil {
							return fmt.Errorf("network %s not defined in the TOML file", selectedNetwork)
						}

						if len(cfg.Network.URLs) == 0 {
							return errors.New("no URLs defined for the network")
						}
					} else {
						url := os.Getenv(seth.URL_ENV_VAR)

						if url == "" {
							return fmt.Errorf("network not selected, set %s=... or %s=..., check TOML config for available networks", seth.NETWORK_ENV_VAR, seth.URL_ENV_VAR)
						}

						//look for default network
						for _, n := range cfg.Networks {
							if n.Name == seth.DefaultNetworkName {
								cfg.Network = n
								cfg.Network.Name = selectedNetwork
								cfg.Network.URLs = []string{url}
								break
							}
						}

						if cfg.Network == nil {
							return errors.New("default network not defined in the TOML file")
						}

						if len(cfg.Network.URLs) == 0 {
							return errors.New("no URLs defined for the network")
						}

						if cfg.Network.DialTimeout == nil {
							cfg.Network.DialTimeout = &seth.Duration{D: seth.DefaultDialTimeout}
						}
						ctx, cancel := context.WithTimeout(context.Background(), cfg.Network.DialTimeout.Duration())
						defer cancel()
						rpcClient, err := rpc.DialOptions(ctx, cfg.MustFirstNetworkURL(), rpc.WithHeaders(cfg.RPCHeaders))
						if err != nil {
							return fmt.Errorf("failed to connect RPC client to '%s' due to: %w", cfg.MustFirstNetworkURL(), err)
						}
						client := ethclient.NewClient(rpcClient)
						defer client.Close()

						if cfg.Network.Name == seth.DefaultNetworkName {
							chainId, err := client.ChainID(context.Background())
							if err != nil {
								return errors.Wrap(err, "failed to get chain ID")
							}
							cfg.Network.ChainID = chainId.Uint64()
						}
					}

					zero := int64(0)
					cfg.EphemeralAddrs = &zero
					cfg.TracingLevel = seth.TracingLevel_All
					if cfg.Network.DialTimeout == nil {
						cfg.Network.DialTimeout = &seth.Duration{D: seth.DefaultDialTimeout}
					}

					client, err := seth.NewClientWithConfig(cfg)
					if err != nil {
						return err
					}

					seth.L.Info().Msgf("Tracing transactions from %s file", file)

					for _, txHash := range transactions {
						seth.L.Info().Msgf("Tracing transaction %s", txHash)
						ctx, cancel := context.WithTimeout(context.Background(), cfg.Network.TxnTimeout.Duration())
						tx, _, err := client.Client.TransactionByHash(ctx, common.HexToHash(txHash))
						cancel()
						if err != nil {
							return errors.Wrapf(err, "failed to get transaction %s", txHash)
						}

						_, err = client.Decode(tx, nil)
						if err != nil {
							seth.L.Info().Msgf("Possible revert reason: %s", err.Error())
						}
					}
					return err
				},
			},
		},
	}
	return app.Run(args)
}
