package platsec_slack_integration_go

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	username   string
	token      string
	apiUrl     string
	awsAccount string
	ssmRole    string
}

type (
	MessagePayload struct {
		ChannelLookup  `json:"channelLookup"`
		MessageDetails `json:"messageDetails"`
	}
	ChannelLookup struct {
		By            string   `json:"by"`
		SlackChannels []string `json:"slackChannels"`
	}
	MessageDetails struct {
		Text        string           `json:"text"`
		Attachments []AttachmentItem `json:"attachments"`
	}
	AttachmentItem struct {
		Color string `json:"color"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}
)

const (
	SLACK_API_URL_ENV_NAME      = "SLACK_API_URL"
	SLACK_USERNAME_KEY_ENV_NAME = "SLACK_USERNAME_KEY"
	SLACK_TOKEN_KEY_ENV_NAME    = "SLACK_TOKEN_KEY"
	SSM_READ_ROLE_ENV_NAME      = "SSM_READ_ROLE"
	AWS_ACCOUNT                 = "AWS_ACCOUNT"
)

type HttpPostAPI interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// notifySlack sends message to Slack via the Platops service.
func notifySlack(config SlackNotifierConfig, message []byte, httpClient HttpPostAPI) (*http.Response, error) {
	responseBody := bytes.NewBuffer(message)
	response, err := httpClient.Post(config.apiUrl, "application/json", responseBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// createSlackMessages generates a list of slack messages depending on the number of channels supplied.
func createSlackMessages(channels []string, header string, title string, text string, colour string) ([]SlackMessage, int) {
	slackMessages := make([]SlackMessage, len(channels))

	if len(channels) < 1 {
		return []SlackMessage{}, 0
	}

	slackMessages = append(slackMessages, SlackMessage{
		channels: channels,
		header:   header,
		title:    title,
		text:     text,
		colour:   colour,
	})

	return slackMessages, len(channels)
}

// NewSlackNotifierConfig returns a config struct.
func NewSlackNotifierConfig(username string, token string, apiUrl string, awsAccount string, ssmRole string) SlackNotifierConfig {
	return SlackNotifierConfig{
		username:   username,
		token:      token,
		apiUrl:     apiUrl,
		awsAccount: awsAccount,
		ssmRole:    ssmRole,
	}
}

func SendMessageWithEnvVars(channels []string, header string, title string, text string, colour string) bool {
	keysToValidate := []string{"SLACK_API_URL", "SLACK_USERNAME_KEY", "SLACK_TOKEN_KEY", "SSM_READ_ROLE", "AWS_ACCOUNT"}
	keysPresent := validateEnvConfig(keysToValidate, 0)

	if !keysPresent {
		return keysPresent
		os.Exit(-1)
	}

	slackConfig := assignConfigItems(getEnvConfig())
	slackService := SlackService{}
	slackMessages, msgCount := createSlackMessages(channels, header, title, text, colour)

	if msgCount == 0 {
		os.Exit(-1)
	}

	for _, slackMessage := range slackMessages {
		msgPayload := generatePayload(slackMessage)
		messageData, err := marshallPayload(msgPayload)
		if err != nil {
			os.Exit(-1)
		}
		resp, err := notifySlack(slackConfig, messageData, slackService)
		if resp.StatusCode != 200 {
			os.Exit(-1)
		}
	}

	return true
}

func SendMessageWithParams() bool {
	return true
}

func buildHeaders(config SlackNotifierConfig) map[string]string {
	src := fmt.Sprintf("%s:%s", config.username, config.token)
	sEnd := b64.StdEncoding.EncodeToString([]byte(src))

	return map[string]string{"ContentType": "application/json", "Authorization": fmt.Sprintf("Basic %s",
		sEnd)}
}

// GetEnvConfig returns environment variables needed to execute the Slack service.
func getEnvConfig() map[string]string {
	return map[string]string{
		SLACK_TOKEN_KEY_ENV_NAME:    os.Getenv(SLACK_TOKEN_KEY_ENV_NAME),
		SLACK_USERNAME_KEY_ENV_NAME: os.Getenv(SLACK_USERNAME_KEY_ENV_NAME),
		SLACK_API_URL_ENV_NAME:      os.Getenv(SLACK_API_URL_ENV_NAME),
		SSM_READ_ROLE_ENV_NAME:      os.Getenv(SSM_READ_ROLE_ENV_NAME),
		AWS_ACCOUNT:                 os.Getenv(AWS_ACCOUNT),
	}
}

// ValidateEnvConfig validates keys in config settings.
func validateEnvConfig(configKeys []string, index int) bool {
	validKey := true
	if index <= len(configKeys)-1 {
		_, result := os.LookupEnv(configKeys[index])
		validKey = result
		if result {
			index++
			if index <= len(configKeys)-1 {
				validKey = validateEnvConfig(configKeys, index)
			}
		}
	}
	return validKey
}

// assignConfigItems creates a new SlackNotifierConfig struct from passed in items.
func assignConfigItems(configItems map[string]string) SlackNotifierConfig {
	if len(configItems) != 5 {
		return SlackNotifierConfig{}
	}

	return SlackNotifierConfig{
		username:   configItems["SLACK_USERNAME_KEY"],
		token:      configItems["SLACK_TOKEN_KEY"],
		apiUrl:     configItems["SLACK_API_URL"],
		ssmRole:    configItems["SSM_READ_ROLE"],
		awsAccount: configItems["AWS_ACCOUNT"],
	}
}

// generatePayload creates a service specific message.
func generatePayload(msg SlackMessage) MessagePayload {
	payload := MessagePayload{
		ChannelLookup{
			By:            "slack-channel",
			SlackChannels: msg.channels,
		},
		MessageDetails{
			Text: msg.header,
			Attachments: []AttachmentItem{
				{
					Color: msg.colour,
					Title: msg.title,
					Text:  msg.text,
				},
			},
		},
	}

	return payload
}

// marshallPayload returns a string representation of JSON data.
func marshallPayload(msg MessagePayload) ([]byte, error) {
	var jsonData []byte

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
