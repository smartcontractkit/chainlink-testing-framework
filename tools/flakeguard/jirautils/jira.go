package jirautils

import (
	"context"
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
)

// GetJiraClient creates and returns a Jira client using env vars.
func GetJiraClient() (*jira.Client, error) {
	domain := os.Getenv("JIRA_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("JIRA_DOMAIN env var is not set")
	}
	email := os.Getenv("JIRA_EMAIL")
	if email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL env var is not set")
	}
	apiKey := os.Getenv("JIRA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("JIRA_API_KEY env var is not set")
	}

	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiKey,
	}
	return jira.NewClient(tp.Client(), fmt.Sprintf("https://%s", domain))
}

// CreateTicketInJira creates a new Jira ticket and returns its issue key.
func CreateTicketInJira(
	client *jira.Client,
	summary, description, projectKey, issueType string,
) (string, error) {
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project:     jira.Project{Key: projectKey},
			Summary:     summary,
			Description: description,
			Type:        jira.IssueType{Name: issueType},
			// Labels:      []string{"flaky_test"}, TODO: enable
		},
	}
	newIssue, resp, err := client.Issue.CreateWithContext(context.Background(), issue)
	if err != nil {
		return "", fmt.Errorf("error creating Jira issue: %w (resp: %v)", err, resp)
	}
	return newIssue.Key, nil
}

// DeleteTicketInJira deletes a Jira ticket with the given ticket key.
func DeleteTicketInJira(client *jira.Client, ticketKey string) error {
	resp, err := client.Issue.DeleteWithContext(context.Background(), ticketKey)
	if err != nil {
		return fmt.Errorf("error deleting Jira ticket %s: %w (resp: %v)", ticketKey, err, resp)
	}
	return nil
}
