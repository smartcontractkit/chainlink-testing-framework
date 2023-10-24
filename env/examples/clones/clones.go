package main

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
)

func main() {
	// Multiple environments of the same type/chart
	err := environment.New(&environment.Config{
		Labels:            []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5)},
		KeepConnection:    true,
		RemoveOnInterrupt: true,
	}).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu": "344m",
					},
					"limits": map[string]interface{}{
						"cpu": "344m",
					},
				},
			},
			"db": map[string]interface{}{
				"stateful": "true",
				"capacity": "1Gi",
			},
		})).
		AddHelm(chainlink.New(1,
			map[string]interface{}{
				"chainlink": map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"cpu": "577m",
						},
						"limits": map[string]interface{}{
							"cpu": "577m",
						},
					},
				},
			})).
		Run()
	if err != nil {
		panic(err)
	}
}
