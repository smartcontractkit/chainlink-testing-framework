package jirautils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
)

const PillarCustomFieldID = "customfield_11016" // Define constant for pillar field

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
		issueFields.Unknowns[PillarCustomFieldID] = map[string]interface{}{"value": pillarName}
		// If it's a select list by ID, it would be map[string]interface{}{"id": "12345"}
	}

	issue := jira.Issue{
		Fields: issueFields,
	}

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

	resp, err := client.Issue.DeleteWithContext(context.Background(), issueKey)
	if err != nil {
		errMsg := ReadJiraErrorResponse(resp)
		log.Error().Err(err).Str("key", issueKey).Str("responseBody", errMsg).Msg("Failed to delete Jira issue")
		return fmt.Errorf("failed to delete issue %s: %w; response: %s", issueKey, err, errMsg)
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

func ExtractPillarValue(issue jira.Issue) string {
	if issue.Fields == nil {
		return ""
	}
	if pillarFieldRaw, ok := issue.Fields.Unknowns[PillarCustomFieldID]; ok && pillarFieldRaw != nil {
		// Handle different possible structures for custom fields (text, select list value)
		if pillarFieldMap, ok := pillarFieldRaw.(map[string]interface{}); ok {
			if value, ok := pillarFieldMap["value"].(string); ok {
				return value // Common for select lists
			}
			if value, ok := pillarFieldMap["name"].(string); ok {
				return value // Sometimes 'name' is used
			}
		}
		// Handle simple text field case
		if value, ok := pillarFieldRaw.(string); ok {
			return value
		}
	}
	return "" // Not found or unexpected format
}

// UpdatePillarName updates the pillar name custom field for a given issue key.
func UpdatePillarName(client *jira.Client, issueKey, targetPillar string) error {
	if client == nil || issueKey == "" || targetPillar == "" {
		return fmt.Errorf("client, issueKey, and targetPillar must be provided")
	}
	// Construct the payload carefully based on field type
	// Assuming it's a select list identified by 'value'
	updatePayload := map[string]interface{}{
		"fields": map[string]interface{}{
			PillarCustomFieldID: map[string]interface{}{
				"value": targetPillar,
			},
		},
	}

	req, err := client.NewRequest("PUT", fmt.Sprintf("rest/api/2/issue/%s", issueKey), updatePayload)
	if err != nil {
		return fmt.Errorf("failed to create Jira update request: %w", err)
	}

	resp, err := client.Do(req, nil) // No need to decode response body for PUT usually
	if err != nil {
		errMsg := ReadJiraErrorResponse(resp) // Use helper if available
		return fmt.Errorf("failed to update pillar name for %s: %w; response: %s", issueKey, err, errMsg)
	}

	return nil // Success
}
