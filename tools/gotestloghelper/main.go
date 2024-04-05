package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper/gotestevent"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/clihelper"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		stop() // restore default exit behavior
		log.Println("Cancelling... interrupt again to exit")
	}()

	config := gotestevent.NewDefaultConfig()
	config.ShouldImmediatelyPrint = true
	config.RemoveTLogPrefix = flag.Bool("tlogprefix", false, "Set to true to remove the go test log prefix")
	config.IsJsonInput = flag.Bool("json", false, "Set to true to enable parsing the input from a go test -json output")
	flag.Var(config.OnlyErrors, "onlyerrors", "Set to true to only print tests that failed, not compatible without -json")
	config.Color = flag.Bool("color", false, "Set to true to enable color output")
	config.CI = flag.Bool("ci", false, "Set to true to enable CI mode, which will print out the logs with groupings when combined with -json")
	config.SinglePackage = flag.Bool("singlepackage", false, "Set to true if the go test output is from a single package only, this will print tests out as they finish instead of waiting for the package to finish")
	config.ErrorAtTopLength = flag.Int("errorattoplength", 100, "If the error message doesn't appear before this many lines, it will be printed at the top of the test output as well. Set to 0 to disable. Only works with -ci")
	flag.Parse()

	err := config.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v\n", err)
	}

	// Add modifiers to the list based on the flags provided, order could be important
	modifiers := gotestevent.SetupModifiers(config)

	err = ReadAndModifyLogs(ctx, os.Stdin, modifiers, config)
	if err != nil {
		log.Fatalf("Error reading and modifying logs: %v\n", err)
	}
	if config.FailuresExist {
		os.Exit(1)
	}
}

func ReadAndModifyLogs(ctx context.Context, r io.Reader, modifiers []gotestevent.TestLogModifier, c *gotestevent.TestLogModifierConfig) error {
	return clihelper.ReadLine(ctx, r, func(b []byte) error {
		var err error
		te := &gotestevent.GoTestEvent{}

		// build a TestEvent from the input line
		if *c.IsJsonInput {
			te, err = gotestevent.ParseTestEvent(b)
			if err != nil {
				log.Fatalf("Error parsing json test event from stdin: %v\n", err)
			}
			if te == nil {
				// got a non json line when expecting json, just print it out and move on
				fmt.Println(string(b))
				return nil
			}
		} else {
			te.Action = gotestevent.ActionOutput
			te.Output = string(b)
		}

		// Run the modifiers on the output
		for _, m := range modifiers {
			err := m(te, c)
			if err != nil {
				log.Fatalf("Error modifying output: %v\nProblematic line: %s\n", err, te.Output)
			}
		}

		// print line back out
		if c.ShouldImmediatelyPrint {
			if *c.IsJsonInput {
				s, err := te.String()
				if err != nil {
					return err
				}
				fmt.Println(s)
			} else {
				fmt.Println(te.Output)
			}
		}
		return nil
	})
}
