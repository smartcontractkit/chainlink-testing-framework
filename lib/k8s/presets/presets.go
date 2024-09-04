package presets

import (
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/foundry"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/reorg"
)

var BaseToml = `[Log]
Level = "debug"
JSONConsole = true
[Log.File]
MaxSize = "0b"
[WebServer]
AllowOrigins = "*"
HTTPPort = 6688
SecureCookies = false
SessionTimeout = "999h0m0s"
[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 100
[WebServer.TLS]
HTTPSPort = 0
[Database]
MaxIdleConns = 20
MaxOpenConns = 40
MigrateOnStartup = true
[OCR2]
Enabled = true
[P2P]
[P2P.V2]
ListenAddresses = ["0.0.0.0:6690"]`

func OnlyRemoteRunner(config *environment.Config) *environment.Environment {
	return environment.New(config)
}

// EVMOneNode local development Chainlink deployment
func EVMOneNode(config *environment.Config) (*environment.Environment, error) {

	c := chainlink.New(0, map[string]any{
		"replicas": 1,
		"toml":     BaseToml,
	})

	return environment.New(config).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(c), nil
}

// EVMMinimalLocalBS local development Chainlink deployment,
// 1 bootstrap + 4 oracles (minimal requirements for OCR) + Blockscout
func EVMMinimalLocalBS(config *environment.Config) (*environment.Environment, error) {
	c := chainlink.New(0, map[string]any{
		"replicas": 5,
		"toml":     BaseToml,
	})
	return environment.New(config).
		AddChart(blockscout.New(&blockscout.Props{})).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(c), nil
}

// EVMMinimalLocal local development Chainlink deployment,
// 1 bootstrap + 4 oracles (minimal requirements for OCR)
func EVMMinimalLocal(config *environment.Config) *environment.Environment {
	return environment.New(config).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": 5,
			"toml":     BaseToml,
		}))
}

// EVMMinimalLocal local development Chainlink deployment,
// 1 bootstrap + 4 oracles (minimal requirements for OCR)
func EVMMultipleNodesWithDiffDBVersion(config *environment.Config) *environment.Environment {
	return environment.New(config).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml": BaseToml,
			"nodes": []map[string]any{
				{
					"name": "node-1",
					"db": map[string]any{
						"image": map[string]any{
							"image":   "postgres",
							"version": "13.12",
						},
					},
				},
				{
					"name": "node-2",
					"db": map[string]any{
						"image": map[string]any{
							"image":   "postgres",
							"version": "14.9",
						},
					},
				},
				{
					"name": "node-3",
					"db": map[string]any{
						"image": map[string]any{
							"image":   "postgres",
							"version": "15.4",
						},
					},
				},
			},
		}))
}

// EVMReorg deployment for two Ethereum networks re-org test
func EVMReorg(config *environment.Config) (*environment.Environment, error) {
	var clToml = `[[EVM]]
ChainID = '1337'
FinalityDepth = 200

[[EVM.Nodes]]
Name = 'geth'
WSURL = 'ws://geth-ethereum-geth:8546'
HTTPURL = 'http://geth-ethereum-geth:8544'

[EVM.HeadTracker]
HistoryDepth = 400
[OCR2]
Enabled = true
[P2P]
[P2P.V2]
ListenAddresses = ["0.0.0.0:6690"]`
	c := chainlink.New(0, map[string]interface{}{
		"replicas": 5,
		"toml":     clToml,
	})

	return environment.New(config).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(reorg.New(&reorg.Props{
			NetworkName: "geth",
			NetworkType: "geth-reorg",
			Values: map[string]interface{}{
				"geth": map[string]interface{}{
					"genesis": map[string]interface{}{
						"networkId": "1337",
					},
				},
			},
		})).
		AddHelm(reorg.New(&reorg.Props{
			NetworkName: "geth-2",
			NetworkType: "geth-reorg",
			Values: map[string]interface{}{
				"geth": map[string]interface{}{
					"genesis": map[string]interface{}{
						"networkId": "2337",
					},
				},
			},
		})).
		AddHelm(c), nil
}

// EVMSoak deployment for a long running soak tests
func EVMSoak(config *environment.Config) *environment.Environment {
	return environment.New(config).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			Simulated: true,
			Values: map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "1000m",
						"memory": "2048Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "1000m",
						"memory": "2048Mi",
					},
				},
			},
		})).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": 5,
			"toml":     BaseToml,
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "1Gi",
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "250m",
						"memory": "256Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "250m",
						"memory": "256Mi",
					},
				},
			},
			"chainlink": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "1000m",
						"memory": "2048Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "1000m",
						"memory": "2048Mi",
					},
				},
			},
		}))
}

func FoundryNetwork(config *environment.Config) *environment.Environment {
	return environment.New(config).
		AddHelm(foundry.New(&foundry.Props{
			Values: map[string]interface{}{
				"fullnameOverride": "foundry",
			},
		}))
}
