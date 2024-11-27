package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type nodeSetForm struct {
	Network          string
	CLVersion        string
	Nodes            int
	Observability    bool
	Blockscout       bool
	BlockscoutRPCURL string
}

func createComponentsFromForm(form *nodeSetForm) error {
	var (
		bc  *blockchain.Output
		_   *simple_node_set.Output
		err error
	)
	if err := framework.DefaultNetwork(&sync.Once{}); err != nil {
		return err
	}
	switch form.Network {
	case "anvil":
		f := func() {
			bc, err = blockchain.NewBlockchainNetwork(&blockchain.Input{
				Type:    "anvil",
				Image:   "f4hrenh9it/foundry",
				Port:    "8545",
				ChainID: "31337",
			})
		}
		err = spinner.New().
			Title("Creating anvil blockchain..").
			Action(f).
			Run()
	}
	switch form.Nodes {
	case 5:
		nspecs := make([]*clnode.Input, 0)
		for i := 0; i < form.Nodes; i++ {
			nspecs = append(nspecs, &clnode.Input{
				Node: &clnode.NodeInput{
					Image: form.CLVersion,
				},
			})
		}
		f := func() {
			_, err = simple_node_set.NewSharedDBNodeSet(&simple_node_set.Input{
				Nodes:        5,
				OverrideMode: "all",
				DbInput: &postgres.Input{
					Image: "postgres:12.0",
				},
				NodeSpecs: nspecs,
			}, bc)
		}
		err = spinner.New().
			Title("Creating node set..").
			Action(f).
			Run()
	}
	switch form.Observability {
	case true:
		if err = framework.NewPromtail(); err != nil {
			return err
		}
		if err := observabilityUp(); err != nil {
			return err
		}
	}
	switch form.Blockscout {
	case true:
		if err := blockscoutUp(form.BlockscoutRPCURL); err != nil {
			return err
		}
	}
	return err
}

func cleanup(form *nodeSetForm) error {
	var err error
	f := func() {
		err = removeTestContainers()
	}
	err = spinner.New().
		Title("Removing docker resources..").
		Action(f).
		Run()
	switch form.Observability {
	case true:
		if err := observabilityDown(); err != nil {
			return err
		}
	}
	switch form.Blockscout {
	case true:
		if err := blockscoutDown(form.BlockscoutRPCURL); err != nil {
			return err
		}
	}
	return err
}

func runSetupForm() error {
	if !framework.IsDockerRunning() {
		return fmt.Errorf(`Docker daemon is not running!
Please set up OrbStack (https://orbstack.dev/)
or
Docker Desktop (https://www.docker.com/products/docker-desktop/)
`)
	}
	f := &nodeSetForm{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which network would you like to connect to?").
				Options(
					huh.NewOption("Anvil", "anvil"),
				).
				Value(&f.Network),

			huh.NewSelect[int]().
				Title("How many nodes do you need?").
				Options(
					huh.NewOption("5", 5),
				).
				Value(&f.Nodes),

			huh.NewSelect[string]().
				Title("Choose Chainlink node version").
				Options(
					huh.NewOption("public.ecr.aws/chainlink/chainlink:v2.17.0-arm64", "public.ecr.aws/chainlink/chainlink:v2.17.0-arm64"),
					huh.NewOption("public.ecr.aws/chainlink/chainlink:v2.17.0", "public.ecr.aws/chainlink/chainlink:v2.17.0"),
					huh.NewOption("public.ecr.aws/chainlink/chainlink:v2.16.0-arm64", "public.ecr.aws/chainlink/chainlink:v2.16.0-arm64"),
					huh.NewOption("public.ecr.aws/chainlink/chainlink:v2.16.0", "public.ecr.aws/chainlink/chainlink:v2.16.0"),
				).
				Value(&f.CLVersion),
			huh.NewConfirm().
				Title("Do you need to spin up an observability stack?").
				Value(&f.Observability),

			huh.NewConfirm().
				Title("Do you need to spin up a Blockscout stack?").
				Value(&f.Blockscout),
			huh.NewSelect[string]().
				Title("To which blockchain node you want Blockscout to connect?").
				Options(
					huh.NewOption("Network 1", "http://host.docker.internal:8545"),
					huh.NewOption("Network 2", "http://host.docker.internal:8550"),
				).
				Value(&f.BlockscoutRPCURL),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("failed to run form: %w", err)
	}
	if err := createComponentsFromForm(f); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		err := cleanup(f)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	framework.L.Info().Msg("Services are up! Press Ctrl+C to remove them..")
	select {}
}
