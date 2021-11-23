package platsec_slack_integration_go

import (
	b64 "encoding/base64"
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
	username   string
	token      string
	apiUrl     string
	awsAccount string
	ssmRole    string
}

type (
	messagePayload struct {
		channelLookup
		messageDetails
	}
	channelLookup struct {
		by            string
		slackChannels []string
	}
	messageDetails struct {
		text        string
		attachments attachments
	}
	attachments struct {
		attachment []attachmentItem
	}
	attachmentItem struct {
		color string
		title string
		text  string
	}
)

const (
	SLACK_API_URL_ENV_NAME      = "SLACK_API_URL"
	SLACK_USERNAME_KEY_ENV_NAME = "SLACK_USERNAME_KEY"
	SLACK_TOKEN_KEY_ENV_NAME    = "SLACK_TOKEN_KEY"
	SSM_READ_ROLE_ENV_NAME      = "SSM_READ_ROLE"
	AWS_ACCOUNT                 = "AWS_ACCOUNT"
)

/*
type HttpPostAPI interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// notifySlack sends message to Slack via the Platops service.
func notifySlack(config SlackNotifierConfig, messages []SlackMessage, httpClient HttpPostAPI) (*http.Response, error) {
	for _, msg := range messages {
		httpClient.Post(config.apiUrl,)
	}
}
*/

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

func SendMessageWithEnvVars() bool {
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
func generatePayload(msg SlackMessage, config SlackNotifierConfig) messagePayload {
	payload := messagePayload{
		channelLookup{
			by:            "slack-channel",
			slackChannels: msg.channels,
		},
		messageDetails{
			text: msg.header,
			attachments: attachments{
				attachment: []attachmentItem{
					{
						color: msg.colour,
						title: msg.title,
						text:  msg.text,
					},
				},
			},
		},
	}

	return payload
}
