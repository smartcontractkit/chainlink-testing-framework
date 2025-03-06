package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/v50/github"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"golang.org/x/oauth2"
)

const (
	WorkflowRateLimitPerSecond = 10
	JobsRateLimitPerSecond     = 10
	MaxBarLength               = 50
	GHResultsPerPage           = 100 // anything above that won't work
)

var (
	SlowTestThreshold          = 5 * time.Minute
	ExtremelySlowTestThreshold = 10 * time.Minute
)

type JobResult struct {
	StepStats map[string]Stat
	JobStats  map[string]Stat
}

type Stat struct {
	Name      string
	Median    time.Duration
	P95       time.Duration
	P99       time.Duration
	Durations []time.Duration
}

// AnalyzeCIRuns analyzes GitHub Actions job runs and prints statistics
func AnalyzeCIRuns(owner, repo, wf string, daysRange int) error {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	framework.L.Info().
		Str("Owner", owner).
		Str("Repo", repo).
		Str("Workflow", wf).
		Msg("Analyzing CI runs")

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Fetch workflow runs for the last N days
	// have GH rate limits in mind, see file constants
	lastMonth := time.Now().AddDate(0, 0, -daysRange)
	runs, err := getAllWorkflowRuns(ctx, client, owner, repo, wf, lastMonth)
	if err != nil {
		return fmt.Errorf("failed to fetch workflow runs: %w", err)
	}

	framework.L.Info().
		Int("Runs", len(runs)).
		Msg("Found matching workflow runs")

	results := make(chan JobResult, len(runs))
	eg := &errgroup.Group{}
	rl := ratelimit.New(JobsRateLimitPerSecond)

	for _, run := range runs {
		eg.Go(func() error {
			rl.Take()
			return analyzeRun(ctx, client, run, results, owner, repo)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	close(results)

	perStepStats := make(map[string]Stat)
	perJobStats := make(map[string]Stat)

	for result := range results {
		// Aggregate step durations
		for stepName, durations := range result.StepStats {
			if existing, ok := perStepStats[stepName]; ok {
				existing.Durations = append(existing.Durations, durations.Durations...)
				perStepStats[stepName] = existing
			} else {
				perStepStats[stepName] = Stat{
					Name:      stepName,
					Durations: durations.Durations,
				}
			}
		}
		// Aggregate job stats
		for jobName, stat := range result.JobStats {
			if existing, ok := perJobStats[jobName]; ok {
				existing.Durations = append(existing.Durations, stat.Durations...)
				perJobStats[jobName] = existing
			} else {
				perJobStats[jobName] = Stat{
					Name:      jobName,
					Durations: stat.Durations,
				}
			}
		}
	}

	for stepName, stat := range perStepStats {
		stat.Median, stat.P95, stat.P99 = calculatePercentiles(stat.Durations)
		perStepStats[stepName] = stat
	}
	for jobName, stat := range perJobStats {
		stat.Median, stat.P95, stat.P99 = calculatePercentiles(stat.Durations)
		perJobStats[jobName] = stat
	}
	fmt.Print("\nSteps:\n")
	printStats(perStepStats)
	fmt.Print("\nJobs:\n")
	printStats(perJobStats)
	return nil
}

func getAllWorkflowRuns(ctx context.Context, client *github.Client, owner, repo, name string, timeRange time.Time) ([]*github.WorkflowRun, error) {
	var allRuns []*github.WorkflowRun
	opts := &github.ListWorkflowRunsOptions{
		Created:     fmt.Sprintf(">%s", timeRange.Format(time.RFC3339)),
		ListOptions: github.ListOptions{PerPage: 100},
	}
	rl := ratelimit.New(WorkflowRateLimitPerSecond)
	for {
		rl.Take()
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
		}
		framework.L.Debug().Int("Runs", len(runs.WorkflowRuns)).Msg("Loading runs")
		for _, wr := range runs.WorkflowRuns {
			if strings.Contains(*wr.Name, name) {
				allRuns = append(allRuns, wr)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRuns, nil
}

// analyzeRun fetches workflow runs that are not skipped and returns their Stat through channel
func analyzeRun(ctx context.Context, client *github.Client, run *github.WorkflowRun, results chan<- JobResult, owner, repo string) error {
	logger := framework.L.With().
		Str("RunID", fmt.Sprintf("%d", *run.ID)).
		Str("CreatedAt", run.CreatedAt.Format(time.RFC3339)).
		Logger()
	logger.Debug().Msg("Analyzing run")

	jobs, _, err := client.Actions.ListWorkflowJobs(ctx, owner, repo, *run.ID, &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{PerPage: GHResultsPerPage},
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch jobs for run")
	}

	stepStats := make(map[string]Stat)
	jobStats := make(map[string]Stat)

	// Analyze each job
	for _, job := range jobs.Jobs {
		logger.Debug().
			Str("job_id", fmt.Sprintf("%d", *job.ID)).
			Str("job_name", *job.Name).
			Msg("Found job")

		// ignore jobs that are in progress or skipped
		if job.Conclusion != nil && *job.Conclusion == "skipped" {
			continue
		}
		if job.CompletedAt == nil {
			continue
		}
		jobDuration := job.CompletedAt.Time.Sub(job.StartedAt.Time)
		// Collect step durations
		for _, step := range job.Steps {
			if step.Conclusion != nil && *step.Conclusion == "skipped" {
				continue
			}
			elapsed := step.CompletedAt.Time.Sub(step.StartedAt.Time)
			if existing, ok := stepStats[*step.Name]; ok {
				existing.Durations = append(existing.Durations, elapsed)
				stepStats[*step.Name] = existing
			} else {
				stepStats[*step.Name] = Stat{
					Name:      *step.Name,
					Durations: []time.Duration{elapsed},
				}
			}
		}
		// Collect per-job statistics
		if existing, ok := jobStats[*job.Name]; ok {
			existing.Durations = append(existing.Durations, jobDuration)
			jobStats[*job.Name] = existing
		} else {
			jobStats[*job.Name] = Stat{
				Name:      *job.Name,
				Durations: []time.Duration{jobDuration},
			}
		}
	}
	results <- JobResult{
		StepStats: stepStats,
		JobStats:  jobStats,
	}
	return nil
}

// calculatePercentiles calculates the median (50th), 95th, and 99th percentiles
func calculatePercentiles(durations []time.Duration) (median, p95, p99 time.Duration) {
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	medianIndex := int(float64(len(durations)) * 50 / 100)
	p95Index := int(float64(len(durations)) * 95 / 100)
	p99Index := int(float64(len(durations)) * 99 / 100)
	return durations[medianIndex], durations[p95Index], durations[p99Index]
}

func printStats(jobStats map[string]Stat) {
	var stats []Stat
	for _, stat := range jobStats {
		sort.Slice(stat.Durations, func(i, j int) bool { return stat.Durations[i] < stat.Durations[j] })
		stats = append(stats, stat)
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Median > stats[j].Median })
	maxNameLen := 0
	for _, stat := range stats {
		if len(stat.Name) > maxNameLen {
			maxNameLen = len(stat.Name)
		}
	}

	for _, stat := range stats {
		colorPrinter := getColorPrinter(stat.Median)
		barLength := int(stat.Median.Seconds())
		if barLength > MaxBarLength {
			barLength = MaxBarLength
		}
		bar := strings.Repeat("=", barLength)
		fmt.Printf("%-*s 50th:%s 95th:%s 99th:%s %s\n",
			maxNameLen,
			stat.Name,
			colorPrinter.Sprintf("%-12s", stat.Median.Round(time.Second)),
			colorPrinter.Sprintf("%-12s", stat.P95.Round(time.Second)),
			colorPrinter.Sprintf("%-12s", stat.P99.Round(time.Second)),
			colorPrinter.Sprint(bar))
	}
}

// getColorPrinter returns a color printer based on the duration
func getColorPrinter(duration time.Duration) *color.Color {
	switch {
	case duration < SlowTestThreshold:
		return color.New(color.FgGreen)
	case duration < ExtremelySlowTestThreshold:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgRed)
	}
}
