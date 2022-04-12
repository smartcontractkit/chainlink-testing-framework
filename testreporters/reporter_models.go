// Package testreporters holds all the tools necessary to report on tests that are run utilizing the testsetups package
package testreporters

import (
	"errors"
	"fmt"
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
	slackAPIKey  = os.Getenv("SLACK_API")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUserID  = os.Getenv("SLACK_USER_ID")
)

// UpdateSlackEnvVars updates the slack environment variables in case they are changed while remote test is running.
// Usually used for unit tests.
func UpdateSlackEnvVars() {
	slackAPIKey = os.Getenv("SLACK_API")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUserID = os.Getenv("SLACK_USER_ID")
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
	if uploadParams.File != "" {
		if _, err := os.Stat(uploadParams.File); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Unable to upload file as it does not exist: %w", err)
		} else if err != nil {
			return err
		}
	}
	_, err := slackClient.UploadFile(uploadParams)
	return err
}

// Sends a slack message, and returns an error and the message timestamp
func sendSlackMessage(slackClient *slack.Client, msgOptions ...slack.MsgOption) (string, error) {
	log.Info().
		Str("Slack API Key", slackAPIKey).
		Str("Slack Channel", slackChannel).
		Msg("Attempting to send message")
	if slackAPIKey == "" {
		return "", errors.New("Unable to send message without a Slack API Key")
	}
	if slackChannel == "" {
		return "", errors.New("Unable to send message without a Slack Channel")
	}
	msgOptions = append(msgOptions, slack.MsgOptionAsUser(true))
	_, timeStamp, err := slackClient.PostMessage(slackChannel, msgOptions...)
	return timeStamp, err
}
