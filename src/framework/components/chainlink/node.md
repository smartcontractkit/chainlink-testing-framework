# Node

Here we provide *full* configuration parameters for `Node`

<div class="warning">
Here we provide full configuration reference, if you want to copy and run it, please remove all .out fields before!
</div>


## Configuration
```toml
[cl_node]

  [cl_node.db]
    # PostgreSQL image version and tag
    image = "postgres:12.0"
    # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
    pull_image = false

  [cl_node.node]
    # custom ports that plugins may need to expose and map to the host machine
    custom_ports = [14000, 14001]
    # A list of paths to capability binaries
    capabilities = ["./capability_1", "./capability_2"]
    # Default capabilities directory inside container
    capabilities_container_dir = "/home/capabilities"
    # Image to use, you can either provide "image" or "docker_file" + "docker_ctx" fields
    image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
    # Path to your Chainlink Dockerfile
    docker_file = "../../core/chainlink.Dockerfile"
    # Path to docker context that should be used to build from
    docker_ctx = "../.."
    # Optional name for image we build, default is "ctftmp"
    docker_image_name = "ctftmp"
    # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
    pull_image = false
    # Overrides Chainlink node TOML configuration
    # can be multiline, see example
    user_config_overrides = """
      [Log]
      level = 'info'
      """
    # Overrides Chainlink node secrets TOML configuration
    # you can only add fields, overriding existing fields is prohibited by Chainlink node
    user_secrets_overrides = """
      [AnotherSecret]
      mySecret = 'a'
      """

  # Outputs are the results of deploying a component that can be used by another component
  [cl_node.out]
    # If 'use_cache' equals 'true' we skip component setup when we run the test and return the outputs
    use_cache = true
    # Describes deployed or external Chainlink node
    [cl_node.out.node]
      # API user name
      api_auth_user = 'notreal@fakeemail.ch'
      # API password
      api_auth_password = 'fj293fbBnlQ!f9vNs'
      # Host Docker URLs the test uses
      # in case of using external component you can replace these URLs with another deployment
      p2p_url = "http://127.0.0.1:32812"
      url = "http://127.0.0.1:32847"

    # Describes deployed or external Chainlink node
    [cl_node.out.postgresql]
      # PostgreSQL connection string
      # in case of using external database can be overriden
      url = "postgresql://chainlink:thispasswordislongenough@127.0.0.1:32846/chainlink?sslmode=disable"
```

## Usage
```golang
package yourpackage_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/stretchr/testify/require"
	"testing"
)

type Step2Cfg struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	CLNode      *clnode.Input     `toml:"cl_node" validate:"required"`
}

func TestMe(t *testing.T) {
	in, err := framework.Load[Step2Cfg](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	networkCfg, err := clnode.NewNetworkCfgOneNetworkAllNodes(bc)
	require.NoError(t, err)
	in.CLNode.Node.TestConfigOverrides = networkCfg

	output, err := clnode.NewNodeWithDB(in.CLNode)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		fmt.Printf("node url: %s\n", output.Node.ExternalURL)
		require.NotEmpty(t, output.Node.ExternalURL)
	})
}
```