package main

import (
	"context"
	"errors"

	"github.com/google/go-github/v50/github"
)

type FakeGHA struct {
	workflowRuns []*github.WorkflowRun
	workflowJobs []*github.WorkflowJob
}

func NewFakeGHA() *FakeGHA {
	return &FakeGHA{
		workflowRuns: make([]*github.WorkflowRun, 0),
		workflowJobs: make([]*github.WorkflowJob, 0),
	}
}

func (c *FakeGHA) WorkflowRun(run *github.WorkflowRun) {
	c.workflowRuns = append(c.workflowRuns, run)
}

func (c *FakeGHA) Job(jobs ...*github.WorkflowJob) {
	c.workflowJobs = append(c.workflowJobs, jobs...)
}

func (c *FakeGHA) ListRepositoryWorkflowRuns(_ context.Context, _, _ string, _ *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	total := len(c.workflowRuns)
	workflowRuns := &github.WorkflowRuns{
		TotalCount:   &total,
		WorkflowRuns: c.workflowRuns,
	}
	return workflowRuns, &github.Response{}, nil
}

func (c *FakeGHA) ListWorkflowJobs(_ context.Context, _, _ string, runID int64, _ *github.ListWorkflowJobsOptions) (*github.Jobs, *github.Response, error) {
	jobs := &github.Jobs{}
	for _, job := range c.workflowJobs {
		if job.RunID != nil && *job.RunID == runID {
			jobs.Jobs = append(jobs.Jobs, job)
		}
	}
	if len(jobs.Jobs) == 0 {
		return nil, nil, errors.New("no jobs found for the given run ID")
	}
	return jobs, &github.Response{}, nil
}
