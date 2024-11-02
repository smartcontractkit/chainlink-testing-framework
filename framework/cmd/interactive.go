package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"sync"
)

type nodeSetForm struct {
	Network       string
	CLVersion     string
	Nodes         int
	Observability bool
	Blockscout    bool
}

func createComponentsFromForm(f *nodeSetForm) error {
	var (
		bc  *blockchain.Output
		_   *simple_node_set.Output
		err error
	)
	if err := framework.DefaultNetwork(&sync.Once{}); err != nil {
		return err
	}

	switch f.Network {
	case "anvil":
		bc, err = blockchain.NewBlockchainNetwork(&blockchain.Input{
			Type:      "anvil",
			Image:     "f4hrenh9it/foundry",
			PullImage: true,
			Port:      "8545",
			ChainID:   "31337",
		})
		if err != nil {
			return err
		}
	}
	switch f.Nodes {
	case 5:
		nspecs := make([]*clnode.Input, 0)
		for i := 0; i < f.Nodes; i++ {
			nspecs = append(nspecs, &clnode.Input{
				DbInput: &postgres.Input{
					Image:     "postgres:15.6",
					PullImage: true,
				},
				Node: &clnode.NodeInput{
					Image:     f.CLVersion,
					PullImage: true,
				},
			})
		}
		_, err = simple_node_set.NewSharedDBNodeSet(&simple_node_set.Input{
			Nodes:        5,
			OverrideMode: "all",
			NodeSpecs:    nspecs,
		}, bc, "")
	}
	return nil
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

	fmt.Println("Configuration Summary:")
	fmt.Printf("Network: %s\n", f.Network)
	fmt.Printf("Chainlink version: %d\n", f.Nodes)
	fmt.Printf("Number of Nodes: %d\n", f.Nodes)
	fmt.Printf("Observability Stack: %v\n", f.Observability)
	fmt.Printf("Blockscout Stack: %v\n", f.Blockscout)

	return createComponentsFromForm(f)
}
