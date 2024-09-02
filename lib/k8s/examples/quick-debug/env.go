package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
)

func main() {
	err := environment.New(&environment.Config{
		Labels:            []string{"type=construction-in-progress"},
		NamespacePrefix:   "new-environment",
		KeepConnection:    true,
		RemoveOnInterrupt: true,
	}).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, nil)).
		Run()
	if err != nil {
		panic(err)
	}
}
