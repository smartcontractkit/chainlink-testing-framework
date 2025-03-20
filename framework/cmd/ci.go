package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/v50/github"
	"github.com/google/uuid"
	"go.uber.org/ratelimit"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	JobsRateLimitPerSecond = 20
	GHResultsPerPage       = 100 // anything above that won't work, check GitHub docs
)

const (
	MaxNameLen     = 120
	SuccessEmoji   = "✅"
	FailureEmoji   = "❌"
	RerunEmoji     = "🔄"
	CancelledEmoji = "🚫"
)

var (
	SlowTestThreshold          = 5 * time.Minute
	ExtremelySlowTestThreshold = 10 * time.Minute

	DebugDirRoot       = "ctf-ci-debug"
	DebugSubDirWF      = filepath.Join(DebugDirRoot, "workflows")
	DebugSubDirJobs    = filepath.Join(DebugDirRoot, "jobs")
	DefaultResultsDir  = "."
	DefaultResultsFile = "ctf-ci"
)

type GitHubActionsClient interface {
	ListRepositoryWorkflowRuns(ctx context.Context, owner, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *github.ListWorkflowJobsOptions) (*github.Jobs, *github.Response, error)
}

type AnalysisConfig struct {
	Debug               bool      `json:"debug"`
	Owner               string    `json:"owner"`
	Repo                string    `json:"repo"`
	WorkflowName        string    `json:"workflow_name"`
	TimeDaysBeforeStart int       `json:"time_days_before_start"`
	TimeStart           time.Time `json:"time_start"`
	TimeDaysBeforeEnd   int       `json:"time_days_before_end"`
	TimeEnd             time.Time `json:"time_end"`
	Typ                 string    `json:"type"`
	ResultsFile         string    `json:"results_file"`
}

type Stat struct {
	Name          string          `json:"name"`
	Successes     int             `json:"successes"`
	Failures      int             `json:"failures"`
	Cancels       int             `json:"cancels"`
	ReRuns        int             `json:"reRuns"`
	P50           time.Duration   `json:"p50"`
	P95           time.Duration   `json:"p95"`
	P99           time.Duration   `json:"p99"`
	TotalDuration time.Duration   `json:"totalDuration"`
	Durations     []time.Duration `json:"-"`
}

type Stats struct {
	Mu            *sync.Mutex      `json:"-"`
	Runs          int              `json:"runs"`
	CancelledRuns int              `json:"cancelled_runs"`
	IgnoredRuns   int              `json:"ignored_runs"`
	Jobs          map[string]*Stat `json:"jobs"`
	Steps         map[string]*Stat `json:"steps"`
}

func AnalyzeJobsSteps(ctx context.Context, client GitHubActionsClient, cfg *AnalysisConfig) (*Stats, error) {
	framework.L.Info().Time("From", cfg.TimeStart).Time("To", cfg.TimeEnd).Msg("Analyzing workflow runs")
	opts := &github.ListWorkflowRunsOptions{
		Created:     fmt.Sprintf("%s..%s", cfg.TimeStart.Format(time.DateOnly), cfg.TimeEnd.Format(time.DateOnly)),
		ListOptions: github.ListOptions{PerPage: GHResultsPerPage},
	}
	rlJobs := ratelimit.New(JobsRateLimitPerSecond)
	stats := &Stats{
		Mu:    &sync.Mutex{},
		Jobs:  make(map[string]*Stat),
		Steps: make(map[string]*Stat),
	}
	refreshDebugDirs()
	for {
		runs, resp, err := client.ListRepositoryWorkflowRuns(ctx, cfg.Owner, cfg.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
		}
		framework.L.Debug().Int("Runs", len(runs.WorkflowRuns)).Msg("Loading runs")

		eg := &errgroup.Group{}
		for _, wr := range runs.WorkflowRuns {
			framework.L.Debug().Str("Name", *wr.Name).Msg("Analyzing workflow run")
			if !strings.Contains(*wr.Name, cfg.WorkflowName) {
				stats.IgnoredRuns++
				continue
			}
			stats.Runs++
			// analyze workflow
			name := *wr.Name
			framework.L.Debug().Str("Name", name).Msg("Analyzing workflow run")
			_ = dumpResults(cfg.Debug, DebugSubDirWF, name, wr)
			eg.Go(func() error {
				rlJobs.Take()
				jobs, _, err := client.ListWorkflowJobs(ctx, cfg.Owner, cfg.Repo, *wr.ID, &github.ListWorkflowJobsOptions{
					ListOptions: github.ListOptions{PerPage: GHResultsPerPage},
				})
				if err != nil {
					return err
				}
				// analyze jobs
				stats.Mu.Lock()
				defer stats.Mu.Unlock()
				for _, j := range jobs.Jobs {
					name := *j.Name
					_ = dumpResults(cfg.Debug, DebugSubDirJobs, name, wr)
					if skippedOrInProgressJob(j) {
						stats.IgnoredRuns++
						continue
					}
					if j.Status != nil && *j.Status == "cancelled" {
						stats.CancelledRuns++
						continue
					}
					dur := j.CompletedAt.Time.Sub(j.StartedAt.Time)
					if _, ok := stats.Jobs[name]; !ok {
						stats.Jobs[name] = &Stat{
							Name:          name,
							Durations:     []time.Duration{dur},
							TotalDuration: dur,
						}
					} else {
						stats.Jobs[name].Durations = append(stats.Jobs[*j.Name].Durations, dur)
						stats.Jobs[name].TotalDuration += dur
					}
					if j.RunAttempt != nil && *j.RunAttempt > 1 {
						stats.Jobs[name].ReRuns++
					}
					if j.Conclusion != nil && *j.Conclusion == "failure" {
						stats.Jobs[name].Failures++
					} else {
						stats.Jobs[name].Successes++
					}
					if j.Conclusion != nil && *j.Conclusion == "cancelled" {
						stats.Jobs[name].Cancels++
					}
					// analyze steps
					for _, s := range j.Steps {
						name := *s.Name
						if skippedOrInProgressStep(s) {
							continue
						}
						dur := s.CompletedAt.Time.Sub(s.StartedAt.Time)
						if _, ok := stats.Steps[name]; !ok {
							stats.Steps[name] = &Stat{
								Name:      name,
								Durations: []time.Duration{dur},
							}
						} else {
							stats.Steps[name].Durations = append(stats.Steps[name].Durations, dur)
						}
						if *s.Conclusion == "failure" {
							stats.Steps[name].Failures++
						} else {
							stats.Steps[name].Successes++
						}
					}
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	jobs := make([]*Stat, 0)
	for _, js := range stats.Jobs {
		stat := calculatePercentiles(js)
		jobs = append(jobs, stat)
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].P50 > jobs[j].P50 })

	steps := make([]*Stat, 0)
	for _, js := range stats.Steps {
		stat := calculatePercentiles(js)
		steps = append(steps, stat)
	}
	sort.Slice(steps, func(i, j int) bool { return steps[i].P50 > steps[j].P50 })

	switch cfg.Typ {
	case "jobs":
		for _, js := range jobs {
			printSummary(js.Name, js, cfg.Typ)
		}
	case "steps":
		for _, js := range steps {
			printSummary(js.Name, js, cfg.Typ)
		}
	default:
		return nil, errors.New("analytics type is not recognized")
	}
	framework.L.Info().
		Int("Runs", stats.Runs).
		Int("Cancelled", stats.CancelledRuns).
		Int("Ignored", stats.IgnoredRuns).
		Msg("Total runs analyzed")
	return stats, nil
}

func AnalyzeCIRuns(cfg *AnalysisConfig) (*Stats, error) {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	framework.L.Info().
		Str("Owner", cfg.Owner).
		Str("Repo", cfg.Repo).
		Str("Workflow", cfg.WorkflowName).
		Msg("Analyzing CI runs")

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Fetch workflow runs for the last N days
	// have GH rate limits in mind, see file constants
	timeStart := time.Now().AddDate(0, 0, -cfg.TimeDaysBeforeStart)
	cfg.TimeStart = timeStart
	timeEnd := time.Now().AddDate(0, 0, -cfg.TimeDaysBeforeEnd)
	cfg.TimeEnd = timeEnd
	stats, err := AnalyzeJobsSteps(ctx, client.Actions, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
	}
	if cfg.ResultsFile != "" {
		_ = dumpResults(true, DefaultResultsDir, DefaultResultsFile, stats)
	}
	return stats, nil
}

// getColor returns a color printer based on the duration
func getColor(duration time.Duration) *color.Color {
	switch {
	case duration < SlowTestThreshold:
		return color.New(color.FgGreen)
	case duration < ExtremelySlowTestThreshold:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgRed)
	}
}

func printSummary(name string, s *Stat, typ string) {
	cp50 := getColor(s.P50)
	cp95 := getColor(s.P95)
	cp99 := getColor(s.P99)
	var cpFlaky *color.Color
	if s.ReRuns > 0 {
		cpFlaky = color.New(color.FgRed)
	} else {
		cpFlaky = color.New(color.FgGreen)
	}
	if len(name) > MaxNameLen {
		name = name[:MaxNameLen]
	}
	switch typ {
	case "jobs":
		fmt.Printf("%s 50th:%-10s 95th:%-10s 99th:%-10s Total:%-10s %s %-2s %s %-2s %s %-2s %s %-2s\n",
			cp50.Sprintf("%-120s", name),
			cp50.Sprintf("%-8s", s.P50),
			cp95.Sprintf("%-8s", s.P95),
			cp99.Sprintf("%-8s", s.P99),
			fmt.Sprintf("%-8s", s.TotalDuration.Round(time.Second)),
			RerunEmoji,
			cpFlaky.Sprintf("%-2d", s.ReRuns),
			FailureEmoji,
			cpFlaky.Sprintf("%-2d", s.Failures),
			SuccessEmoji,
			cpFlaky.Sprintf("%-2d", s.Successes),
			CancelledEmoji,
			cpFlaky.Sprintf("%-2d", s.Cancels),
		)
	case "steps":
		fmt.Printf("%s 50th:%-10s 95th:%-10s 99th:%-10s %s %-2s %s %-2s\n",
			cp50.Sprintf("%-120s", name),
			cp50.Sprintf("%-8s", s.P50),
			cp95.Sprintf("%-8s", s.P95),
			cp99.Sprintf("%-8s", s.P99),
			FailureEmoji,
			cpFlaky.Sprintf("%-2d", s.Failures),
			SuccessEmoji,
			cpFlaky.Sprintf("%-2d", s.Successes),
		)
	}
}

func skippedOrInProgressStep(s *github.TaskStep) bool {
	if s.Conclusion == nil || s.CompletedAt == nil || *s.Conclusion == "skipped" {
		return true
	}
	return false
}

func skippedOrInProgressJob(s *github.WorkflowJob) bool {
	if s.Conclusion == nil || s.CompletedAt == nil || (*s.Conclusion == "skipped") {
		return true
	}
	return false
}

func calculatePercentiles(stat *Stat) *Stat {
	sort.Slice(stat.Durations, func(i, j int) bool { return stat.Durations[i] < stat.Durations[j] })
	q := func(d *Stat, quantile float64) int {
		return int(float64(len(d.Durations)) * quantile / 100)
	}
	stat.P50 = stat.Durations[q(stat, 50)].Round(time.Second)
	stat.P95 = stat.Durations[q(stat, 95)].Round(time.Second)
	stat.P99 = stat.Durations[q(stat, 99)].Round(time.Second)
	return stat
}

func refreshDebugDirs() {
	_ = os.RemoveAll(DebugDirRoot)
	if _, err := os.Stat(DebugDirRoot); os.IsNotExist(err) {
		_ = os.MkdirAll(DebugSubDirWF, os.ModePerm)
		_ = os.MkdirAll(DebugSubDirJobs, os.ModePerm)
	}
}

func dumpResults(enabled bool, dir, name string, data interface{}) error {
	if enabled {
		d, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}
		return os.WriteFile(fmt.Sprintf("%s/%s-%s.json", dir, name, uuid.NewString()[0:5]), d, os.ModeAppend|os.ModePerm)
	}
	return nil
}
