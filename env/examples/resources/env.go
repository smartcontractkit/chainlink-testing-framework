package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/env/pkg"
	"github.com/smartcontractkit/chainlink-testing-framework/env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/env/pkg/helm/ethereum"
)

func main() {
	e := environment.New(&environment.Config{
		Labels: []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5)},
	}).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, nil))
	err := e.Run()
	if err != nil {
		panic(err)
	}
	// default k8s selector
	summ, err := e.ResourcesSummary("app in (chainlink-0, geth)")
	if err != nil {
		panic(err)
	}
	log.Warn().Interface("Resources", summ).Send()
	e.Shutdown()
}
