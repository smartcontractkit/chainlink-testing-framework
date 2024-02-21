package gotestevent

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/clitext"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/github"
)

const testingLogPrefix = `^(\s+)(\w+\.go:\d+: )`

// func debugLog(msg string) {
// 	fmt.Println(clitext.Color(clitext.ColorRed, "TATATATAATATATAT"+msg))
// }

type Action string

const (
	ActionRun    Action = "run"
	ActionPass   Action = "pass"
	ActionFail   Action = "fail"
	ActionOutput Action = "output"
	ActionSkip   Action = "skip"
	ActionStart  Action = "start"
	ActionPause  Action = "pause"
	ActionCont   Action = "cont"
)

type TestEvent struct {
	Time    time.Time `json:"Time,omitempty"`
	Action  Action    `json:"Action,omitempty"`
	Package string    `json:"Package,omitempty"`
	Test    string    `json:"Test,omitempty"`
	Output  string    `json:"Output,omitempty"`
	Elapsed float64   `json:"Elapsed,omitempty"`
}

func (te *TestEvent) String() string {
	// Convert the TestEvent instance to JSON
	jsonBytes, err := json.Marshal(te)
	if err != nil {
		log.Fatalf("Error marshalling TestEvent to JSON: %v", err)
	}

	// Convert bytes to string to get the JSON string representation
	return string(jsonBytes)
}

type TestLogModifierConfig struct {
	IsJsonInput            *bool
	RemoveTLogPrefix       *bool
	RemoveTLogRegexp       *regexp.Regexp
	OnlyErrors             *bool
	CI                     *bool
	ShouldImmediatelyPrint bool
}

// ValidateConfig validates the TestLogModifierConfig does not have any invalid combinations
func (c *TestLogModifierConfig) Validate() error {
	if *c.OnlyErrors {
		if !*c.IsJsonInput {
			return fmt.Errorf("OnlyErrors flag is only valid when run with -json flag")
		}
	}
	return nil
}

type TestLogModifier func(*TestEvent, *TestLogModifierConfig) error

// parseTestEvent parses a byte slice into a TestEvent
func ParseTestEvent(b []byte) (*TestEvent, error) {
	// If a non json line is encountered return nil
	if len(b) <= 0 || b[0] != '{' {
		return nil, nil
	}
	te := &TestEvent{}
	err := json.Unmarshal(b, te)
	return te, err
}

func RemoveTestLogPrefix(te *TestEvent, c *TestLogModifierConfig) error {
	if c.RemoveTLogRegexp == nil {
		c.RemoveTLogRegexp = regexp.MustCompile(testingLogPrefix)
	}
	if te.Action == ActionOutput {
		if len(te.Output) > 0 && c.RemoveTLogRegexp.MatchString(te.Output) {
			te.Output = c.RemoveTLogRegexp.ReplaceAllString(te.Output, "$1")
		}
	}
	return nil
}

type TestPackage struct {
	Name        string
	NonTestLogs []TestEvent
	TestLogs    map[string][]TestEvent
	TestOrder   []string
	FailedTests []string
	PanicTests  []string
	Failed      bool
	Elapsed     float64
}

var TestPackageFailures = map[string]*TestPackage{}

func JsonTestOutputToStandard(te *TestEvent, c *TestLogModifierConfig) error {
	if len(te.Package) == 0 {
		return nil
	}

	// does this package exist in the map
	_, ok := TestPackageFailures[te.Package]
	if !ok {
		TestPackageFailures[te.Package] = &TestPackage{
			Name:        te.Package,
			NonTestLogs: []TestEvent{},
			TestLogs:    map[string][]TestEvent{},
			TestOrder:   []string{},
			FailedTests: []string{},
			PanicTests:  []string{},
		}
	}

	p := TestPackageFailures[te.Package]

	// if this is a test log then make sure it is ordered correctly
	if len(te.Test) > 0 {
		if _, ok := p.TestLogs[te.Test]; !ok {
			p.TestLogs[te.Test] = []TestEvent{}
			p.TestOrder = append(p.TestOrder, te.Test)
		}
		p.TestLogs[te.Test] = append(p.TestLogs[te.Test], *te)

		// if we have a test failure then we add it to the test failures
		if te.Action == ActionFail && len(te.Test) > 0 {
			p.FailedTests = append(p.FailedTests, te.Test)
			p.Failed = true
		}

	} else if (te.Action == ActionFail || te.Action == ActionPass) && len(te.Test) == 0 {
		// if we have a package completed then we can print out the errors if any
		if te.Action == ActionFail {
			p.Failed = true
		}
		p.Elapsed = te.Elapsed
		printPackage(p, c)

		// remove package from map since it has been printed and is no longer needed
		delete(TestPackageFailures, te.Package)
		return nil
	} else {
		p.NonTestLogs = append(p.NonTestLogs, *te)
	}

	// check output for a panic which is a rare case where a panic failed the test but the test was not marked as failed
	// if ('Output' in result && result.Output.includes('panic:')) {
	// 	const pattern = /^panic:.* (Test[A-Z]\w*)/
	// p.PanicTests = append(p.PanicTests, te.Test)

	return nil
}

func printPackage(p *TestPackage, c *TestLogModifierConfig) {
	// if package passed
	if !p.Failed {
		// if we only want errors then skip
		if *c.OnlyErrors {
			return
		}
		// right here is where we would print the passed package with elapsed time if needed
	}

	// start package group
	if *c.CI {
		if p.Failed {
			StartGroupFail(fmt.Sprintf("FAIL  \t%s\t%f", p.Name, p.Elapsed), c)
		} else {
			StartGroupPass(fmt.Sprintf("ok  \t%s\t%f", p.Name, p.Elapsed), c)
		}
	}

	// print the tests in the order of first seen to last seen according to the json logs
	for _, testName := range p.TestOrder {
		test := p.TestLogs[testName]
		shouldPrintLine := false
		// if we only want errors
		if *c.OnlyErrors && p.Failed {
			if len(p.FailedTests) == 0 {
				// we had a package fail without a test fail, we want all the logs for triage in this case
				shouldPrintLine = true
			} else if SliceContains(p.FailedTests, test[0].Test) {
				shouldPrintLine = true
			}
		} else {
			// we want all the logs since we aren't specifying otherwise
			shouldPrintLine = true
		}

		if shouldPrintLine {
			printTest(test, !p.Failed, c)
		}
	}

	// now print the non test logs for the package
	for _, log := range p.NonTestLogs {
		l := log
		PrintEvent(&l, c)
	}

	// end the group if we are in CI mode
	if *c.CI {
		github.EndGroup()
	}
}

func printTest(logs []TestEvent, pass bool, c *TestLogModifierConfig) {
	if !*c.IsJsonInput {
		return // not compatible with non json input
	}
	// start the group
	if *c.CI {
		if pass {
			StartGroupPass(logs[0].Test, c)
		} else {
			StartGroupFail(logs[0].Test, c)
		}
	}

	// print out the test logs
	for _, log := range logs {
		l := log
		PrintEvent(&l, c)
	}

	// end the group if we are in CI mode
	if *c.CI {
		github.EndGroup()
	}
}

func PrintEvent(te *TestEvent, c *TestLogModifierConfig) {
	if te.Output != "" {
		fmt.Print(te.Output)
	}
}

func StartGroupPass(title string, c *TestLogModifierConfig) {
	t := clitext.Color(clitext.ColorGreen, title)
	if *c.CI {
		github.StartGroup(t)
	} else {
		fmt.Println(title)
	}
}
func StartGroupFail(title string, c *TestLogModifierConfig) {
	t := clitext.Color(clitext.ColorRed, title)
	if *c.CI {
		github.StartGroup(t)
	} else {
		fmt.Println(title)
	}
}

func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
