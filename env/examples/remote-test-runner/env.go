package main

import (
	"os"

	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
)

func main() {
	// see REMOTE_RUN.md for tutorial
	e := environment.New(&environment.Config{
		NamespacePrefix: "zmytest",
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": 1,
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   os.Getenv(config.EnvVarCLImage),
					"version": os.Getenv(config.EnvVarCLTag),
				},
			},
		}))
	if err := e.Run(); err != nil {
		panic(err)
	}
}
