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
	// Use the more precise original regex to avoid capturing .func suffixes
	nestedTestNameRe = regexp.MustCompile(`\.(Test[^\s]+?)(?:\.[^(]+)?\s*\(`)
	// Other regexes remain the same
	testLogAfterTestRe = regexp.MustCompile(`^panic: Log in goroutine after (Test[^\s]+) has completed:`)
	didTestTimeoutRe   = regexp.MustCompile(`^panic: test timed out after ([^\s]+)`)
	timedOutTestNameRe = regexp.MustCompile(`^\s*(Test[^\s]+)\s+\((.*)\)`)

	// Exported Errors
	ErrFailedToAttributePanicToTest              = errors.New("failed to attribute panic to test")
	ErrFailedToAttributeRaceToTest               = errors.New("failed to attribute race to test")
	ErrFailedToParseTimeoutDuration              = errors.New("failed to parse timeout duration")
	ErrFailedToExtractTimeoutDuration            = errors.New("failed to extract timeout duration")
	ErrDetectedLogAfterCompleteFailedAttribution = errors.New("detected a log after test has completed panic, but failed to properly attribute it")
	ErrDetectedTimeoutFailedParse                = errors.New("detected test timeout, but failed to parse the duration from the test")
	ErrDetectedTimeoutFailedAttribution          = errors.New("detected test timeout, but failed to attribute the timeout to a specific test")
)

// AttributePanicToTest properly attributes panics to the test that caused them.
func AttributePanicToTest(outputs []string) (test string, timeout bool, err error) {
	var (
		timeoutDurationStr string
		timeoutDuration    time.Duration
		foundTestName      string // Store first plausible test name found as fallback
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
				log.Warn().Str("duration_str", timeoutDurationStr).Err(parseErr).Msg("Failed to parse timeout duration from initial panic line")
				// Use Errorf here to wrap the error
				return "", true, fmt.Errorf("%w: %w using output line: %s", ErrFailedToParseTimeoutDuration, parseErr, output)
			}
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
					// If duration parsing fails for a candidate test, immediately return a specific error
					return "", true, fmt.Errorf("%w: test '%s' listed with unparseable duration '%s': %w", ErrDetectedTimeoutFailedParse, testName, testDurationStr, parseErr)
				} else if testDuration >= timeoutDuration {
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Attributed timeout panic via duration match")
					return testName, true, nil // Found a valid match!
				} else {
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Ignoring test line, duration too short for timeout")
				}
			}
		}

		// General check for test names using the more precise regex
		matchNestedTestName := nestedTestNameRe.FindStringSubmatch(output)
		if len(matchNestedTestName) > 1 {
			testName := strings.TrimSpace(matchNestedTestName[1]) // Group 1 captures the core test name
			if !strings.HasPrefix(testName, "Test") {
				continue // Should not happen with this regex, but safety check
			}
			// Store the first plausible name found as a fallback
			if foundTestName == "" {
				log.Debug().Str("test", testName).Str("line", output).Msg("Found potential test name in panic output")
				foundTestName = testName
			}
		}
	} // End loop over outputs

	// Post-loop evaluation
	if timeout {
		// If we reach here, timeout was detected, but NO line matched BOTH name and duration threshold.
		// Return the generic attribution failure error.
		var errMsg string
		if foundTestName != "" {
			// Include the fallback name if found, even though its duration didn't match/parse.
			errMsg = fmt.Sprintf("timeout duration %s detected, found candidate test '%s' but duration did not meet threshold or failed parsing earlier", timeoutDurationStr, foundTestName)
		} else {
			errMsg = fmt.Sprintf("timeout duration %s detected, but no matching test found in output", timeoutDurationStr)
		}
		return "", true, fmt.Errorf("%w: %s: %w", ErrDetectedTimeoutFailedAttribution, errMsg, errors.New(strings.Join(outputs, "\n")))
	}

	if foundTestName != "" {
		// If not a timeout, but we found a test name via the general regex
		log.Debug().Str("test", foundTestName).Msg("Attributed non-timeout panic via test name found in stack")
		return foundTestName, false, nil
	}

	// If we reach here, no pattern matched successfully for non-timeout panic
	// Use Errorf for the final error wrapping
	return "", false, fmt.Errorf("%w: using output: %w", ErrFailedToAttributePanicToTest, errors.New(strings.Join(outputs, "\n")))
}

// AttributeRaceToTest properly attributes races to the test that caused them.
func AttributeRaceToTest(outputs []string) (string, error) {
	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output == "" {
			continue
		}
		// Use the precise regex here too
		match := nestedTestNameRe.FindStringSubmatch(output)
		if len(match) > 1 {
			testName := strings.TrimSpace(match[1]) // Group 1 captures the core test name
			if strings.HasPrefix(testName, "Test") {
				log.Debug().Str("test", testName).Str("line", output).Msg("Attributed race via test name match")
				return testName, nil
			}
		}
	}
	// Use Errorf for the final error wrapping if loop completes without match
	return "", fmt.Errorf("%w: using output: %w",
		ErrFailedToAttributeRaceToTest, errors.New(strings.Join(outputs, "\n")))
}
