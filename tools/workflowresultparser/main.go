package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Job struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
	URL   string `json:"html_url"`
}

type Step struct {
	Name       string `json:"name"`
	Conclusion string `json:"conclusion"`
}

type GitHubResponse struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

type ParsedResult struct {
	Conclusion string `json:"conclusion"`
	Cap        string `json:"cap"`
	URL        string `json:"html_url"`
}

type ResultsMap map[string][]ParsedResult

func main() {
	githubToken := flag.String("githubToken", "", "GitHub token for authentication")
	githubRepo := flag.String("githubRepo", "", "GitHub repository in the format owner/repo")
	workflowRunID := flag.String("workflowRunID", "", "ID of the GitHub Actions workflow run")
	jobNameRegex := flag.String("jobNameRegex", "", "Regex pattern to match job names")
	namedKey := flag.String("namedKey", "", "Optional named key under which results will be stored")
	outputFile := flag.String("outputFile", "", "Optional output file to save results")

	flag.Parse()

	if *githubToken == "" || *githubRepo == "" || *workflowRunID == "" || *jobNameRegex == "" {
		panic(fmt.Errorf("Please provide all required flags: --githubToken, --githubRepo, --workflowRunID, --jobNameRegex"))
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs/%s/jobs?per_page=100", *githubRepo, *workflowRunID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(fmt.Errorf("error creating HTTP request:", err))
	}
	req.Header.Set("Authorization", "Bearer "+*githubToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("error making HTTP request:", err))
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("GitHub API request failed with status:", resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("error reading response body:", err))
	}

	var githubResponse GitHubResponse
	err = json.Unmarshal(body, &githubResponse)
	if err != nil {
		panic(fmt.Errorf("error unmarshalling JSON response:", err))
	}

	var parsedResults []ParsedResult
	re := regexp.MustCompile(*jobNameRegex)
	for _, job := range githubResponse.Jobs {
		if re.MatchString(job.Name) {
			for _, step := range job.Steps {
				if step.Name == "Run Tests" {
					conclusion := ":x:"
					if step.Conclusion == "success" {
						conclusion = ":white_check_mark:"
					}
					captureGroup := fmt.Sprintf("%s", re.FindStringSubmatch(job.Name)[1])
					parsedResults = append(parsedResults, ParsedResult{
						Conclusion: conclusion,
						Cap:        captureGroup,
						URL:        job.URL,
					})
				}
			}
		}
	}

	if len(parsedResults) == 0 {
		fmt.Printf("No results found for '%s' regex in workflow id %s\n", *jobNameRegex, *workflowRunID)
		return
	}

	results := ResultsMap{}

	if *outputFile != "" {
		if _, statErr := os.Stat(*outputFile); statErr == nil {
			existingData, readErr := os.ReadFile(*outputFile)
			if readErr == nil {
				jsonErr := json.Unmarshal(existingData, &results)
				if jsonErr != nil {
					panic(fmt.Errorf("error unmarshalling existing data:", jsonErr))
				}
			}
		}
	}

	key := "results"
	if *namedKey != "" {
		key = *namedKey
	}
	results[key] = append(results[key], parsedResults...)

	formattedResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(fmt.Errorf("error marshalling formatted results:", err))
	}

	if *outputFile != "" {
		err = os.WriteFile(*outputFile, formattedResults, 0644)
		if err != nil {
			panic(fmt.Errorf("error writing results to file:", err))
		} else {
			fmt.Printf("Results for '%s' regex and workflow id %s saved to %s\n", *jobNameRegex, *workflowRunID, *outputFile)
		}
	} else {
		fmt.Println(string(formattedResults))
	}
}
