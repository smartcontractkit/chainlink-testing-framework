package logstream

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

type LogProcessor[ReturnType any] struct {
	logStream *LogStream
}

// NewLogProcessor creates a new LogProcessor instance
func NewLogProcessor[ReturnType any](logStream *LogStream) *LogProcessor[ReturnType] {
	return &LogProcessor[ReturnType]{
		logStream: logStream,
	}
}

type LogProcessorFn[ReturnType any] func(content LogContent, returnValue *ReturnType) error

// ProcessContainerLogs reads the logs of a container and processes them with the provided function
func (l *LogProcessor[ReturnType]) ProcessContainerLogs(containerName string, processFn func(content LogContent, returnValue *ReturnType) error) (*ReturnType, error) {
	containerName = strings.Replace(containerName, "/", "", 1)
	var consumer *ContainerLogConsumer
	for _, c := range l.logStream.consumers {
		if c.name == containerName {
			consumer = c
			break
		}
	}

	if consumer == nil {
		return new(ReturnType), fmt.Errorf("no consumer found for container %s", containerName)
	}

	// Duplicate the file descriptor for independent reading, so we don't mess up writing the file by moving the cursor
	fd, err := syscall.Dup(int(consumer.tempFile.Fd()))
	if err != nil {
		return new(ReturnType), err
	}
	readFile := os.NewFile(uintptr(fd), "name_doesnt_matter.txt")
	_, err = readFile.Seek(0, 0)
	if err != nil {
		return new(ReturnType), err
	}

	decoder := gob.NewDecoder(readFile)
	var returnValue ReturnType

	for {
		var log LogContent
		decodeErr := decoder.Decode(&log)
		if decodeErr == nil {
			processErr := processFn(log, &returnValue)
			if processErr != nil {
				l.logStream.log.Error().
					Err(processErr).
					Str("Container", consumer.name).
					Msg("Failed to process log")
				return new(ReturnType), processErr
			}
		} else if errors.Is(decodeErr, io.EOF) {
			l.logStream.log.Debug().
				Str("Container", consumer.name).
				Str("Processing result", fmt.Sprint(returnValue)).
				Msg("Finished scanning logs")
			break
		} else {
			l.logStream.log.Error().
				Err(decodeErr).
				Str("Container", consumer.name).
				Msg("Failed to decode log")
			return new(ReturnType), decodeErr
		}
	}

	return &returnValue, nil
}

// GetRegexMatchingProcessor creates a LogProcessor that counts the number of logs matching a regex pattern. Function returns
// the LogProcessor, the processing function, and an error if the regex pattern is invalid.
func GetRegexMatchingProcessor(logStream *LogStream, pattern string) (*LogProcessor[int], LogProcessorFn[int], error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, nil, err
	}

	logProcessor := NewLogProcessor[int](logStream)

	processFn := func(logContent LogContent, returnValue *int) error {
		if re.MatchString(string(logContent.Content)) {
			newVal := *returnValue + 1
			*returnValue = newVal
		}
		return nil
	}

	return logProcessor, processFn, nil
}
