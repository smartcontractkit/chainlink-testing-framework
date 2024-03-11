package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/cmd/internal"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-test-runner",
	Short: "K8s Test Runner",
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	internal.Create.SetContext(ctx)

	rootCmd.AddCommand(internal.Create)
	rootCmd.AddCommand(internal.Run)
	rootCmd.AddCommand(internal.ECR)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
