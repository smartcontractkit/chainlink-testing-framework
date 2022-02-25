// Package testreporters holds all the tools necessary to report on tests that are run utilizing the testsetups package
package testreporters

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

// TestReporter is a general interface for all test reporters
type TestReporter interface {
	WriteReport(folderLocation string) error
	SendSlackNotification(slackClient *slack.Client) error
	SetNamespace(namespace string)
}

// Common Slack Notification Helpers

// Values for reporters to use slack to notify user of test end
var (
	slackWebhook = os.Getenv("SLACK_WEBHOOK")
	slackAPIKey  = os.Getenv("SLACK_API")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUserID  = os.Getenv("SLACK_USER_ID")
)

// UpdateSlackEnvVars updates the slack environment variables in case they are changed while remote test is running.
// Usually used for unit tests.
func UpdateSlackEnvVars() {
	slackWebhook = os.Getenv("SLACK_WEBHOOK")
	slackAPIKey = os.Getenv("SLACK_API")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUserID = os.Getenv("SLACK_USER_ID")
}

// Sends a slack webhook message
func sendSlackWebhook(webhookBlocks *slack.WebhookMessage) error {
	msgBytes, err := json.Marshal(webhookBlocks)
	if err != nil {
		log.Error().Err(err).Interface("Webhook Message", webhookBlocks).Msg("Error marshalling webhook message to JSON")
	}
	log.Info().Str("Webhook URL", slackWebhook).Str("Message Body", string(msgBytes)).Msg("Sending Slack Notification")
	return slack.PostWebhook(slackWebhook, webhookBlocks)
}

// Uploads a slack file to the designated channel using the API key
func uploadSlackFile(slackClient *slack.Client, uploadParams slack.FileUploadParameters) error {
	log.Info().
		Str("Slack API Key", slackAPIKey).
		Str("Slack Channel", slackChannel).
		Str("User Id to Notify", slackUserID).
		Str("File", uploadParams.File).
		Msg("Attempting to upload file")
	if slackAPIKey == "" {
		return errors.New("Unable to upload file without a Slack API Key")
	}
	if slackChannel == "" {
		return errors.New("Unable to upload file without a Slack Channel")
	}
	if uploadParams.Channels == nil || uploadParams.Channels[0] == "" {
		uploadParams.Channels = []string{slackChannel}
	}
	_, err := slackClient.UploadFile(uploadParams)
	return err
}
