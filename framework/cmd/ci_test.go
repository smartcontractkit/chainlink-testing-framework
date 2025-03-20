package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/go-github/v50/github"
)

var (
	DefaultRunID           = github.Int64(1337)
	DefaultWorkflowRunName = "test-workflow"
	DefaultJobName         = "job"
	DefaultJobConfig       = &AnalysisConfig{
		Owner:               "test",
		Repo:                "test",
		WorkflowName:        DefaultWorkflowRunName,
		TimeDaysBeforeStart: 1,
		TimeStart:           time.Now().Add(-24 * time.Hour),
		TimeDaysBeforeEnd:   0,
		TimeEnd:             time.Now(),
		Typ:                 "jobs",
	}
	DefaultStepsConfig = &AnalysisConfig{
		Owner:               "test",
		Repo:                "test",
		WorkflowName:        DefaultWorkflowRunName,
		TimeDaysBeforeStart: 1,
		TimeStart:           time.Now().Add(-24 * time.Hour),
		TimeDaysBeforeEnd:   0,
		TimeEnd:             time.Now(),
		Typ:                 "steps",
	}
)

func TestSmokeCLIGitHubAnalytics(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *FakeGHA
		config   *AnalysisConfig
		validate func(t *testing.T, res *Stats, err error)
	}{
		{
			name:   "no jobs found for workflow",
			config: DefaultJobConfig,
			setup: func() *FakeGHA {
				c := NewFakeGHA()
				c.WorkflowRun(&github.WorkflowRun{
					ID:   DefaultRunID,
					Name: github.String(DefaultWorkflowRunName),
				})
				return c
			},
			validate: func(t *testing.T, res *Stats, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "no jobs found")
				require.Empty(t, res)
			},
		},
		{
			name:   "ignoring skipped or in progress jobs",
			config: DefaultJobConfig,
			setup: func() *FakeGHA {
				c := NewFakeGHA()
				runID1 := github.Int64(1337)
				runID2 := github.Int64(1338)
				c.WorkflowRun(&github.WorkflowRun{
					ID:   runID1,
					Name: github.String(DefaultWorkflowRunName),
				})
				c.WorkflowRun(&github.WorkflowRun{
					ID:   runID2,
					Name: github.String(DefaultWorkflowRunName),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String("job-1"),
					RunID:       runID1,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + time.Minute)},
					Conclusion:  github.String("skipped"),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String("job-2"),
					RunID:       runID2,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + time.Minute)},
					// no conclusion yet
				})
				return c
			},
			validate: func(t *testing.T, res *Stats, err error) {
				require.NoError(t, err)
				require.Equal(t, 0, len(res.Jobs))
			},
		},
		{
			name:   "successful analysis with one job",
			config: DefaultJobConfig,
			setup: func() *FakeGHA {
				c := NewFakeGHA()
				c.WorkflowRun(&github.WorkflowRun{
					ID:         DefaultRunID,
					Name:       github.String(DefaultWorkflowRunName),
					Conclusion: github.String("failure"),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String(DefaultJobName),
					RunID:       DefaultRunID,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + time.Minute)},
					Conclusion:  github.String("failure"),
				})
				return c
			},
			validate: func(t *testing.T, res *Stats, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, 1, len(res.Jobs))
				require.Equal(t, 1, res.Jobs[DefaultJobName].Failures)
				require.Equal(t, 0, res.Jobs[DefaultJobName].Successes)
				require.Equal(t, 1*time.Minute, res.Jobs[DefaultJobName].P50)
				require.Equal(t, 1*time.Minute, res.Jobs[DefaultJobName].P95)
				require.Equal(t, 1*time.Minute, res.Jobs[DefaultJobName].P99)
			},
		},
		{
			name:   "correct time frames filtering",
			config: DefaultJobConfig,
			setup: func() *FakeGHA {
				// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
				// there is no need to re-implement this logic in a mock
				// TODO: max runs per page is 100, max pages is 10
				// TODO: for precision of this tool is more than enough and after 10 pages GHA returns page 0, investigate if needed
				return NewFakeGHA()
			},
			validate: func(t *testing.T, res *Stats, err error) {},
		},
		{
			name:   "correct job quantiles for multiple workflow runs",
			config: DefaultJobConfig,
			setup: func() *FakeGHA {
				c := NewFakeGHA()
				runID1 := github.Int64(1337)
				runID2 := github.Int64(1338)
				runID3 := github.Int64(1339)
				c.WorkflowRun(&github.WorkflowRun{
					ID:   runID1,
					Name: github.String(DefaultWorkflowRunName),
				})
				c.WorkflowRun(&github.WorkflowRun{
					ID:   runID2,
					Name: github.String(DefaultWorkflowRunName),
				})
				c.WorkflowRun(&github.WorkflowRun{
					ID:   runID3,
					Name: github.String(DefaultWorkflowRunName),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String(DefaultJobName),
					RunID:       runID1,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + 1*time.Minute)},
					Conclusion:  github.String("failure"),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String(DefaultJobName),
					RunID:       runID2,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + 3*time.Minute)},
					Conclusion:  github.String("success"),
				})
				c.Job(&github.WorkflowJob{
					Name:        github.String(DefaultJobName),
					RunID:       runID3,
					StartedAt:   &github.Timestamp{Time: time.Now().Add(-time.Hour)},
					CompletedAt: &github.Timestamp{Time: time.Now().Add(-time.Hour + 5*time.Minute)},
					Conclusion:  github.String("success"),
				})
				return c
			},
			validate: func(t *testing.T, res *Stats, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, 1, len(res.Jobs))
				require.Equal(t, 1, res.Jobs[DefaultJobName].Failures)
				require.Equal(t, 2, res.Jobs[DefaultJobName].Successes)
				require.Equal(t, 3*time.Minute, res.Jobs[DefaultJobName].P50)
				require.Equal(t, 5*time.Minute, res.Jobs[DefaultJobName].P95)
				require.Equal(t, 5*time.Minute, res.Jobs[DefaultJobName].P99)
			},
		},
		{
			name:   "successful analysis for job steps",
			config: DefaultStepsConfig,
			setup: func() *FakeGHA {
				c := NewFakeGHA()
				c.WorkflowRun(&github.WorkflowRun{
					ID:         DefaultRunID,
					Name:       github.String("test-workflow"),
					Conclusion: github.String("success"),
				})
				hago := time.Now()
				c.Job(&github.WorkflowJob{
					Name:        github.String(DefaultJobName),
					RunID:       DefaultRunID,
					StartedAt:   &github.Timestamp{Time: hago},
					CompletedAt: &github.Timestamp{Time: hago.Add(10 * time.Minute)},
					Conclusion:  github.String("success"),
					Steps: []*github.TaskStep{
						{
							Name:        github.String("step-1"),
							StartedAt:   &github.Timestamp{Time: hago.Add(1 * time.Minute)},
							CompletedAt: &github.Timestamp{Time: hago.Add(2 * time.Minute)},
							Conclusion:  github.String("failure"),
						},
						{
							Name:        github.String("step-2"),
							StartedAt:   &github.Timestamp{Time: hago.Add(2 * time.Minute)},
							CompletedAt: &github.Timestamp{Time: hago.Add(4 * time.Minute)},
							Conclusion:  github.String("success"),
						},
					},
				})
				return c
			},
			validate: func(t *testing.T, res *Stats, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, 1, len(res.Jobs))
				require.Equal(t, 10*time.Minute, res.Jobs[DefaultJobName].P50)
				require.Equal(t, 10*time.Minute, res.Jobs[DefaultJobName].P95)
				require.Equal(t, 10*time.Minute, res.Jobs[DefaultJobName].P99)
				// steps
				require.Equal(t, 0, res.Steps["step-1"].Successes)
				require.Equal(t, 1, res.Steps["step-1"].Failures)
				require.Equal(t, 1*time.Minute, res.Steps["step-1"].P50)
				require.Equal(t, 1*time.Minute, res.Steps["step-1"].P95)
				require.Equal(t, 1*time.Minute, res.Steps["step-1"].P99)
				require.Equal(t, 1, res.Steps["step-2"].Successes)
				require.Equal(t, 0, res.Steps["step-2"].Failures)
				require.Equal(t, 2*time.Minute, res.Steps["step-2"].P50)
				require.Equal(t, 2*time.Minute, res.Steps["step-2"].P95)
				require.Equal(t, 2*time.Minute, res.Steps["step-2"].P99)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			res, err := AnalyzeJobsSteps(context.Background(), c, tt.config)
			tt.validate(t, res, err)
		})
	}
}
