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
	l.logStream.consumerMutex.RLock()
	for _, c := range l.logStream.consumers {
		if c.name == containerName {
			consumer = c
			break
		}
	}
	l.logStream.consumerMutex.RUnlock()

	if consumer == nil {
		return new(ReturnType), fmt.Errorf("no consumer found for container %s", containerName)
	}

	// Create a temporary snapshot of the log file and temporarily lock accept mutex to prevent new logs from being written
	// as that might corrupt the gob file due to saving of incomplete logs
	l.logStream.acceptMutex.Lock()
	tempSnapshotFile, err := createTemporarySnapshot(consumer.tempFile)
	l.logStream.acceptMutex.Unlock()
	if err != nil {
		return new(ReturnType), err
	}
	defer func() { _ = os.Remove(tempSnapshotFile.Name()) }()

	decoder := gob.NewDecoder(tempSnapshotFile)
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

func createTemporarySnapshot(file *os.File) (*os.File, error) {
	// Duplicate the file descriptor (so that when we work with the file, we don't affect the original file descriptor, especially cursor position)
	fd, err := syscall.Dup(int(file.Fd()))
	if err != nil {
		return nil, err
	}
	// We are not creating a new file here, but creating a new file descriptor, but filename is still required
	readFile := os.NewFile(uintptr(fd), "snapshot.txt")

	// Move the cursor of the duplicated file descriptor to the beginning as otherwise, the file will be read from the current cursor position, which is at the end of the file
	if _, err := readFile.Seek(0, 0); err != nil {
		return nil, err
	}

	tempSnapshot, err := os.CreateTemp("", "snapshot")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(tempSnapshot, readFile); err != nil {
		return nil, err
	}

	snapshotStat, err := tempSnapshot.Stat()
	if err != nil {
		return nil, err
	}
	if snapshotStat.Size() == 0 {
		return nil, fmt.Errorf("temporary log snapshot is empty")
	}

	// Compare the snapshot size with the original file size to make sure everything was copied
	// and nothing was added in the meantime
	originalStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if snapshotStat.Size() != originalStat.Size() {
		return nil, fmt.Errorf("temporary log snapshot size (%d) does not match original log file size (%d)", snapshotStat.Size(), originalStat.Size())
	}

	// Move cursor to the beginning of the temporary snapshot, so that it will be read from the beginning
	if _, err := tempSnapshot.Seek(0, 0); err != nil {
		return nil, err
	}

	return tempSnapshot, nil
}
