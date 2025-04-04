package jirautils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
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

// CreateTicketInJira creates a Jira issue with the specified details, including priority.
func CreateTicketInJira(
	client *jira.Client,
	summary, description, projectKey, issueType, assigneeId, priorityName string,
) (string, error) {

	// --- Prepare Issue Fields ---
	fields := &jira.IssueFields{
		Project:     jira.Project{Key: projectKey},
		Summary:     summary,
		Description: description,
		Type:        jira.IssueType{Name: issueType},
		Labels:      []string{"flaky_test"}, // Default label
	}

	// Set Assignee only if assigneeId is provided
	if assigneeId != "" {
		fields.Assignee = &jira.User{AccountID: assigneeId}
	}

	// Find and Set Priority
	if priorityName != "" {
		priorities, resp, err := client.Priority.GetList()
		if err != nil {
			// Log the error but proceed without priority if fetching fails
			status := "unknown"
			if resp != nil {
				status = resp.Status
			}
			log.Warn().Err(err).Str("status", status).Msgf("Failed to fetch Jira priorities. Creating ticket without setting priority '%s'.", priorityName)
		} else {
			foundPriority := false
			for _, p := range priorities {
				if p.Name == priorityName {
					fields.Priority = &p // Set the Priority field with the found object
					foundPriority = true
					break
				}
			}
			if !foundPriority {
				// Log a warning if the specified priority name doesn't exist in Jira
				log.Warn().Msgf("Priority '%s' not found in Jira instance. Creating ticket without this priority.", priorityName)
			}
		}
	}

	// Create the issue
	issue := &jira.Issue{
		Fields: fields,
	}
	newIssue, resp, err := client.Issue.CreateWithContext(context.Background(), issue)
	if err != nil {
		// Read response body for more detailed error context
		errMsg := readResponseBody(resp)
		log.Error().Err(err).Str("response_body", errMsg).Msg("Failed to create Jira issue")
		// Return a more informative error message
		return "", fmt.Errorf("error creating Jira issue (status: %s): %w; response: %s", getResponseStatus(resp), err, errMsg)
	}

	return newIssue.Key, nil
}

// Helper function to safely read response body
func readResponseBody(resp *jira.Response) string {
	if resp == nil || resp.Body == nil {
		return "[no response body]"
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("[error reading response body: %v]", err)
	}
	// Limit body size in logs/errors if necessary
	// const maxBodyLog = 512
	// if len(bodyBytes) > maxBodyLog {
	//  return string(bodyBytes[:maxBodyLog]) + "... (truncated)"
	// }
	return string(bodyBytes)
}

// Helper function to safely get response status
func getResponseStatus(resp *jira.Response) string {
	if resp != nil {
		return resp.Status
	}
	return "unknown"
}

// DeleteTicketInJira deletes a Jira ticket with the given ticket key.
func DeleteTicketInJira(client *jira.Client, ticketKey string) error {
	resp, err := client.Issue.DeleteWithContext(context.Background(), ticketKey)
	if err != nil {
		return fmt.Errorf("error deleting Jira ticket %s: %w (resp: %v)", ticketKey, err, resp)
	}
	return nil
}

// PostCommentToTicket posts a comment to a Jira ticket identified by ticketKey.
// It returns an error if the comment cannot be added.
func PostCommentToTicket(client *jira.Client, ticketKey, comment string) error {
	cmt := jira.Comment{
		Body: comment,
	}
	_, resp, err := client.Issue.AddComment(ticketKey, &cmt)
	if err != nil {
		return fmt.Errorf("failed to add comment to ticket %s: %w (response: %+v)", ticketKey, err, resp)
	}
	return nil
}

// getJiraLink returns the full Jira URL for a given ticket key if JIRA_DOMAIN is set.
func GetJiraLink(ticketKey string) string {
	domain := os.Getenv("JIRA_DOMAIN")
	if domain != "" {
		return fmt.Sprintf("https://%s/browse/%s", domain, ticketKey)
	}
	return ticketKey
}
