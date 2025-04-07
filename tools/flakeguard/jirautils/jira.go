package jirautils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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

const pillarCustomFieldID = "customfield_11016" // Define constant for pillar field

// CreateTicketInJira creates a new Jira ticket with specified details.
// Adds support for assignee, priority, labels, and pillar name.
func CreateTicketInJira(
	ctx context.Context, // Add context for cancellation/timeout
	client *jira.Client,
	summary, description, projectKey, issueType, assigneeID, priorityName string,
	labels []string, // Add labels parameter
	pillarName string, // Add pillarName parameter
) (string, error) {
	if client == nil {
		return "", fmt.Errorf("jira client is nil")
	}

	// Basic validation
	if summary == "" || projectKey == "" || issueType == "" {
		return "", fmt.Errorf("summary, projectKey, and issueType are required")
	}

	// Prepare issue fields
	issueFields := &jira.IssueFields{
		Project:     jira.Project{Key: projectKey},
		Summary:     summary,
		Description: description,
		Type:        jira.IssueType{Name: issueType},
		Labels:      labels, // Add labels
	}

	// Add Assignee if provided
	if assigneeID != "" {
		issueFields.Assignee = &jira.User{AccountID: assigneeID} // Use AccountID typically
		// Note: Sometimes 'Name' (username) is used instead of AccountID. Verify for your Jira instance.
		// issueFields.Assignee = &jira.User{Name: assigneeID}
	}

	// Add Priority if provided
	if priorityName != "" {
		issueFields.Priority = &jira.Priority{Name: priorityName}
	}

	// Add Pillar Name (custom field) if provided
	// Custom fields often require a specific structure (e.g., map[string]string{"value": "Pillar"})
	if pillarName != "" {
		if issueFields.Unknowns == nil {
			issueFields.Unknowns = make(map[string]interface{})
		}
		// The exact structure depends on the custom field type in Jira (e.g., text, select list)
		// For a simple text field or select list (by value):
		issueFields.Unknowns[pillarCustomFieldID] = map[string]interface{}{"value": pillarName}
		// If it's a select list by ID, it would be map[string]interface{}{"id": "12345"}
		log.Debug().Str("fieldId", pillarCustomFieldID).Str("value", pillarName).Msg("Adding pillar custom field")
	}

	issue := jira.Issue{
		Fields: issueFields,
	}

	log.Debug().Interface("issuePayload", issue).Msg("Jira creation payload")

	// Use context with the API call
	createdIssue, resp, err := client.Issue.CreateWithContext(ctx, &issue)
	if err != nil {
		// Try to read response body for more details
		errMsg := ReadJiraErrorResponse(resp)
		log.Error().Err(err).Str("responseBody", errMsg).Msg("Failed to create Jira issue")
		// Return wrapped error with response details
		return "", fmt.Errorf("failed to create Jira issue: %w; response: %s", err, errMsg)
	}
	if createdIssue == nil {
		return "", fmt.Errorf("jira API returned success but issue object is nil")
	}

	log.Info().Str("key", createdIssue.Key).Str("id", createdIssue.ID).Msg("Jira issue created successfully")
	return createdIssue.Key, nil
}

// DeleteTicketInJira permanently deletes a Jira ticket. Use with caution.
func DeleteTicketInJira(client *jira.Client, issueKey string) error {
	if client == nil {
		return fmt.Errorf("jira client is nil")
	}
	if issueKey == "" {
		return fmt.Errorf("issue key cannot be empty")
	}

	log.Warn().Str("key", issueKey).Msg("Attempting to permanently delete Jira ticket")

	resp, err := client.Issue.DeleteWithContext(context.Background(), issueKey)
	if err != nil {
		errMsg := ReadJiraErrorResponse(resp)
		log.Error().Err(err).Str("key", issueKey).Str("responseBody", errMsg).Msg("Failed to delete Jira issue")
		return fmt.Errorf("failed to delete issue %s: %w; response: %s", issueKey, err, errMsg)
	}

	// Check status code explicitly? Jira library might already handle non-2xx as errors.
	// if resp.StatusCode != http.StatusNoContent { ... }

	log.Info().Str("key", issueKey).Msg("Jira ticket deleted successfully")
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

// ReadJiraErrorResponse tries to extract error messages from Jira's response body.
func ReadJiraErrorResponse(resp *jira.Response) string {
	if resp == nil || resp.Body == nil {
		return "(No response body)"
	}
	defer resp.Body.Close() // Ensure body is closed

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("(Failed to read response body: %v)", err)
	}

	bodyString := string(bodyBytes)

	// Attempt to unmarshal into Jira's standard error structure
	var jiraErr jira.Error
	if err := json.Unmarshal(bodyBytes, &jiraErr); err == nil {
		// Extract messages if available
		var messages []string
		if len(jiraErr.ErrorMessages) > 0 {
			messages = append(messages, jiraErr.ErrorMessages...)
		}
		if len(jiraErr.Errors) > 0 {
			for key, val := range jiraErr.Errors {
				messages = append(messages, fmt.Sprintf("%s: %s", key, val))
			}
		}
		if len(messages) > 0 {
			return strings.Join(messages, "; ")
		}
	}

	// If JSON parsing fails or yields no messages, return the raw body (truncated if too long)
	maxLen := 500
	if len(bodyString) > maxLen {
		return bodyString[:maxLen] + "..."
	}
	return bodyString
}

// getJiraLink returns the full Jira URL for a given ticket key if JIRA_DOMAIN is set.
func GetJiraLink(ticketKey string) string {
	domain := os.Getenv("JIRA_DOMAIN")
	if domain != "" {
		return fmt.Sprintf("https://%s/browse/%s", domain, ticketKey)
	}
	return ticketKey
}
