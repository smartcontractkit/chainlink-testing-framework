package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func createMockResponse(totalCount, jobCount int) string {
	jobs := make([]Job, jobCount)
	for i := 0; i < jobCount; i++ {
		jobs[i] = Job{
			Name: fmt.Sprintf("Test Job %d", i+1),
			Steps: []Step{
				{Name: "Run Tests", Conclusion: "success"},
			},
			URL: fmt.Sprintf("http://example.com/job%d", i+1),
		}
	}
	response := GitHubResponse{
		TotalCount: totalCount,
		Jobs:       jobs,
	}
	respJSON, _ := json.Marshal(response)
	return string(respJSON)
}

func TestFetchGitHubJobs(t *testing.T) {
	tests := []struct {
		name       string
		totalCount int
		jobCount   int
		pageCount  int
		wantJobs   int
		wantErr    bool
		malformed  bool
	}{
		{
			name:       "Single Page",
			totalCount: 2,
			jobCount:   2,
			pageCount:  1,
			wantJobs:   2,
			wantErr:    false,
		},
		{
			name:       "Multiple Pages",
			totalCount: 300,
			jobCount:   100,
			pageCount:  3,
			wantJobs:   300,
			wantErr:    false,
		},
		{
			name:       "Empty Response",
			totalCount: 0,
			jobCount:   0,
			pageCount:  1,
			wantJobs:   0,
			wantErr:    false,
		},
		{
			name:       "Error Response",
			totalCount: 0,
			jobCount:   0,
			pageCount:  1,
			wantJobs:   0,
			wantErr:    true,
		},
		{
			name:       "Malformed JSON Response",
			totalCount: 0,
			jobCount:   0,
			pageCount:  1,
			wantJobs:   0,
			wantErr:    true,
			malformed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageCounter := 0
			client := &MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					pageCounter++
					var statusCode int
					var body string
					if tt.wantErr {
						statusCode = http.StatusInternalServerError
						body = `{ "message": "something went wrong" }`
					} else if tt.malformed {
						statusCode = http.StatusOK
						body = `{"total_count": 1, "jobs": [`
					} else {
						statusCode = http.StatusOK
						if pageCounter <= tt.pageCount {
							body = createMockResponse(tt.totalCount, tt.jobCount)
						} else {
							body = createMockResponse(tt.totalCount, 0)
						}
					}
					r := io.NopCloser(strings.NewReader(body))
					return &http.Response{
						StatusCode: statusCode,
						Body:       r,
						Header:     make(http.Header),
					}, nil
				},
			}

			apiURL := "https://api.github.com/repos/owner/repo/actions/runs/1/jobs?per_page=100"
			jobs, err := fetchGitHubJobs(apiURL, "dummy_token", client)

			if (err != nil) != tt.wantErr {
				t.Errorf("fetchGitHubJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(jobs) != tt.wantJobs {
				t.Errorf("fetchGitHubJobs() got %d jobs, want %d", len(jobs), tt.wantJobs)
			}
			if pageCounter != tt.pageCount {
				t.Errorf("fetchGitHubJobs() fetched %d pages, want %d", pageCounter, tt.pageCount)
			}
		})
	}
}

func TestParseJobs(t *testing.T) {
	tests := []struct {
		name            string
		jobNameRegex    string
		mockJobs        []Job
		expectedResults []ParsedResult
	}{
		{
			name:         "Matching Regex",
			jobNameRegex: "Test Job (\\d)",
			mockJobs: []Job{
				{
					Name: "Test Job 1",
					Steps: []Step{
						{Name: "Run Tests", Conclusion: "success"},
					},
					URL: "http://example.com/job1",
				},
				{
					Name: "Test Job 2",
					Steps: []Step{
						{Name: "Run Tests", Conclusion: "failure"},
					},
					URL: "http://example.com/job2",
				},
			},
			expectedResults: []ParsedResult{
				{
					Conclusion: ":white_check_mark:",
					Cap:        "1",
					URL:        "http://example.com/job1",
				},
				{
					Conclusion: ":x:",
					Cap:        "2",
					URL:        "http://example.com/job2",
				},
			},
		},
		{
			name:         "Non-Matching Regex",
			jobNameRegex: "NonMatchingJob (\\d)",
			mockJobs: []Job{
				{
					Name: "Test Job 1",
					Steps: []Step{
						{Name: "Run Tests", Conclusion: "success"},
					},
					URL: "http://example.com/job1",
				},
				{
					Name: "Test Job 2",
					Steps: []Step{
						{Name: "Run Tests", Conclusion: "failure"},
					},
					URL: "http://example.com/job2",
				},
			},
			expectedResults: []ParsedResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			one := "1"
			parsedResults, err := parseResults(&tt.jobNameRegex, &one, tt.mockJobs)
			if len(tt.expectedResults) == 0 {
				require.Error(t, err)
				require.Equal(t, fmt.Sprintf("No results found for '%s' regex in workflow id 1\n", tt.jobNameRegex), err.Error())
			} else {
				require.NoError(t, err)
			}

			for i, result := range parsedResults {
				if result != tt.expectedResults[i] {
					t.Errorf("Expected result %+v, got %+v", tt.expectedResults[i], result)
				}
			}
		})
	}
}

func TestMainOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{
			"total_count": 2,
			"jobs": [
				{
					"name": "Test Job 1",
					"steps": [{"name": "Run Tests", "conclusion": "success"}],
					"html_url": "http://example.com/job1"
				},
				{
					"name": "Test Job 2",
					"steps": [{"name": "Run Tests", "conclusion": "failure"}],
					"html_url": "http://example.com/job2"
				}
			]
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	jobs, err := fetchGitHubJobs(server.URL, "dummy_token", client)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("Expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].Name != "Test Job 1" {
		t.Fatalf("Expected job name 'Test Job 1', got %s", jobs[0].Name)
	}
}

func TestMainFunction(t *testing.T) {
	// Backup original arguments and restore them after the test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Mock the necessary flags
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "Missing Required Flags",
			args:       []string{"cmd"},
			wantOutput: "Please provide all required flags: --githubToken, --githubRepo, --workflowRunID, --jobNameRegex",
			wantErr:    true,
		},
		{
			name:    "Valid Flags",
			args:    []string{"cmd", "--githubToken=dummy_token", "--githubRepo=owner/repo", "--workflowRunID=1", "--jobNameRegex=Test Job (\\d)"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			var buf bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &buf)
			stdout := os.Stdout
			stderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			done := make(chan bool)
			go func() {
				_, _ = io.Copy(mw, r)
				done <- true
			}()

			defer func() {
				_ = w.Close()
				<-done
				os.Stdout = stdout
				os.Stderr = stderr
			}()

			// Run the main function and capture its output
			if tt.wantErr {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok {
							if !strings.Contains(err.Error(), tt.wantOutput) && tt.name == "Missing Required Flags" {
								t.Errorf("Expected error message: %s, got: %s", tt.wantOutput, err.Error())
							} else if tt.name == "Valid Flags" && !strings.Contains(err.Error(), "401 Unauthorized") {
								t.Errorf("Expected error message: 401 Unauthorized, got: %s", err.Error())
							}
						} else {
							t.Errorf("Expected an error, got: %v", r)
						}
					}
				}()
				main()
			} else {
				main()
			}

			// Check for expected output for Missing Required Flags
			if tt.name == "Missing Required Flags" && !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("Expected output: %s, got: %s", tt.wantOutput, buf.String())
			}
		})
	}
}
