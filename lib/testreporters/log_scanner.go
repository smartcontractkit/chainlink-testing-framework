package testreporters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

var (
	OneLogAtLogLevelErr       = "found log at level"
	MultipleLogsAtLogLevelErr = "found too many logs at level"
)

// ScanLogLine scans a log line for a failing log level, returning the number of failing logs found so far. It returns an error if the failure threshold is reached or if any panic is found
// or if there's no log level found. It also takes a list of allowed messages that are ignored if found.
func ScanLogLine(log zerolog.Logger, jsonLogLine string, failingLogLevel zapcore.Level, foundSoFar, failureThreshold uint, allowedMessages []AllowedLogMessage) (uint, error) {
	var zapLevel zapcore.Level
	var err error

	if !strings.HasPrefix(jsonLogLine, "{") { // don't bother with non-json lines
		if strings.HasPrefix(jsonLogLine, "panic") { // unless it's a panic
			return 0, fmt.Errorf("found panic: %s", jsonLogLine)
		}
		return foundSoFar, nil
	}
	jsonMapping := map[string]any{}

	if err = json.Unmarshal([]byte(jsonLogLine), &jsonMapping); err != nil {
		// This error can occur anytime someone uses %+v in a log message, ignoring
		return foundSoFar, nil
	}
	logLevel, ok := jsonMapping["level"].(string)
	if !ok {
		return 0, fmt.Errorf("found no log level in chainlink log line: %s", jsonLogLine)
	}

	if logLevel == "crit" { // "crit" is a custom core type they map to DPanic
		zapLevel = zapcore.DPanicLevel
	} else {
		zapLevel, err = zapcore.ParseLevel(logLevel)
		if err != nil {
			return 0, fmt.Errorf("'%s' not a valid zapcore level", logLevel)
		}
	}

	if zapLevel >= failingLogLevel {
		logErr := fmt.Errorf("%s '%s', failing any log level higher than %s: %s", OneLogAtLogLevelErr, logLevel, zapLevel.String(), jsonLogLine)
		if failureThreshold > 1 {
			logErr = fmt.Errorf("%s '%s' or above; failure threshold of %d reached; last error found: %s", MultipleLogsAtLogLevelErr, logLevel, failureThreshold, jsonLogLine)
		}
		logMessage, hasMessage := jsonMapping["msg"]
		if !hasMessage {
			foundSoFar++
			if foundSoFar >= failureThreshold {
				return foundSoFar, logErr
			}
			return foundSoFar, nil
		}

		for _, allowedLog := range allowedMessages {
			if strings.Contains(logMessage.(string), allowedLog.message) {
				if allowedLog.logWhenFound {
					log.Warn().
						Str("Reason", allowedLog.reason).
						Str("Level", allowedLog.level.CapitalString()).
						Str("Msg", logMessage.(string)).
						Msg("Found allowed log message, ignoring")
				}

				return foundSoFar, nil
			}
		}

		foundSoFar++
		if foundSoFar >= failureThreshold {
			return foundSoFar, logErr
		}
	}

	return foundSoFar, nil
}
