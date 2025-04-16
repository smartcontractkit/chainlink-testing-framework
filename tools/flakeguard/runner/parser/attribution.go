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
	nestedTestNameRe   = regexp.MustCompile(`\.(Test[^\s]+?)(?:\.[^(]+)?\s*\(`)
	testLogAfterTestRe = regexp.MustCompile(`^panic: Log in goroutine after (Test[^\s]+) has completed:`)
	didTestTimeoutRe   = regexp.MustCompile(`^panic: test timed out after ([^\s]+)`)
	timedOutTestNameRe = regexp.MustCompile(`^\s*(Test[^\s]+)\s+\((.*)\)`)

	ErrFailedToAttributePanicToTest              = errors.New("failed to attribute panic to test")
	ErrFailedToAttributeRaceToTest               = errors.New("failed to attribute race to test")
	ErrFailedToParseTimeoutDuration              = errors.New("failed to parse timeout duration")
	ErrFailedToExtractTimeoutDuration            = errors.New("failed to extract timeout duration")
	ErrDetectedLogAfterCompleteFailedAttribution = errors.New("detected a log after test has completed panic, but failed to properly attribute it")
	ErrDetectedTimeoutFailedParse                = errors.New("detected test timeout, but failed to parse the duration from the test")
	ErrDetectedTimeoutFailedAttribution          = errors.New("detected test timeout, but failed to attribute the timeout to a specific test")
)

// AttributePanicToTest attributes panics to the test that caused them.
func AttributePanicToTest(outputs []string) (test string, timeout bool, err error) {
	var (
		timeoutDurationStr string
		timeoutDuration    time.Duration
		foundTestName      string
	)

	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output == "" {
			continue
		}

		if match := testLogAfterTestRe.FindStringSubmatch(output); len(match) > 1 {
			testName := strings.TrimSpace(match[1])
			log.Debug().Str("test", testName).Str("line", output).Msg("Attributed panic via LogAfterTest pattern")
			return testName, false, nil
		}

		if match := didTestTimeoutRe.FindStringSubmatch(output); len(match) > 1 {
			timeout = true
			timeoutDurationStr = match[1]
			var parseErr error
			timeoutDuration, parseErr = time.ParseDuration(timeoutDurationStr)
			if parseErr != nil {
				log.Warn().Str("duration_str", timeoutDurationStr).Err(parseErr).Msg("Failed to parse timeout duration from initial panic line")
				return "", true, fmt.Errorf("%w: %w using output line: %s", ErrFailedToParseTimeoutDuration, parseErr, output)
			}
			log.Debug().Dur("duration", timeoutDuration).Msg("Detected timeout panic")
			continue
		}

		if timeout {
			if match := timedOutTestNameRe.FindStringSubmatch(output); len(match) > 2 {
				testName := strings.TrimSpace(match[1])
				testDurationStr := strings.TrimSpace(match[2])
				testDuration, parseErr := time.ParseDuration(testDurationStr)
				if parseErr != nil {
					log.Warn().Str("test", testName).Str("duration_str", testDurationStr).Err(parseErr).Msg("Failed to parse duration from timed-out test line")
					return "", true, fmt.Errorf("%w: test '%s' listed with unparseable duration '%s': %w", ErrDetectedTimeoutFailedParse, testName, testDurationStr, parseErr)
				} else if testDuration >= timeoutDuration {
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Attributed timeout panic via duration match")
					return testName, true, nil
				} else {
					log.Debug().Str("test", testName).Dur("test_duration", testDuration).Dur("timeout_duration", timeoutDuration).Msg("Ignoring test line, duration too short for timeout")
				}
			}
		}

		matchNestedTestName := nestedTestNameRe.FindStringSubmatch(output)
		if len(matchNestedTestName) > 1 {
			testName := strings.TrimSpace(matchNestedTestName[1])
			if !strings.HasPrefix(testName, "Test") {
				continue
			}
			if foundTestName == "" {
				log.Debug().Str("test", testName).Str("line", output).Msg("Found potential test name in panic output")
				foundTestName = testName
			}
		}
	}

	if timeout {
		var errMsg string
		if foundTestName != "" {
			errMsg = fmt.Sprintf("timeout duration %s detected, found candidate test '%s' but duration did not meet threshold or failed parsing earlier", timeoutDurationStr, foundTestName)
		} else {
			errMsg = fmt.Sprintf("timeout duration %s detected, but no matching test found in output", timeoutDurationStr)
		}
		return "", true, fmt.Errorf("%w: %s: %w", ErrDetectedTimeoutFailedAttribution, errMsg, errors.New(strings.Join(outputs, "\n")))
	}

	if foundTestName != "" {
		log.Debug().Str("test", foundTestName).Msg("Attributed non-timeout panic via test name found in stack")
		return foundTestName, false, nil
	}

	return "", false, fmt.Errorf("%w: using output: %w", ErrFailedToAttributePanicToTest, errors.New(strings.Join(outputs, "\n")))
}

// AttributeRaceToTest attributes races to the test that caused them.
func AttributeRaceToTest(outputs []string) (string, error) {
	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output == "" {
			continue
		}
		match := nestedTestNameRe.FindStringSubmatch(output)
		if len(match) > 1 {
			testName := strings.TrimSpace(match[1])
			if strings.HasPrefix(testName, "Test") {
				log.Debug().Str("test", testName).Str("line", output).Msg("Attributed race via test name match")
				return testName, nil
			}
		}
	}
	return "", fmt.Errorf("%w: using output: %w", ErrFailedToAttributeRaceToTest, errors.New(strings.Join(outputs, "\n")))
}
