package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SlackPayload represents the JSON payload structure for sending a message to Slack
type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// sendSlackMessage sends a message to a Slack channel using the Slack API
func sendSlackMessage(message string) {
	url := "https://slack.com/api/chat.postMessage"
	slack_token := config.SlackOAuthToken
	slack_channel := config.SlackChannelID

	if slack_token == "" {
		resultQueue <- TaskResult{application: "slack", status: false, reason: "SLACK_OAUTH_TOKEN is not set"}
		return
	}

	if slack_channel == "" {
		resultQueue <- TaskResult{application: "slack", status: false, reason: "SLACK_CHANNEL_ID is not set"}
		return
	}

	payload := SlackPayload{
		Channel: slack_channel,
		Text:    message,
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		resultQueue <- TaskResult{application: "slack", status: false, reason: fmt.Sprintf("failed to create request: %v", err)}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slack_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultQueue <- TaskResult{application: "slack", status: false, reason: fmt.Sprintf("failed to send request: %v", err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultQueue <- TaskResult{application: "slack", status: false, reason: fmt.Sprintf("failed to read response body: %v", err)}
		return
	}

	if resp.StatusCode != http.StatusOK {
		resultQueue <- TaskResult{application: "slack", status: false, reason: fmt.Sprintf("received non-OK response: %s", body)}
		return
	} else {
		resultQueue <- TaskResult{application: "slack", status: true, reason: ""}
		return
	}
}
