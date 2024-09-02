package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "envresolve",
		Short: "Resolve environment variables in a string. Example envresolve '{{ env.CHAINLINK_IMAGE }}'",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Fprintln(cmd.OutOrStdout(), "Error: No input provided")
				return
			}
			input := args[0]

			fmt.Fprintln(cmd.OutOrStdout(), utils.MustResolveEnvPlaceholder(input))
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
