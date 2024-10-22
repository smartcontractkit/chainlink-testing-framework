package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	ts "github.com/smartcontractkit/chainlink-testing-framework/lib/testsummary"
)

var PrintKeyCmd = &cobra.Command{
	Use:   "print-key [key]",
	Short: "Prints all values for the given key from test summary file",
	RunE:  printKeyRunE,
}

func init() {
	PrintKeyCmd.Flags().Bool("json", true, "print as json")
	PrintKeyCmd.Flags().Bool("md", true, "print as mardkown")
}

func printKeyRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || args[0] == "" {
		return cmd.Help()
	}

	key := strings.ToLower(args[0])

	f, err := os.OpenFile(ts.SUMMARY_FILE, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer f.Close()

	fc, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var sk ts.SummaryKeys
	err = json.Unmarshal(fc, &sk)
	if err != nil {
		return err
	}

	if entry, ok := sk[key]; ok {
		if cmd.Flag("json").Value.String() == "true" {
			fmt.Println(prettyPrint(entry))
		} else if cmd.Flag("md").Value.String() == "true" {
			panic("not implemented")
		} else {
			fmt.Printf("%+v\n", entry)
		}
		return nil
	}

	return fmt.Errorf("no entry for key '%s' found", args[0])
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
