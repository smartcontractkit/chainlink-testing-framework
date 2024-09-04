package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/cmd/internal"
)

var rootCmd = &cobra.Command{
	Use:   "ctf-utils",
	Short: "CTF Utils Tool",
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	internal.GetLatestEthereumClientVersionCmd.SetContext(ctx)

	rootCmd.AddCommand(internal.GetLatestEthereumClientVersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
