// File: internal/mock_logger.go
package internal

import (
	"fmt"
	"strings"
	"sync"
)

// MockLogger captures logs for testing purposes.
type MockLogger struct {
	mu     sync.Mutex
	Logs   []string
	Errors []string
}

// NewMockLogger initializes a new MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Logs:   []string{},
		Errors: []string{},
	}
}

func (ml *MockLogger) Reset() {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.Logs = []string{}
	ml.Errors = []string{}
}

func (ml *MockLogger) Debug(args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.Logs = append(ml.Logs, fmt.Sprint(args...))
}

func (ml *MockLogger) Info(args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.Logs = append(ml.Logs, fmt.Sprint(args...))
}

func (ml *MockLogger) Warn(args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.Logs = append(ml.Logs, fmt.Sprint(args...))
}

func (ml *MockLogger) Error(args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.Errors = append(ml.Errors, fmt.Sprint(args...))
}

// Debugf logs a formatted debug message.
func (ml *MockLogger) Debugf(format string, args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	logMsg := fmt.Sprintf(format, args...)
	ml.Logs = append(ml.Logs, logMsg)
}

// Infof logs a formatted info message.
func (ml *MockLogger) Infof(format string, args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	logMsg := fmt.Sprintf(format, args...)
	ml.Logs = append(ml.Logs, logMsg)
}

// Warnf logs a formatted warning message.
func (ml *MockLogger) Warnf(format string, args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	logMsg := fmt.Sprintf(format, args...)
	ml.Logs = append(ml.Logs, logMsg)
}

// Errorf logs a formatted error message.
func (ml *MockLogger) Errorf(format string, args ...interface{}) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	logMsg := fmt.Sprintf(format, args...)
	ml.Errors = append(ml.Errors, logMsg)
}

// ContainsLog checks if any log in Logs contains the specified substring.
func (ml *MockLogger) ContainsLog(substring string) bool {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	for _, log := range ml.Logs {
		if strings.Contains(log, substring) {
			return true
		}
	}
	return false
}

// ContainsLog checks if any log in Logs contains the specified substring.
func (ml *MockLogger) ContainsError(substring string) bool {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	for _, error := range ml.Errors {
		if strings.Contains(error, substring) {
			return true
		}
	}
	return false
}

// NumLogs returns the total number of logs.
func (ml *MockLogger) NumLogs() int {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	return len(ml.Logs)
}
