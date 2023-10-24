# How to create environments

- [Getting started](#getting-started)
- [Connect to environment](#connect-to-environment)
- [Creating environments](#creating-environments)
  - [Debugging a new integration environment](#debugging-a-new-integration-environment)
  - [Creating a new deployment part in Helm](#creating-a-new-deployment-part-in-helm)
  - [Creating a new deployment part in cdk8s](#creating-a-new-deployment-part-in-cdk8s)
  - [Using multi-stage environment](#using-multi-stage-environment)
- [Modifying environments](#modifying-environments)
  - [Modifying environment from code](#modifying-environment-from-code)
  - [Modifying environment part from code](#modifying-environment-part-from-code)
- [Configuring](#configuring)
  - [Environment variables](#environment-variables)
  - [Environment config](#environment-config)
- [Utilities](#utilities)
  - [Collecting logs](#collecting-logs)
  - [Resources summary](#resources-summary)
- [Chaos](#chaos)
- [Coverage](#coverage)
- [Remote run](./REMOTE_RUN.md)

## Getting started

Read [here](KUBERNETES.md) about how to spin up a local cluster if you don't have one

Let's create a simple environment by combining different deployment parts

Create `examples/simple/env.go`

```golang
package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
)

func main() {
	err := environment.New(&environment.Config{
      KeepConnection:    false,
      RemoveOnInterrupt: false,
    }).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, nil)).
		Run()
	if err != nil {
		panic(err)
	}
}
```

Then run `go run examples/simple/env.go`

Now you have your environment running, you can [connect](#connect-to-environment) to it later

## Connect to environment

We've already created an environment [previously](#getting-started), now we can connect

If you are planning to use environment locally not in tests and keep connection, modify `KeepConnection` in `environment.Config` we used

```
      KeepConnection:    true,
```

Add `ENV_NAMESPACE=${your_env_namespace}` var and run `go run examples/simple/env.go` again

You can get the namespace name from logs on creation time

# Creating environments

## Debugging a new integration environment

You can spin up environment and block on forwarder if you'd like to run some other code

```golang
package main

import (
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
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
```

Send any signal to remove the namespace then, for example `Ctrl+C` `SIGINT`

## Creating a new deployment part in Helm

Let's add a new [deployment part](examples/deployment_part/sol.go), it should implement an interface

```golang
// ConnectedChart interface to interact both with cdk8s apps and helm charts
type ConnectedChart interface {
	// IsDeploymentNeeded
	// true - we deploy/connect and expose environment data
	// false - we are using external environment, but still exposing data
	IsDeploymentNeeded() bool
	// GetName name of the deployed part
	GetName() string
	// GetPath get Helm chart path, repo or local path
	GetPath() string
	// GetProps get code props if it's typed environment
	GetProps() any
	// GetValues get values.yml props as map, if it's Helm
	GetValues() *map[string]any
	// ExportData export deployment part data in the env
	ExportData(e *Environment) error
}
```

When creating new deployment part, you can use any public Helm chart or a local path in Helm props

```golang
func New(props *Props) environment.ConnectedChart {
	if props == nil {
		props = defaultProps()
	}
	return Chart{
		HelmProps: &HelmProps{
			Name:   "sol",
			Path:   "chainlink-qa/solana-validator", // ./local_path/chartdir will work too
			Values: &props.Values,
		},
		Props: props,
	}
}
```

Now let's tie them together

```golang
package main

import (
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/examples/deployment_part"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"time"
)

func main() {
	e := environment.New(&environment.Config{
		NamespacePrefix:   "adding-new-deployment-part",
		TTL:               3 * time.Hour,
		KeepConnection:    true,
		RemoveOnInterrupt: true,
	}).
		AddHelm(deployment_part.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 5,
			"env": map[string]any{
				"SOLANA_ENABLED":              "true",
				"EVM_ENABLED":                 "false",
				"EVM_RPC_ENABLED":             "false",
				"CHAINLINK_DEV":               "false",
				"FEATURE_OFFCHAIN_REPORTING2": "true",
				"feature_offchain_reporting":  "false",
				"P2P_NETWORKING_STACK":        "V2",
				"P2PV2_LISTEN_ADDRESSES":      "0.0.0.0:6690",
				"P2PV2_DELTA_DIAL":            "5s",
				"P2PV2_DELTA_RECONCILE":       "5s",
				"p2p_listen_port":             "0",
			},
		}))
	if err := e.Run(); err != nil {
		panic(err)
	}
}
```

Then run it `examples/deployment_part/cmd/env.go`

## Creating a new deployment part in cdk8s

Let's add a new [deployment part](examples/deployment_part/sol.go), it should implement the same interface

```golang
// ConnectedChart interface to interact both with cdk8s apps and helm charts
type ConnectedChart interface {
	// IsDeploymentNeeded
	// true - we deploy/connect and expose environment data
	// false - we are using external environment, but still exposing data
	IsDeploymentNeeded() bool
	// GetName name of the deployed part
	GetName() string
	// GetPath get Helm chart path, repo or local path
	GetPath() string
	// GetProps get code props if it's typed environment
	GetProps() any
	// GetValues get values.yml props as map, if it's Helm
	GetValues() *map[string]any
	// ExportData export deployment part data in the env
	ExportData(e *Environment) error
}
```
Now let's tie them together

```golang
package main

import (
  "github.com/smartcontractkit/chainlink-env/environment"
  "github.com/smartcontractkit/chainlink-env/examples/deployment_part_cdk8s"
  "github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
  "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
)

func main() {
  e := environment.New(nil).
    AddChart(deployment_part_cdk8s.New(&deployment_part_cdk8s.Props{})).
    AddHelm(ethereum.New(nil)).
          AddHelm(chainlink.New(0, map[string]any{
            "replicas": 2,
          }))
  if err := e.Run(); err != nil {
    panic(err)
  }
  e.Shutdown()
}
```

Then run it `examples/deployment_part_cdk8s/cmd/env.go`

## Using multi-stage environment

You can split [environment](examples/multistage/env.go) deployment in several parts if you need to first copy something into a pod or use connected clients first

```golang
package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
)

func main() {
	e := environment.New(nil)
	err := e.
		AddChart(blockscout.New(&blockscout.Props{})). // you can also add cdk8s charts if you like Go code
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, nil)).
		Run()
	if err != nil {
		panic(err)
	}
	// do some other stuff with deployed charts
	pl, err := e.Client.ListPods(e.Cfg.Namespace, "app=chainlink-0")
	if err != nil {
		panic(err)
	}
	dstPath := fmt.Sprintf("%s/%s:/", e.Cfg.Namespace, pl.Items[0].Name)
	if _, _, _, err = e.Client.CopyToPod(e.Cfg.Namespace, "./examples/multistage/someData.txt", dstPath, "node"); err != nil {
		panic(err)
	}
	// deploy another part
	err = e.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		Run()
	defer func() {
		errr := e.Shutdown()
		panic(errr)
	}()
	if err != nil {
		panic(err)
	}
}
```

# Modifying environments

## Modifying environment from code

In case you need to [modify](examples/modify_cdk8s/env.go) environment in tests you can always construct manifest again and apply it

That's working for `cdk8s` components too

```golang
package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
)

func main() {
	e := environment.New(&environment.Config{
		NamespacePrefix: "modified-env",
		Labels:          []string{fmt.Sprintf("envType=Modified")},
	}).
		AddChart(blockscout.New(&blockscout.Props{
			WsURL:   "ws://geth:8546",
			HttpURL: "http://geth:8544",
		})).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
		}))
	err := e.Run()
	if err != nil {
		panic(err)
	}
	e.ClearCharts()
	err = e.
		AddChart(blockscout.New(&blockscout.Props{
			HttpURL: "http://geth:9000",
		})).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
		})).
		Run()
	if err != nil {
		panic(err)
	}
}
```

## Modifying environment part from code

We can [modify](examples/modify_helm/env.go) only a part of environment

```golang
package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
)

func main() {
	e := environment.New(&environment.Config{
		NamespacePrefix: "modified-env",
		Labels:          []string{fmt.Sprintf("envType=Modified")},
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
		}))
	err := e.Run()
	if err != nil {
		panic(err)
	}
	e.Cfg.KeepConnection = true
	e.Cfg.RemoveOnInterrupt = true
	e, err = e.
		ReplaceHelm("chainlink-0", chainlink.New(0, map[string]any{
			"replicas": 2,
		}))
	if err != nil {
		panic(err)
	}
	err = e.Run()
	if err != nil {
		panic(err)
	}
}
```

# Configuring

## Environment variables

List of environment variables available

```golang
const (
	EnvVarNamespace            = "ENV_NAMESPACE"
	EnvVarNamespaceDescription = "Namespace name to connect to"
	EnvVarNamespaceExample     = "chainlink-test-epic"

	EnvVarCLImage            = "CHAINLINK_IMAGE"
	EnvVarCLImageDescription = "Chainlink image repository"
	EnvVarCLImageExample     = "public.ecr.aws/chainlink/chainlink"

	EnvVarCLTag            = "CHAINLINK_VERSION"
	EnvVarCLTagDescription = "Chainlink image tag"
	EnvVarCLTagExample     = "1.5.1-root"

	EnvVarUser            = "CHAINLINK_ENV_USER"
	EnvVarUserDescription = "Owner of an environment"
	EnvVarUserExample     = "Satoshi"

	EnvVarCLCommitSha            = "CHAINLINK_COMMIT_SHA"
	EnvVarCLCommitShaDescription = "The sha of the commit that you're running tests on. Mostly used for CI"
	EnvVarCLCommitShaExample     = "${{ github.sha }}"

	EnvVarTestTrigger            = "TEST_TRIGGERED_BY"
	EnvVarTestTriggerDescription = "How the test was triggered, either manual or CI."
	EnvVarTestTriggerExample     = "CI"

	EnvVarLogLevel            = "TEST_LOG_LEVEL"
	EnvVarLogLevelDescription = "Environment logging level"
	EnvVarLogLevelExample     = "info | debug | trace"

	EnvVarSlackKey            = "SLACK_API_KEY"
	EnvVarSlackKeyDescription = "The OAuth Slack API key to report tests results with"
	EnvVarSlackKeyExample     = "xoxb-example-key"

	EnvVarSlackChannel            = "SLACK_CHANNEL"
	EnvVarSlackChannelDescription = "The Slack code for the channel you want to send the notification to"
	EnvVarSlackChannelExample     = "C000000000"

	EnvVarSlackUser            = "SLACK_USER"
	EnvVarSlackUserDescription = "The Slack code for the user you want to notify"
	EnvVarSlackUserExample     = "U000000000"
)
```
### Environment config

```golang
// Config is an environment common configuration, labels, annotations, connection types, readiness check, etc.
type Config struct {
	// TTL is time to live for the environment, used with kube-janitor
	TTL time.Duration
	// NamespacePrefix is a static namespace prefix
	NamespacePrefix string
	// Namespace is full namespace name
	Namespace string
	// Labels is a set of labels applied to the namespace in a format of "key=value"
	Labels            []string
	nsLabels          *map[string]*string
	// ReadyCheckData is settings for readiness probes checks for all deployment components
	// checking that all pods are ready by default with 8 minutes timeout
	//	&client.ReadyCheckData{
	//		ReadinessProbeCheckSelector: "",
	//		Timeout:                     15 * time.Minute,
	//	}
	ReadyCheckData    *client.ReadyCheckData
	// DryRun if true, app will just generate a manifest in local dir
	DryRun            bool
	// InsideK8s used for long-running soak tests where you connect to env from the inside
	InsideK8s         bool
	// KeepConnection keeps connection until interrupted with a signal, useful when prototyping and debugging a new env
	KeepConnection    bool
	// RemoveOnInterrupt automatically removes an environment on interrupt
	RemoveOnInterrupt bool
}
```

# Utilities

## Collecting logs

You can collect the [logs](examples/dump/env.go) while running tests, or if you have created an enrionment [already](#connect-to-environment)

```golang
package main

import (
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
)

func main() {
	e := environment.New(nil).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, nil))
	if err := e.Run(); err != nil {
		panic(err)
	}
	if err := e.DumpLogs("logs/mytest"); err != nil {
		panic(err)
	}
}
```

## Resources summary

It can be useful to get current env [resources](examples/resources/env.go) summary for test reporting

```golang
package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
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
```

# Chaos

Check our [tests](https://github.com/smartcontractkit/chainlink/blob/develop/integration-tests/chaos/chaos_test.go) to see how we using Chaosmesh

# Coverage

Build your target image with those 2 steps to allow automatic coverage discovery

```Dockerfile
FROM ...

# add those 2 steps to instrument the code
RUN curl -s https://api.github.com/repos/qiniu/goc/releases/latest | grep "browser_download_url.*-linux-amd64.tar.gz" | cut -d : -f 2,3 | tr -d \" | xargs -n 1 curl -L | tar -zx && chmod +x goc && mv goc /usr/local/bin
# -o my_service means service will be called "my_service" in goc coverage service
# --center http://goc:7777 means that on deploy, your instrumented service will automatically register to a local goc node inside your deployment (namespace)
RUN goc build -o my_service . --center http://goc:7777

CMD ["./my_service"]
```
Add `goc` to your deployment, check example with `dummy` service deployment:

```golang
package main

import (
  "time"

  "github.com/smartcontractkit/chainlink-env/environment"
  goc "github.com/smartcontractkit/chainlink-env/pkg/cdk8s/goc"
  dummy "github.com/smartcontractkit/chainlink-env/pkg/cdk8s/http_dummy"
)

func main() {
  e := environment.New(nil).
    AddChart(goc.New()).
    AddChart(dummy.New())
  if err := e.Run(); err != nil {
    panic(err)
  }
  // run your test logic here
  time.Sleep(1 * time.Minute)
  if err := e.SaveCoverage(); err != nil {
    panic(err)
  }
  // clear the coverage, rerun the tests again if needed
  if err := e.ClearCoverage(); err != nil {
    panic(err)
  }
}

```

After tests are finished, coverage is collected for every service, check `cover` directory
