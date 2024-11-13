# NodeSet

Here we provide *full* configuration parameters for `NodeSet`

<div class="warning">
Here we provide full configuration reference, if you want to copy and run it, please remove all .out fields before!
</div>

## Configuration

This component requires some Blockchain to be deployed, add this to config
```toml
[blockchain_a]
  # Blockchain node type, can be "anvil" or "geth"
  type = "anvil"
  # Chain ID
  chain_id = "31337"
  # Anvil command line params, ex.: docker_cmd_params = ['--block-time=1', '...']
  docker_cmd_params = []
  # Docker image and tag
  image = "f4hrenh9it/foundry:latest"
  # External port to expose
  port = "8545"
  # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
  pull_image = false

  # Outputs are the results of deploying a component that can be used by another component
  [blockchain_a.out]
    chain_id = "31337"
    # If 'use_cache' equals 'true' we skip component setup when we run the test and return the outputs
    use_cache = true

    [[blockchain_a.out.nodes]]
      # URLs to access the node(s) inside docker network, used by other components
      docker_internal_http_url = "http://anvil-14411:8545"
      docker_internal_ws_url = "ws://anvil-14411:8545"
      # URLs to access the node(s) on your host machine or in CI
      http_url = "http://127.0.0.1:33955"
      ws_url = "ws://127.0.0.1:33955"
```

Then configure NodeSet
```toml
[nodeset]
  # amount of Chainlink nodes to spin up
  nodes = 5
  # Override mode: can be "all" or "each"
  # defines how we override configs, either we apply first node fields to all of them
  # or we define each node custom configuration (used in compatibility testing)
  override_mode = "all"
  # HTTP API port range start, each new node get port incremented (host machine)
  http_port_range_start = 10000
  # P2P API port range start, each new node get port incremented (host machine)
  p2p_port_range_start = 12000


  [[nodeset.node_specs]]
    # Optional URL for fake data provider URL
    # usually set up in test with local mock server
    data_provider_url = "http://example.com"

    [nodeset.node_specs.db]
      # PostgreSQL image version and tag
      image = "postgres:15.6"
      # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
      pull_image = true
      # PostgreSQL volume name
      volume_name = ""

    [nodeset.node_specs.node]
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
      pull_image = true
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
  [nodeset.out]
    # If 'use_cache' equals 'true' we skip component setup when we run the test and return the outputs
    use_cache = true
    
    # Describes deployed or external Chainlink nodes
    [[nodeset.out.cl_nodes]]
      use_cache = true

      # Describes deployed or external Chainlink node
      [nodeset.out.cl_nodes.node]
        # API user name
        api_auth_user = 'notreal@fakeemail.ch'
        # API password
        api_auth_password = 'fj293fbBnlQ!f9vNs'
        # Host Docker URLs the test uses
        # in case of using external component you can replace these URLs with another deployment
        p2p_url = "http://127.0.0.1:32996"
        url = "http://127.0.0.1:33096"
      # Describes PostgreSQL instance
      [nodeset.out.cl_nodes.postgresql]
        # PostgreSQL connection string
        # in case of using external database can be overriden
        url = "postgresql://chainlink:thispasswordislongenough@127.0.0.1:33094/chainlink?sslmode=disable"
    
    # Can have more than one node, fields are the same, see above ^^
    [[nodeset.out.cl_nodes]]
      [nodeset.out.cl_nodes.node]
      [nodeset.out.cl_nodes.postgresql]
    ...
```

## Usage
```golang
package capabilities_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestMe(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.HostURL)
			require.NotEmpty(t, n.Node.HostP2PURL)
		}
	})
}
```