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

	"github.com/smartcontractkit/chainlink-testing-framework/utils/clihelper"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/github"
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
const testRunPrefix = `^=== (RUN|PAUSE|CONT)`

//nolint:gosec
const testPassFailPrefix = `^--- (PASS|FAIL|SKIP):(.*)`

//nolint:gosec
const packagePassFailPrefix = `^(PASS|FAIL)\n$`
const packageCoveragePrefix = `^coverage:`
const testingLogPrefix = `^(\s+)(\w+\.go:\d+: )`
const testPanic = `^panic:.* (Test[A-Z]\w*)`
const testErrorPrefix = `^\s+(Error\sTrace|Error|Test|Messages):\s+`
const testErrorPrefix2 = `        \t            \t.*`
const packageOkFailPrefix = `^(ok|FAIL)\s*\t(.*)` //`^(FAIL|ok).*\t(.*)$`

var testRunPrefixRegexp = regexp.MustCompile(testRunPrefix)
var testPassFailPrefixRegexp = regexp.MustCompile(testPassFailPrefix)
var packagePassFailPrefixRegexp = regexp.MustCompile(packagePassFailPrefix)
var packageCoveragePrefixRegexp = regexp.MustCompile(packageCoveragePrefix)
var removeTLogRegexp = regexp.MustCompile(testingLogPrefix)
var testPanicRegexp = regexp.MustCompile(testPanic)
var testErrorPrefixRegexp = regexp.MustCompile(testErrorPrefix)
var testErrorPrefix2Regexp = regexp.MustCompile(testErrorPrefix2)
var packageOkFailPrefixRegexp = regexp.MustCompile(packageOkFailPrefix)

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
func (t Test) Print(c *TestLogModifierConfig) {
	if !*c.IsJsonInput {
		return // not compatible with non json input
	}

	// preprocess the logs
	message := t[0].Test
	outcomeType := "pass"
	toRemove := []int{}
	errorMessages := []GoTestEvent{}
	for i, log := range t {
		if testPassFailPrefixRegexp.MatchString(log.Output) {
			match := testPassFailPrefixRegexp.FindStringSubmatch(log.Output)
			if strings.Contains(log.Output, "PASS") {
				message = fmt.Sprintf("✅%s", match[2])
				outcomeType = "pass"
			} else if strings.Contains(log.Output, "SKIP") {
				message = fmt.Sprintf("🚧%s", match[2])
				outcomeType = "skip"
			} else {
				message = fmt.Sprintf("❌%s", match[2])
				outcomeType = "fail"
			}
			toRemove = append(toRemove, i)
		} else if *c.CI && (testErrorPrefixRegexp.MatchString(log.Output) || testErrorPrefix2Regexp.MatchString(log.Output) || testPanicRegexp.MatchString(log.Output)) {
			errorMessages = append(errorMessages, log)
		}
	}

	// remove the logs that we don't want to print
	for i := len(toRemove) - 1; i >= 0; i-- {
		t = append(t[:toRemove[i]], t[toRemove[i]+1:]...)
	}

	// start the group
	hasLogs := len(t) > 0
	if outcomeType == "pass" {
		StartGroupPass(message, c, hasLogs)
	} else if outcomeType == "skip" {
		StartGroupSkip(message, c, hasLogs)
	} else {
		StartGroupFail(message, c, hasLogs)
	}

	// print out the error message at the top if the logs are longer than the specified length
	if *c.CI && *c.ErrorAtTopLength > 0 && len(t) > *c.ErrorAtTopLength {
		fmt.Println("❌ Error found:")
		for _, log := range errorMessages {
			log.Print()
		}
	}

	// print out the test logs
	for _, log := range t {
		log.Print()
	}

	// end the group if we are in CI mode
	if *c.CI && hasLogs {
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

	if !*c.SinglePackage {
		// Add color to the output if needed
		if *c.CI || *c.Color {
			match := packageOkFailPrefixRegexp.FindStringSubmatch(p.Message)
			if p.Failed {
				fmt.Printf("📦 %s", clihelper.Color(clihelper.ColorRed, match[2]))
			} else {
				fmt.Printf("📦 %s", clihelper.Color(clihelper.ColorGreen, match[2]))
			}
		}

		p.printTestsInOrder(c)
	}

	// now print the non test logs for the package
	for _, log := range p.NonTestLogs {
		log.Print()
	}
}

func (p TestPackage) printTestsInOrder(c *TestLogModifierConfig) {
	// print the tests in the order of first seen to last seen according to the json logs
	for _, testName := range p.TestOrder {
		test := p.TestLogs[testName]
		if p.ShouldPrintTest(test, c) {
			test.Print(c)
		}
	}
}

func (p TestPackage) ShouldPrintTest(test Test, c *TestLogModifierConfig) bool {
	shouldPrintTest := false
	testFailed := SliceContains(p.FailedTests, test[0].Test)
	// if we only want errors
	if c.OnlyErrors.Value {
		// if the test failed or if we had a package fail without a test fail, we want all the logs for triage in this case
		if (len(p.FailedTests) == 0 || testFailed) && p.Failed {
			shouldPrintTest = true
		}
	} else {
		// we want all the logs since we aren't specifying otherwise
		shouldPrintTest = true
	}
	return shouldPrintTest
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
	OnlyErrors             *clihelper.BoolFlag
	Color                  *bool
	CI                     *bool
	SinglePackage          *bool
	ShouldImmediatelyPrint bool
	TestPackageMap         TestPackageMap
	FailuresExist          bool
	ErrorAtTopLength       *int
}

func NewDefaultConfig() *TestLogModifierConfig {
	return &TestLogModifierConfig{
		IsJsonInput:            ptr.Ptr(false),
		RemoveTLogPrefix:       ptr.Ptr(false),
		OnlyErrors:             &clihelper.BoolFlag{},
		Color:                  ptr.Ptr(false),
		CI:                     ptr.Ptr(false),
		SinglePackage:          ptr.Ptr(false),
		ShouldImmediatelyPrint: false,
		ErrorAtTopLength:       ptr.Ptr(100),
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
		if !*c.IsJsonInput {
			return fmt.Errorf("OnlyErrors flag is only valid when run with -json flag")
		}
	}
	if *c.ErrorAtTopLength < 0 {
		return fmt.Errorf("ErrorAtTopLength must be greater than or equal to 0")
	}

	return nil
}

// SetupModifiers sets up the modifiers based on the flags provided
func SetupModifiers(c *TestLogModifierConfig) []TestLogModifier {
	modifiers := []TestLogModifier{}
	if *c.CI {
		c.Color = ptr.Ptr(true)
		c.IsJsonInput = ptr.Ptr(true)
		c.ShouldImmediatelyPrint = false
		if !c.OnlyErrors.IsSet {
			// nolint errcheck
			c.OnlyErrors.Set("true")
		}
		c.RemoveTLogPrefix = ptr.Ptr(true)
	}
	if *c.RemoveTLogPrefix {
		modifiers = append(modifiers, RemoveTestLogPrefix)
	}
	if *c.Color {
		modifiers = append(modifiers, HighlightErrorOutput)
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
		if testErrorPrefixRegexp.MatchString(te.Output) ||
			testErrorPrefix2Regexp.MatchString(te.Output) {
			te.Output = clihelper.Color(clihelper.ColorRed, te.Output)
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
			c.FailuresExist = true
		}
		t := p.TestLogs[te.Test]
		if *c.SinglePackage && (te.Action == ActionFail || te.Action == ActionPass) && p.ShouldPrintTest(t, c) {
			t.Print(c)
			return nil
		}

	} else if (te.Action == ActionFail || te.Action == ActionPass) && len(te.Test) == 0 {
		// if we have a package completed then we can print out the errors if any
		if te.Action == ActionFail {
			p.Failed = true
			c.FailuresExist = true
		}
		p.Elapsed = te.Elapsed
		p.Print(c)

		// remove package from map since it has been printed and is no longer needed
		delete(c.TestPackageMap, te.Package)
		return nil
	} else {
		// this is package output
		// remove noise from the logs
		if packagePassFailPrefixRegexp.MatchString(te.Output) ||
			packageCoveragePrefixRegexp.MatchString(te.Output) {
			return nil
		}
		if len(te.Output) > 0 && testPanicRegexp.MatchString(te.Output) {
			p.Failed = true
			c.FailuresExist = true
			match := testPanicRegexp.FindStringSubmatch(te.Output)
			te.Output = clihelper.Color(clihelper.ColorRed, te.Output)
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
func StartGroupPass(title string, c *TestLogModifierConfig, hasLogs bool) {
	if *c.Color {
		title = clihelper.Color(clihelper.ColorGreen, title)
	}
	if *c.CI && hasLogs {
		github.StartGroup(title)
	} else {
		fmt.Print(title)
	}
}

// StartGroupSkip starts a group in the CI environment with a green title
func StartGroupSkip(title string, c *TestLogModifierConfig, hasLogs bool) {
	if *c.Color {
		title = clihelper.Color(clihelper.ColorYellow, title)
	}
	if *c.CI && hasLogs {
		github.StartGroup(title)
	} else {
		fmt.Print(title)
	}
}

// StartGroupFail starts a group in the CI environment with a red title
func StartGroupFail(title string, c *TestLogModifierConfig, hasLogs bool) {
	if *c.Color {
		title = clihelper.Color(clihelper.ColorRed, title)
	}
	if *c.CI && hasLogs {
		github.StartGroup(title)
	} else {
		fmt.Print(title)
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
	return clihelper.ReadLine(ctx, r, func(b []byte) error {
		var te *GoTestEvent
		var err error

		// build a TestEvent from the input line
		if *c.IsJsonInput {
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
