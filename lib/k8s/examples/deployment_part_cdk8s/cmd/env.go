package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/examples/deployment_part_cdk8s"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
)

func main() {
	e := environment.New(nil).
		AddChart(deployment_part_cdk8s.New(&deployment_part_cdk8s.Props{})).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": 2,
		}))
	if err := e.Run(); err != nil {
		panic(err)
	}
	e.Shutdown()
}
