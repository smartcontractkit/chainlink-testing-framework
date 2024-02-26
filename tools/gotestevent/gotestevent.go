package gotestevent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/clireader"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/clitext"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/github"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

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

// Represntation of a go test -json event
type GoTestEvent struct {
	Time    time.Time `json:"Time,omitempty"`
	Action  Action    `json:"Action,omitempty"`
	Package string    `json:"Package,omitempty"`
	Test    string    `json:"Test,omitempty"`
	Output  string    `json:"Output,omitempty"`
	Elapsed float64   `json:"Elapsed,omitempty"`
}

// String returns the JSON string representation of the GoTestEvent
func (gte GoTestEvent) String() (string, error) {
	// Convert the TestEvent instance to JSON
	jsonBytes, err := json.Marshal(gte)
	if err != nil {
		return "", err
	}

	// Convert bytes to string to get the JSON string representation
	return string(jsonBytes), nil
}

// Print prints the GoTestEvent to the console
func (gte GoTestEvent) Print() {
	if gte.Output != "" {
		fmt.Print(gte.Output)
	}
}

type Test []GoTestEvent

// Print prints the Test to the console
func (t Test) Print(pass bool, c *TestLogModifierConfig) {
	if !ptr.Val(c.IsJsonInput) {
		return // not compatible with non json input
	}
	// start the group
	if ptr.Val(c.CI) {
		if pass {
			StartGroupPass(t[0].Test, c)
		} else {
			StartGroupFail(t[0].Test, c)
		}
	}

	// print out the test logs
	for _, log := range t {
		l := log
		l.Print()
	}

	// end the group if we are in CI mode
	if ptr.Val(c.CI) {
		github.EndGroup()
	}
}

type TestPackage struct {
	Name        string
	NonTestLogs []GoTestEvent
	TestLogs    map[string]Test
	TestOrder   []string
	FailedTests []string
	PanicTests  []string
	Failed      bool
	Elapsed     float64
}

func (p *TestPackage) AddTestEvent(te *GoTestEvent) {
	if _, ok := p.TestLogs[te.Test]; !ok {
		p.TestLogs[te.Test] = []GoTestEvent{}
		p.TestOrder = append(p.TestOrder, te.Test)
	}
	p.TestLogs[te.Test] = append(p.TestLogs[te.Test], *te)
}

// Print prints the TestPackage to the console
func (p TestPackage) Print(c *TestLogModifierConfig) {
	// if package passed
	if !p.Failed {
		// if we only want errors then skip
		if ptr.Val(c.OnlyErrors) {
			return
		}
		// right here is where we would print the passed package with elapsed time if needed
	}

	// start package group
	if ptr.Val(c.CI) {
		if p.Failed {
			StartGroupFail(fmt.Sprintf("FAIL  \t%s\t%f", p.Name, p.Elapsed), c)
		} else {
			StartGroupPass(fmt.Sprintf("ok  \t%s\t%f", p.Name, p.Elapsed), c)
		}
	}

	p.printTestsInOrder(c)

	// now print the non test logs for the package
	for _, log := range p.NonTestLogs {
		l := log
		l.Print()
	}

	// end the group if we are in CI mode
	if ptr.Val(c.CI) {
		github.EndGroup()
	}
}

func (p TestPackage) printTestsInOrder(c *TestLogModifierConfig) {
	// print the tests in the order of first seen to last seen according to the json logs
	for _, testName := range p.TestOrder {
		test := p.TestLogs[testName]
		shouldPrintLine := false
		// if we only want errors
		if ptr.Val(c.OnlyErrors) && p.Failed {
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
			test.Print(!p.Failed, c)
		}
	}
}

type TestPackageMap map[string]*TestPackage

func (m TestPackageMap) InitPackageInMap(packageName string) {
	_, ok := m[packageName]
	if !ok {
		m[packageName] = &TestPackage{
			Name:        packageName,
			NonTestLogs: []GoTestEvent{},
			TestLogs:    map[string]Test{},
			TestOrder:   []string{},
			FailedTests: []string{},
			PanicTests:  []string{},
		}
	}
}

type TestLogModifierConfig struct {
	IsJsonInput            *bool
	RemoveTLogPrefix       *bool
	OnlyErrors             *bool
	CI                     *bool
	ShouldImmediatelyPrint bool
	TestPackageMap         TestPackageMap
}

// ValidateConfig validates the TestLogModifierConfig does not have any invalid combinations
func (c TestLogModifierConfig) Validate() error {
	if ptr.Val(c.OnlyErrors) {
		if !ptr.Val(c.IsJsonInput) {
			return fmt.Errorf("OnlyErrors flag is only valid when run with -json flag")
		}
	}
	return nil
}

// SetupModifiers sets up the modifiers based on the flags provided
func SetupModifiers(c *TestLogModifierConfig) []TestLogModifier {
	modifiers := []TestLogModifier{}
	if *c.RemoveTLogPrefix {
		modifiers = append(modifiers, RemoveTestLogPrefix)
	}
	if *c.IsJsonInput {
		c.ShouldImmediatelyPrint = false
		modifiers = append(modifiers, JsonTestOutputToStandard)
	}
	return modifiers
}

// TestLogModifier is a generic function interface that modifies a GoTestEvent
type TestLogModifier func(*GoTestEvent, *TestLogModifierConfig) error

// parseTestEvent parses a byte slice into a TestEvent
func ParseTestEvent(b []byte) (*GoTestEvent, error) {
	// If a non json line is encountered return nil
	if len(b) <= 0 || b[0] != '{' {
		return nil, nil
	}
	te := &GoTestEvent{}
	err := json.Unmarshal(b, te)
	return te, err
}

const testingLogPrefix = `^(\s+)(\w+\.go:\d+: )`
const testPanic = `^panic:.* (Test[A-Z]\w*)`

var removeTLogRegexp = regexp.MustCompile(testingLogPrefix)
var testPanicRegexp = regexp.MustCompile(testPanic)

// RemoveTestLogPrefix is a TestLogModifier that takes a GoTestEvent and removes the test log prefix
func RemoveTestLogPrefix(te *GoTestEvent, _ *TestLogModifierConfig) error {
	if te.Action == ActionOutput {
		if len(te.Output) > 0 && removeTLogRegexp.MatchString(te.Output) {
			te.Output = removeTLogRegexp.ReplaceAllString(te.Output, "$1")
		}
	}
	return nil
}

// JsonTestOutputToStandard is a TestLogModifier that takes a GoTestEvent and modifies the output as configured
func JsonTestOutputToStandard(te *GoTestEvent, c *TestLogModifierConfig) error {
	if len(te.Package) == 0 {
		return nil
	}

	if c.TestPackageMap == nil {
		c.TestPackageMap = make(TestPackageMap)
	}
	// does this package exist in the map
	c.TestPackageMap.InitPackageInMap(te.Package)

	p := c.TestPackageMap[te.Package]

	// if this is a test log then make sure it is ordered correctly
	if len(te.Test) > 0 {
		p.AddTestEvent(te)

		// if we have a test failure or panic then we add it to the test failures
		if te.Action == ActionFail || testPanicRegexp.MatchString(te.Output) {
			p.FailedTests = append(p.FailedTests, te.Test)
			p.Failed = true
		}

	} else if (te.Action == ActionFail || te.Action == ActionPass) && len(te.Test) == 0 {
		// if we have a package completed then we can print out the errors if any
		if te.Action == ActionFail {
			p.Failed = true
		}
		p.Elapsed = te.Elapsed
		p.Print(c)

		// remove package from map since it has been printed and is no longer needed
		delete(c.TestPackageMap, te.Package)
		return nil
	} else {
		p.NonTestLogs = append(p.NonTestLogs, *te)
	}

	return nil
}

// StartGroupPass starts a group in the CI environment with a green title
func StartGroupPass(title string, c *TestLogModifierConfig) {
	t := clitext.Color(clitext.ColorGreen, title)
	if ptr.Val(c.CI) {
		github.StartGroup(t)
	} else {
		fmt.Println(title)
	}
}

// StartGroupFail starts a group in the CI environment with a red title
func StartGroupFail(title string, c *TestLogModifierConfig) {
	t := clitext.Color(clitext.ColorRed, title)
	if ptr.Val(c.CI) {
		github.StartGroup(t)
	} else {
		fmt.Println(title)
	}
}

// SliceContains checks if a slice contains a given item
func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ReadAndModifyLogs(ctx context.Context, r io.Reader, modifiers []TestLogModifier, c *TestLogModifierConfig) error {
	return clireader.ReadLine(ctx, r, func(b []byte) error {
		var te *GoTestEvent
		var err error

		// build a TestEvent from the input line
		if ptr.Val(c.IsJsonInput) {
			te, err = ParseTestEvent(b)
			if err != nil {
				log.Fatalf("Error parsing json test event from stdin: %v\n", err)
			}
			if te == nil {
				// got a non json line when expecting json, just print it out and move on
				fmt.Println(string(b))
				return nil
			}
		} else {
			te = &GoTestEvent{}
			te.Action = ActionOutput
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
			if ptr.Val(c.IsJsonInput) {
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
