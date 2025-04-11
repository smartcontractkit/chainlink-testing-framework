package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// Regex to extract a valid test function name from a panic message.
	// This is the most common situation for test panics, e.g.
	// github.com/smartcontractkit/chainlink/deployment/keystone/changeset_test.TestDeployBalanceReader(0xc000583c00)
	nestedTestNameRe = regexp.MustCompile(`\.(Test[^\s(]+)`) // Simpler regex, matches TestName directly after a dot

	// Regex to check if the panic is from a log after a goroutine, e.g.
	// panic: Log in goroutine after Test_workflowRegisteredHandler/skips_fetch_if_secrets_url_is_missing has completed: <Log line>
	testLogAfterTestRe = regexp.MustCompile(`^panic: Log in goroutine after (Test[^\s]+) has completed:`)

	// Check if the panic message indicates a timeout, e.g.
	// panic: test timed out after 10m0s
	didTestTimeoutRe = regexp.MustCompile(`^panic: test timed out after ([^\s]+)`)
	// Regex to extract a valid test function name from a panic message if the panic is a timeout, e.g.
	// TestTimedOut (10m0s)
	timedOutTestNameRe = regexp.MustCompile(`^\s*(Test[^\s]+)\s+\((.*)\)`) // Added optional leading space

	ErrFailedToAttributePanicToTest              = errors.New("failed to attribute panic to test")
	ErrFailedToAttributeRaceToTest               = errors.New("failed to attribute race to test")
	ErrFailedToParseTimeoutDuration              = errors.New("failed to parse timeout duration")
	ErrFailedToExtractTimeoutDuration            = errors.New("failed to extract timeout duration")
	ErrDetectedLogAfterCompleteFailedAttribution = errors.New("detected a log after test has completed panic, but failed to properly attribute it")
	ErrDetectedTimeoutFailedParse                = errors.New("detected test timeout, but failed to parse the duration from the test")
	ErrDetectedTimeoutFailedAttribution          = errors.New("detected test timeout, but failed to attribute the timeout to a specific test")
)

// attributePanicToTest properly attributes panics to the test that caused them.
// There are a lot of edge cases and strange behavior in Go test output when it comes to panics.
func attributePanicToTest(outputs []string) (test string, timeout bool, err error) {
	var (
		timeoutDurationStr string
		timeoutDuration    time.Duration
		foundTestName      string // Store first plausible test name found
	)

	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output == "" {
			continue
		}

		// Check for specific panic patterns first
		if match := testLogAfterTestRe.FindStringSubmatch(output); len(match) > 1 {
			testName := strings.TrimSpace(match[1])
			log.Debug().Str("test", testName).Str("line", output).Msg("Attributed panic via LogAfterTest pattern")
			return testName, false, nil // Found definitive match
		}

		if match := didTestTimeoutRe.FindStringSubmatch(output); len(match) > 1 {
			timeout = true
			timeoutDurationStr = match[1]
			var parseErr error
			timeoutDuration, parseErr = time.ParseDuration(timeoutDurationStr)
			if parseErr != nil {
				// Log error but continue searching, maybe timeout reported differently later
				log.Warn().Str("duration_str", timeoutDurationStr).Err(parseErr).Msg("Failed to parse timeout duration from initial panic line")
				// Return error immediately? Or hope timedOutTestNameRe finds it? Let's return error.
				return "", true, fmt.Errorf("%w: %w using output line: %s", ErrFailedToParseTimeoutDuration, parseErr, output)
			}
			// Don't return yet, need to find the specific timed-out test name below
			log.Debug().Dur("duration", timeoutDuration).Msg("Detected timeout panic")
			continue // Continue scanning for the test name line
		}

		// If timeout detected, look for the "TestName (duration)" pattern
		if timeout {
			if match := timedOutTestNameRe.FindStringSubmatch(output); len(match) > 2 {
				testName := strings.TrimSpace(match[1])
				testDurationStr := strings.TrimSpace(match[2])
				testDuration, parseErr := time.ParseDuration(testDurationStr)
				if parseErr != nil {
					log.Warn().Str("test", testName).Str("duration_str", testDurationStr).Err(parseErr).Msg("Failed to parse duration from timed-out test line")
					// If we already have a timeoutDuration, maybe use this test name anyway?
					// Let's continue searching for a perfect match first. Store this as potential.
					if foundTestName == "" {
						foundTestName = testName
					}
				} else if testDuration >= timeoutDuration {
					// Found the test that likely caused the timeout
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Attributed timeout panic via duration match")
					return testName, true, nil // Found definitive match
				} else {
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Ignoring test line, duration too short for timeout")
				}
			}
		}

		// General check for test names within stack trace lines (less reliable but a fallback)
		if match := nestedTestNameRe.FindStringSubmatch(output); len(match) > 1 {
			testName := strings.TrimSpace(match[1])
			// Avoid standard library test runners or internal functions
			if !strings.HasPrefix(testName, "Test") {
				continue
			}
			// Prioritize longer, more specific names if multiple matches found?
			// For now, store the first plausible one found if we haven't found one yet.
			if foundTestName == "" {
				log.Debug().Str("test", testName).Str("line", output).Msg("Found potential test name in panic output")
				foundTestName = testName
				// Don't return yet, keep searching for more specific patterns (like timeout or log after test)
			}
		}
	} // End loop over outputs

	// Post-loop evaluation
	if timeout {
		if foundTestName != "" {
			// If timeout was detected, and we found a potential test name (maybe without duration match), use it.
			log.Warn().Str("test", foundTestName).Msg("Attributing timeout to test name found, but duration match was inconclusive or missing.")
			return foundTestName, true, nil
		}
		// If timeout detected but no test name found anywhere
		return "", true, fmt.Errorf("%w in package context using output:\n%s", ErrDetectedTimeoutFailedAttribution, strings.Join(outputs, "\n"))
	}

	if foundTestName != "" {
		// If not a timeout, but we found a test name in the stack trace
		log.Debug().Str("test", foundTestName).Msg("Attributed non-timeout panic via test name found in stack")
		return foundTestName, false, nil
	}

	// If we reach here, no pattern matched successfully
	return "", false, fmt.Errorf("%w using output:\n%s", ErrFailedToAttributePanicToTest, strings.Join(outputs, "\n"))
}

// attributeRaceToTest properly attributes races to the test that caused them.
// Race output often includes stack traces mentioning the test function.
func attributeRaceToTest(outputs []string) (string, error) {
	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output == "" {
			continue
		}
		// Use the same regex as panic attribution fallback
		if match := nestedTestNameRe.FindStringSubmatch(output); len(match) > 1 {
			testName := strings.TrimSpace(match[1])
			if strings.HasPrefix(testName, "Test") {
				log.Debug().Str("test", testName).Str("line", output).Msg("Attributed race via test name match")
				return testName, nil
			}
		}
	}
	// If no match found in any line
	return "", fmt.Errorf("%w using output:\n%s",
		ErrFailedToAttributeRaceToTest, strings.Join(outputs, "\n"),
	)
}
