# NodeSet

## Configuration
```toml
[nodeset]
  # amount of Chainlink nodes to spin up
  nodes = 5
  # Override mode: can be "all" or "each"
  # defines how we override configs, either we apply first node fields to all of them
  # or we define each node custom configuration (used in compatibility testing)
  override_mode = "all"

  [[nodeset.node_specs]]
    # Optional URL for fake data provider URL
    # usually set up in test with local mock server
    data_provider_url = ""

    [nodeset.node_specs.db]
      # PostgreSQL image version and tag
      image = "postgres:15.6"
      # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
      pull_image = true

    [nodeset.node_specs.node]
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
package yourpackage_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestNodeSet(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	...
}
```