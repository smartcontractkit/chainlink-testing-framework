package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

var L = framework.L

// LokiPushRequest represents the payload format expected by Loki's push API.
type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream represents one log stream.
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

const (
	lokiURL     = "http://localhost:3030/loki/api/v1/push"
	grafanaURL  = "http://localhost:3000/explore?panes=%7B%22V0P%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22"
	grafanaURL2 = "%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1"
)

// processAndUploadDir traverses the given directory recursively and
// processes every file (ignoring directories) by calling processAndUploadLog.
func processAndUploadDir(dirPath string, limiter ratelimit.Limiter, chunks int, jobID string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			L.Error().Err(err).Msgf("Error accessing %s", path)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		L.Info().Msgf("Processing file: %s", path)
		f, err := os.Open(path)
		if err != nil {
			L.Error().Err(err).Msgf("Error opening file %s", path)
			return nil
		}
		defer f.Close()

		if err := processAndUploadLog(path, f, limiter, chunks, jobID); err != nil {
			L.Error().Err(err).Msgf("Error processing file %s", path)
			// Continue processing other files even if one fails.
		}
		return nil
	})
}

// processAndUploadLog reads log lines from the provided reader,
// splits them into chunks (if more than 10,000 lines) and uploads each chunk concurrently.
func processAndUploadLog(source string, r io.Reader, limiter ratelimit.Limiter, chunks int, jobID string) error {
	scanner := bufio.NewScanner(r)
	var values [][2]string
	baseTime := time.Now()

	// Read all log lines; each line gets a unique timestamp.
	for scanner.Scan() {
		line := scanner.Text()
		ts := baseTime.UnixNano()
		values = append(values, [2]string{fmt.Sprintf("%d", ts), line})
		baseTime = baseTime.Add(time.Nanosecond)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning logs from %s: %w", source, err)
	}

	totalLines := len(values)
	if totalLines == 0 {
		L.Info().Msgf("No log lines found in %s", source)
		return nil
	}
	// Use one chunk if there are 10,000 or fewer lines.
	if totalLines <= 10000 {
		chunks = 1
	}
	if chunks > totalLines {
		chunks = totalLines
	}
	chunkSize := totalLines / chunks
	remainder := totalLines % chunks
	L.Debug().Int("total_lines", totalLines).
		Int("chunks", chunks).
		Msgf("Starting chunk processing for %s", source)
	var wg sync.WaitGroup
	errCh := make(chan error, chunks)
	start := 0
	for i := 0; i < chunks; i++ {
		extra := 0
		if i < remainder {
			extra = 1
		}
		end := start + chunkSize + extra
		chunkValues := values[start:end]
		startLine := start + 1
		endLine := end
		start = end

		// Use the unique jobID as the "job" label.
		labels := map[string]string{
			"job":    jobID,
			"chunk":  fmt.Sprintf("%d", i+1),
			"source": source,
		}
		reqBody := LokiPushRequest{
			Streams: []LokiStream{
				{
					Stream: labels,
					Values: chunkValues,
				},
			},
		}
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("error marshaling JSON for chunk %d: %w", i+1, err)
		}
		chunkMB := float64(len(data)) / (1024 * 1024)
		L.Debug().Int("chunk", i+1).
			Float64("chunk_size_MB", chunkMB).
			Int("start_line", startLine).
			Int("end_line", endLine).
			Msg("Prepared chunk for upload")

		wg.Add(1)
		go func(chunkNum, sLine, eLine int, payload []byte, sizeMB float64) {
			defer wg.Done()
			const maxRetries = 50
			const retryDelay = 1 * time.Second

			var resp *http.Response
			var attempt int
			var err error
			for attempt = 1; attempt <= maxRetries; attempt++ {
				limiter.Take()
				resp, err = http.Post(lokiURL, "application/json", bytes.NewReader(payload))
				if err != nil {
					L.Error().Err(err).Int("attempt", attempt).
						Int("chunk", chunkNum).
						Float64("chunk_size_MB", sizeMB).
						Msg("Error sending POST request")
					time.Sleep(retryDelay)
					continue
				}

				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				if resp.StatusCode == 429 {
					L.Debug().Int("attempt", attempt).
						Int("chunk", chunkNum).
						Float64("chunk_size_MB", sizeMB).
						Msg("Received 429, retrying...")
					time.Sleep(retryDelay)
					continue
				}

				if resp.StatusCode/100 != 2 {
					err = fmt.Errorf("loki error: %s - %s", resp.Status, body)
					L.Error().Err(err).Int("chunk", chunkNum).
						Float64("chunk_size_MB", sizeMB).
						Msg("Chunk upload failed")
					time.Sleep(retryDelay)
					continue
				}

				L.Info().Int("chunk", chunkNum).
					Float64("chunk_size_MB", sizeMB).
					Msg("Successfully uploaded chunk")
				return
			}
			errCh <- fmt.Errorf("max retries reached for chunk %d; last error: %v", chunkNum, err)
		}(i+1, startLine, endLine, data, chunkMB)
	}

	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}
