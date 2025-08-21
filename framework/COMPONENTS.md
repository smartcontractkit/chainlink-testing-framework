## Developing Components

To build a scalable framework that enables the reuse of our product deployments (contracts or services in Docker), we need to establish a clear component structure.
```
package mycomponent

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Input struct {
    // inputs fields that component exposes for configuration
    ...
    // outputs are embedded into inputs so framework can automatically save them
	Out                      *Output  `toml:"out"`
}

type Output struct {
    UseCache bool             `toml:"use_cache"`
    // outputs that will be dumped to config and cached
}


func NewComponent(input *Input) (*Output, error) {
	if input.Out.UseCache {
		return input.Out, nil
	}
	
	// component logic here
	// deploy a docker container(s)
	// or deploy a set of smart contracts
	
	input.Out = &Output{
	    UseCache: true,
	    // other fields
	    ...
	}
	return out, nil
}
```
Each component can define inputs and outputs, following these rules:

- Outputs should be included within inputs.
- If your component is used for side effects output can be omitted.
- `input.Out.UseCache` should be added if you'd like to use caching, see more [here](CACHING.md)

### Docker components good practices for [testcontainers-go](https://golang.testcontainers.org/):

An example [simple component](components/blockchain/anvil.go)

An example of [complex component](components/clnode/clnode.go)

An example of [composite component](components/node_set_extended/don.go)

- Inputs should include at least `image`, `tag` and `pull_image` field
```
	Image                string `toml:"image" validate:"required"`
	Tag                  string `toml:"tag" validate:"required"`
	PullImage            bool   `toml:"pull_image" validate:"required"`
```

- `ContainerRequest` must contain labels, network and alias required for local observability stack and deployment isolation
```
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
```
- In order to copy files into container use `framework.WriteTmpFile(data string, fileName string)`
```
	userSecretsOverridesFile, err := WriteTmpFile(in.Node.UserSecretsOverrides, "user-secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
```
- Output of docker component must contain all the URLs component exposes for access, both for internal docker usage and external test (host) usage
```
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}

	return &NodeOut{
	    UseCache: true,
		InternalURL: fmt.Sprintf("http://%s:%s", containerName, in.Node.Port),
		ExternalURL:   fmt.Sprintf("http://%s:%s", host, mp.Port()),
	}, nil
```
