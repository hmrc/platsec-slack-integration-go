package platsec_slack_integration_go

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"os"
)

type SlackMessage struct {
	channels []string
	header   string
	title    string
	text     string
	colour   string
}

type SlackNotifierConfig struct {
	username string
	token    string
	apiUrl   string
}

const (
	SLACK_API_URL_ENV_NAME      = "SLACK_API_URL"
	SLACK_USERNAME_KEY_ENV_NAME = "SLACK_USERNAME_KEY"
	SLACK_TOKEN_KEY_ENV_NAME    = "SLACK_TOKEN_KEY"
	SSM_READ_ROLE_ENV_NAME      = "SSM_READ_ROLE"
)

func NewSlackMessage(channels []string, header string, title string, text string, colour string) (SlackMessage, error) {
	if len(channels) < 1 {
		return SlackMessage{}, errors.New("no channels specified")
	}

	return SlackMessage{
		channels: channels,
		header:   header,
		title:    title,
		text:     text,
		colour:   colour,
	}, nil
}

func NewSlackNotifierConfig(username string, token string, apiUrl string) SlackNotifierConfig {
	return SlackNotifierConfig{
		username: username,
		token:    token,
		apiUrl:   apiUrl,
	}
}

func SendMessage(slackConfig SlackNotifierConfig) bool {
	return true
}

func buildHeaders(config SlackNotifierConfig) map[string]string {
	src := fmt.Sprintf("%s:%s", config.username, config.token)
	sEnd := b64.StdEncoding.EncodeToString([]byte(src))

	return map[string]string{"ContentType": "application/json", "Authorization": fmt.Sprintf("Basic %s",
		sEnd)}
}

// GetEnvConfig returns environment variables needed to execute the Slack service.
func GetEnvConfig() map[string]string {
	return map[string]string{
		SLACK_TOKEN_KEY_ENV_NAME:    os.Getenv(SLACK_TOKEN_KEY_ENV_NAME),
		SLACK_USERNAME_KEY_ENV_NAME: os.Getenv(SLACK_USERNAME_KEY_ENV_NAME),
		SLACK_API_URL_ENV_NAME:      os.Getenv(SLACK_API_URL_ENV_NAME),
		SSM_READ_ROLE_ENV_NAME:      os.Getenv(SSM_READ_ROLE_ENV_NAME),
	}
}

// ValidateEnvConfig validates keys in config settings.
func ValidateEnvConfig(configKeys []string, index int) bool {
	validKey := true
	if index <= len(configKeys)-1 {
		_, result := os.LookupEnv(configKeys[index])
		validKey = result
		if result {
			index++
			if index <= len(configKeys)-1 {
				validKey = ValidateEnvConfig(configKeys, index)
			}
		}
	}
	return validKey
}

var s = SendMessage(NewSlackNotifierConfig("mteasdal", "token", "s"))
