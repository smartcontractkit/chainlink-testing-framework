package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
)

var GetLatestEthereumClientVersionCmd = &cobra.Command{
	Use:   "get-latest-ethereum-client-version",
	Short: "Get the latest Ethereum client release version from Github",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("please provide the repository in the format 'org/repo:tag'")
		}

		repo := args[0]

		latest, err := test_env.FetchLatestEthereumClientDockerImageVersionIfNeed(repo)
		if err != nil {
			return fmt.Errorf("error fetching release information: %v", err)
		}

		fmt.Println(strings.Split(latest, ":")[1])
		return nil
	},
}
