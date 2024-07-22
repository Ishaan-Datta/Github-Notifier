package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DiscordBotToken  string `json:"discord_bot_token"`
	DiscordChannelID string `json:"discord_channel_id"`
	GithubOAuthToken string `json:"github_oauth_token"`
	SlackChannelID   string `json:"slack_channel_id"`
	SlackOAuthToken  string `json:"slack_oauth_token"`
}

func readConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("could not read config file: %v", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("could not unmarshal config JSON: %v", err)
	}

	return config, nil
}
