package platsec_slack_integration_go

import (
	"errors"
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

func NewSlackNotifierConfig(username string, token string, apiUrl string) (SlackNotifierConfig) {
	return SlackNotifierConfig{
		username: username,
		token:    token,
		apiUrl:   apiUrl,
	}
}

func SendMessage(slackConfig SlackNotifierConfig)bool{
  return true
}

var s = SendMessage(NewSlackNotifierConfig("mteasdal","token","s"))