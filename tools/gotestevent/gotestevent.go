package gotestevent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"dario.cat/mergo"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/clireader"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/clitext"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flags"
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

// regex strings
const testRunPrefix = `^=== (RUN|PAUSE|CONT)   `

//nolint:gosec
const testPassFailPrefix = `^--- (PASS|FAIL):(.*)`

//nolint:gosec
const packagePassFailPrefix = `^(PASS|FAIL)\n$`
const testingLogPrefix = `^(\s+)(\w+\.go:\d+: )`
const testPanic = `^panic:.* (Test[A-Z]\w*)`
const testErrorPrefix = `^\s+(Error\sTrace|Error|Test):\s+`

var testRunPrefixRegexp = regexp.MustCompile(testRunPrefix)
var testPassFailPrefixRegexp = regexp.MustCompile(testPassFailPrefix)
var packagePassFailPrefixRegexp = regexp.MustCompile(packagePassFailPrefix)
var removeTLogRegexp = regexp.MustCompile(testingLogPrefix)
var testPanicRegexp = regexp.MustCompile(testPanic)
var testErrorPrefixRegexp = regexp.MustCompile(testErrorPrefix)

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
	if strings.TrimSpace(gte.Output) != "" {
		fmt.Print(gte.Output)
	}
}

type Test []GoTestEvent

// Print prints the Test to the console
func (t Test) Print(pass bool, c *TestLogModifierConfig) {
	if !ptr.Val(c.IsJsonInput) {
		return // not compatible with non json input
	}

	// preprocess the logs
	message := t[0].Test
	toRemove := []int{}
	for i, log := range t {
		if testPassFailPrefixRegexp.MatchString(log.Output) {

			match := testPassFailPrefixRegexp.FindStringSubmatch(log.Output)
			if strings.Contains(log.Output, "PASS") {
				message = fmt.Sprintf("✅ %s", match[1])
			} else {
				message = fmt.Sprintf("X %s", match[1])
			}
			toRemove = append(toRemove, i)
		}
	}

	// remove the logs that we don't want to print
	for i := len(toRemove) - 1; i >= 0; i-- {
		t = append(t[:toRemove[i]], t[toRemove[i]+1:]...)
	}

	// start the group
	if ptr.Val(c.CI) {
		if pass {
			StartGroupPass(message, c)
		} else {
			StartGroupFail(message, c)
		}
	}

	// print out the test logs
	for _, log := range t {
		log.Print()
	}

	// end the group if we are in CI mode
	if ptr.Val(c.CI) {
		github.EndGroup()
	}
}

const packageOkFailPrefix = `^(ok|FAIL)\s*\t`

var packageOkFailPrefixRegexp = regexp.MustCompile(packageOkFailPrefix)

type TestPackage struct {
	Name        string
	NonTestLogs []GoTestEvent
	TestLogs    map[string]Test
	TestOrder   []string
	FailedTests []string
	PanicTests  []string
	Failed      bool
	Elapsed     float64
	Message     string
}

func (p *TestPackage) AddTestEvent(te *GoTestEvent) {
	// stop noise from being added to the logs
	if testRunPrefixRegexp.MatchString(te.Output) {
		return
	}

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
		if c.OnlyErrors.Value {
			return
		}
		// right here is where we would print the passed package with elapsed time if needed
	}

	// Add color to the output if needed
	if ptr.Val(c.CI) || ptr.Val(c.Color) {
		if p.Failed {
			fmt.Printf("📦 %s", clitext.Color(clitext.ColorRed, p.Message))
		} else {
			fmt.Printf("📦 %s", clitext.Color(clitext.ColorGreen, p.Message))
		}
	}

	p.printTestsInOrder(c)

	// now print the non test logs for the package
	for _, log := range p.NonTestLogs {
		log.Print()
	}
}

func (p TestPackage) printTestsInOrder(c *TestLogModifierConfig) {
	// print the tests in the order of first seen to last seen according to the json logs
	for _, testName := range p.TestOrder {
		test := p.TestLogs[testName]
		shouldPrintLine := false
		testFailed := SliceContains(p.FailedTests, test[0].Test)
		// if we only want errors
		if c.OnlyErrors.Value && p.Failed {
			if len(p.FailedTests) == 0 {
				// we had a package fail without a test fail, we want all the logs for triage in this case
				shouldPrintLine = true
			} else if testFailed {
				shouldPrintLine = true
			}
		} else {
			// we want all the logs since we aren't specifying otherwise
			shouldPrintLine = true
		}

		if shouldPrintLine {
			test.Print(!testFailed, c)
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
	OnlyErrors             *flags.BoolFlag
	Color                  *bool
	CI                     *bool
	ShouldImmediatelyPrint bool
	TestPackageMap         TestPackageMap
}

func NewDefaultConfig() *TestLogModifierConfig {
	return &TestLogModifierConfig{
		IsJsonInput:            ptr.Ptr(false),
		RemoveTLogPrefix:       ptr.Ptr(false),
		OnlyErrors:             &flags.BoolFlag{},
		Color:                  ptr.Ptr(false),
		CI:                     ptr.Ptr(false),
		ShouldImmediatelyPrint: false,
	}
}

// ValidateConfig validates the TestLogModifierConfig does not have any invalid combinations
func (c *TestLogModifierConfig) Validate() error {
	defaultConfig := NewDefaultConfig()
	err := mergo.Merge(c, defaultConfig)
	if err != nil {
		return err
	}
	if c.OnlyErrors.Value {
		if !ptr.Val(c.IsJsonInput) {
			return fmt.Errorf("OnlyErrors flag is only valid when run with -json flag")
		}
	}

	return nil
}

// SetupModifiers sets up the modifiers based on the flags provided
func SetupModifiers(c *TestLogModifierConfig) []TestLogModifier {
	modifiers := []TestLogModifier{}
	if ptr.Val(c.CI) {
		c.Color = ptr.Ptr(true)
		c.IsJsonInput = ptr.Ptr(true)
		c.ShouldImmediatelyPrint = false
		if !c.OnlyErrors.IsSet {
			c.OnlyErrors.Set("true")
		}
		c.RemoveTLogPrefix = ptr.Ptr(true)
	}
	if ptr.Val(c.RemoveTLogPrefix) {
		modifiers = append(modifiers, RemoveTestLogPrefix)
	}
	if ptr.Val(c.IsJsonInput) {
		c.ShouldImmediatelyPrint = false
		modifiers = append(modifiers, JsonTestOutputToStandard)
	}
	if ptr.Val(c.Color) {
		modifiers = append(modifiers, HighlightErrorOutput)
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

// RemoveTestLogPrefix is a TestLogModifier that takes a GoTestEvent and removes the test log prefix
func RemoveTestLogPrefix(te *GoTestEvent, _ *TestLogModifierConfig) error {
	if te.Action == ActionOutput && len(te.Output) > 0 {
		if removeTLogRegexp.MatchString(te.Output) {
			te.Output = removeTLogRegexp.ReplaceAllString(te.Output, "$1")
		}
	}
	return nil
}

func HighlightErrorOutput(te *GoTestEvent, _ *TestLogModifierConfig) error {
	if te.Action == ActionOutput && len(te.Output) > 0 {
		if testErrorPrefixRegexp.MatchString(te.Output) {
			te.Output = clitext.Color(clitext.ColorRed, te.Output)
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
		// this is package output
		// remove noise from the logs
		if packagePassFailPrefixRegexp.MatchString(te.Output) {
			return nil
		}
		if len(te.Output) > 0 && testPanicRegexp.MatchString(te.Output) {
			p.Failed = true
			match := testPanicRegexp.FindStringSubmatch(te.Output)
			te.Output = clitext.Color(clitext.ColorRed, te.Output)
			p.NonTestLogs = append(p.NonTestLogs, *te)
			// Check if there is a match for the test name
			if len(match) > 1 {
				// the second element should have the test name
				p.FailedTests = append(p.FailedTests, match[1])
			} else {
				fmt.Println("What is wrong with this panic???", te.Output)
			}
		} else if len(te.Output) > 0 && packageOkFailPrefixRegexp.MatchString(te.Output) {
			p.Message = te.Output
		} else {
			p.NonTestLogs = append(p.NonTestLogs, *te)
		}
	}

	return nil
}

// StartGroupPass starts a group in the CI environment with a green title
func StartGroupPass(title string, c *TestLogModifierConfig) {
	if ptr.Val(c.Color) {
		title = clitext.Color(clitext.ColorGreen, title)
	}
	if ptr.Val(c.CI) {
		github.StartGroup(title)
	} else {
		fmt.Println(title)
	}
}

// StartGroupFail starts a group in the CI environment with a red title
func StartGroupFail(title string, c *TestLogModifierConfig) {
	if ptr.Val(c.Color) {
		title = clitext.Color(clitext.ColorRed, title)
	}
	if ptr.Val(c.CI) {
		github.StartGroup(title)
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
