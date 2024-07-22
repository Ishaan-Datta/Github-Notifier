package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DiscordPayload represents the JSON payload structure for sending a message to Discord
type DiscordPayload struct {
	Content string `json:"content"`
}

// sendDiscordMessage sends a message to a Discord channel using the Discord API
func sendDiscordMessage(payload DiscordPayload) {
	discord_token := config.DiscordBotToken

	if discord_token == "" {
		resultQueue <- TaskResult{application: "discord", status: false, reason: "DISCORD_OAUTH_TOKEN is not set"}
		return
	}

	channelID := config.DiscordChannelID

	if channelID == "" {
		resultQueue <- TaskResult{application: "discord", status: false, reason: "DISCORD_CHANNEL_ID is not set"}
		return
	}

	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		resultQueue <- TaskResult{application: "discord", status: false, reason: fmt.Sprintf("failed to marshal payload: %v", err)}
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		resultQueue <- TaskResult{application: "discord", status: false, reason: fmt.Sprintf("failed to create request: %v", err)}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+discord_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultQueue <- TaskResult{application: "discord", status: false, reason: fmt.Sprintf("failed to send request: %v", err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultQueue <- TaskResult{application: "discord", status: false, reason: fmt.Sprintf("failed to read response body: %v", err)}
		return
	}

	if resp.StatusCode != http.StatusOK {
		resultQueue <- TaskResult{application: "discord", status: false, reason: fmt.Sprintf("received non-OK response: %s", body)}
		return
	} else {
		resultQueue <- TaskResult{application: "discord", status: true, reason: "None"}
		return
	}
}
