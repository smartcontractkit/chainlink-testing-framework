package framework

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestFanoutContainerLogsReplicatesBytesToAllConsumers(t *testing.T) {
	containerName := "container-a"
	payload := dockerMuxPayload("line-1\n", "line-2\n")
	logStream := map[string]io.ReadCloser{
		containerName: io.NopCloser(bytes.NewReader(payload)),
	}

	type consumeResult struct {
		name string
		data []byte
	}
	var results []consumeResult
	var resultsMu sync.Mutex

	consumers := []LogStreamConsumer{
		{
			Name: "consumer-1",
			Consume: func(streams map[string]io.ReadCloser) error {
				data, err := io.ReadAll(streams[containerName])
				if err != nil {
					return err
				}
				resultsMu.Lock()
				results = append(results, consumeResult{name: "consumer-1", data: data})
				resultsMu.Unlock()
				return nil
			},
		},
		{
			Name: "consumer-2",
			Consume: func(streams map[string]io.ReadCloser) error {
				data, err := io.ReadAll(streams[containerName])
				if err != nil {
					return err
				}
				resultsMu.Lock()
				results = append(results, consumeResult{name: "consumer-2", data: data})
				resultsMu.Unlock()
				return nil
			},
		},
	}

	err := fanoutContainerLogs(logStream, consumers...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 consumer results, got: %d", len(results))
	}

	slices.SortFunc(results, func(a, b consumeResult) int {
		if a.name < b.name {
			return -1
		}
		if a.name > b.name {
			return 1
		}
		return 0
	})
	for _, result := range results {
		if !bytes.Equal(result.data, payload) {
			t.Fatalf("consumer %s did not receive full payload", result.name)
		}
	}
}

func TestSaveContainerLogsFromStreams(t *testing.T) {
	tDir := t.TempDir()
	logStreams := map[string]io.ReadCloser{
		"node-a": io.NopCloser(bytes.NewReader(dockerMuxPayload("a-1\n", "a-2\n"))),
		"node-b": io.NopCloser(bytes.NewReader(dockerMuxPayload("b-1\n"))),
		"node-c": io.NopCloser(bytes.NewReader(dockerMuxPayload("c-1\n", "c-2\n", "c-3\n"))),
	}

	paths, err := SaveContainerLogsFromStreams(tDir, logStreams)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(paths) != len(logStreams) {
		t.Fatalf("expected %d files, got %d", len(logStreams), len(paths))
	}

	for _, p := range paths {
		content, readErr := os.ReadFile(filepath.Clean(p))
		if readErr != nil {
			t.Fatalf("failed to read log file %s: %v", p, readErr)
		}
		if strings.TrimSpace(string(content)) == "" {
			t.Fatalf("expected non-empty log file at %s", p)
		}
	}
}

func TestPrintFailedContainerLogsFromStreams(t *testing.T) {
	logStreams := map[string]io.ReadCloser{
		"node-a": io.NopCloser(bytes.NewReader(dockerMuxPayload("error line\n"))),
		"node-b": io.NopCloser(bytes.NewReader(dockerMuxPayload("warn line\n", "trace line\n"))),
	}

	if err := PrintFailedContainerLogsFromStreams(logStreams, 30); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func dockerMuxPayload(lines ...string) []byte {
	var out bytes.Buffer
	for _, line := range lines {
		msg := []byte(line)
		header := make([]byte, 8)
		header[0] = 1
		binary.BigEndian.PutUint32(header[4:], uint32(len(msg)))
		out.Write(header)
		out.Write(msg)
	}
	return out.Bytes()
}
