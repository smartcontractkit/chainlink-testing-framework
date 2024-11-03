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
	Network       string
	CLVersion     string
	Nodes         int
	Observability bool
	Blockscout    bool
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
				Type:      "anvil",
				Image:     "f4hrenh9it/foundry",
				PullImage: true,
				Port:      "8545",
				ChainID:   "31337",
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
				DbInput: &postgres.Input{
					Image:     "postgres:15.6",
					PullImage: true,
				},
				Node: &clnode.NodeInput{
					Image:     form.CLVersion,
					PullImage: true,
				},
			})
		}
		f := func() {
			_, err = simple_node_set.NewSharedDBNodeSet(&simple_node_set.Input{
				Nodes:        5,
				OverrideMode: "all",
				NodeSpecs:    nspecs,
			}, bc, "")
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
		if err := blockscoutUp(); err != nil {
			return err
		}
	}
	return err
}

func cleanup(form *nodeSetForm) error {
	var err error
	f := func() {
		err = cleanDockerResources()
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
		if err := blockscoutDown(); err != nil {
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
				Value(&f.Network), // stores the selected network

			huh.NewSelect[int]().
				Title("How many nodes do you need?").
				Options(
					huh.NewOption("5", 5),
				).
				Value(&f.Nodes), // stores the selected number of nodes

			huh.NewSelect[string]().
				Title("Choose Chainlink node version").
				Options(
					huh.NewOption("public.ecr.aws/chainlink/chainlink:v2.17.0", "public.ecr.aws/chainlink/chainlink:v2.17.0")).
				Value(&f.CLVersion),
			huh.NewConfirm().
				Title("Do you need to spin up an observability stack?").
				Value(&f.Observability), // stores the observability option

			huh.NewConfirm().
				Title("Do you need to spin up a Blockscout stack?").
				Value(&f.Blockscout), // stores the Blockscout option
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
		fmt.Println("\nReceived Ctrl+C, starting custom cleanup...")
		err := cleanup(f)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	framework.L.Info().Msg("Press Ctrl+C to remove the stack..")
	select {}
}
