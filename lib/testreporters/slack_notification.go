// Package testreporters holds all the tools necessary to report on tests that are run utilizing the testsetups package
package testreporters

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
)

// Common Slack Notification Helpers

// Values for reporters to use slack to notify user of test end
var (
	SlackAPIKey  = os.Getenv(config.EnvVarSlackKey)
	SlackChannel = os.Getenv(config.EnvVarSlackChannel)
	SlackUserID  = os.Getenv(config.EnvVarSlackUser)
)

// Uploads a slack file to the designated channel using the API key
func UploadSlackFile(slackClient *slack.Client, uploadParams slack.UploadFileV2Parameters) error {
	log.Info().
		Str("Slack API Key", SlackAPIKey).
		Str("Slack Channel", SlackChannel).
		Str("User Id to Notify", SlackUserID).
		Str("File", uploadParams.File).
		Msg("Attempting to upload file")
	if SlackAPIKey == "" {
		return fmt.Errorf("unable to upload file without a Slack API Key")
	}
	if SlackChannel == "" {
		return fmt.Errorf("unable to upload file without a Slack Channel")
	}
	if uploadParams.Channel == "" {
		uploadParams.Channel = SlackChannel
	}
	if uploadParams.File != "" {
		file, err := os.Stat(uploadParams.File)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("unable to upload file as it does not exist: %w", err)
		} else if err != nil {
			return err
		}
		// file size is now mandatory, so we need to set if it's empty
		if uploadParams.FileSize == 0 {
			uploadParams.FileSize = int(file.Size())
		}
	}
	_, err := slackClient.UploadFileV2(uploadParams)
	return err
}

// Sends a slack message, and returns an error and the message timestamp
func SendSlackMessage(slackClient *slack.Client, msgOptions ...slack.MsgOption) (string, error) {
	log.Info().
		Str("Slack API Key", SlackAPIKey).
		Str("Slack Channel", SlackChannel).
		Msg("Attempting to send message")
	if SlackAPIKey == "" {
		return "", fmt.Errorf("unable to send message without a Slack API Key")
	}
	if SlackChannel == "" {
		return "", fmt.Errorf("unable to send message without a Slack Channel")
	}
	msgOptions = append(msgOptions, slack.MsgOptionAsUser(true))
	_, timeStamp, err := slackClient.PostMessage(SlackChannel, msgOptions...)
	return timeStamp, err
}

// creates a directory if it doesn't already exist
func MkdirIfNotExists(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %s err: %w", dirName, err)
		}
	}
	return nil
}

func CommonSlackNotificationBlocks(
	headerText, namespace,
	reportCsvLocation string,
) []slack.Block {
	return SlackNotifyBlocks(headerText, namespace, []string{
		fmt.Sprintf("Summary CSV created on _remote-test-runner_ at _%s_\nNotifying <@%s>",
			reportCsvLocation, SlackUserID)})
}

// SlackNotifyBlocks creates a slack payload and writes into the specified json
func SlackNotifyBlocks(headerText, namespace string, msgtext []string) []slack.Block {
	var notificationBlocks slack.Blocks
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet, slack.NewDividerBlock())
	msgtexts := ""
	for _, text := range msgtext {
		msgtexts = fmt.Sprintf("%s%s\n", msgtexts, text)
	}
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		msgtexts, false, true), nil, nil))
	return notificationBlocks.BlockSet
}
