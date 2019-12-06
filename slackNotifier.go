package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type SlackNotifierEnvironment struct {
	SlackToken      string	`required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel 	string	`required:"true" envconfig:"SLACK_CHANNEL"`
}

type SlackNotifier struct {
	SlackToken string
	SlackChannel string
}

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Token   string `json:"token"`
}

func makeSlackNotifier() SlackNotifier {
	var env SlackNotifierEnvironment
	err := envconfig.Process("poundcheck", &env)
	if err != nil {
		log.Fatalf("slackNotifier envconfig.Process: %w", err)
	}

	return SlackNotifier{SlackToken: env.SlackToken, SlackChannel: env.SlackChannel}
}

func (notifier SlackNotifier) Notify(detailsUrl string, vehiclePlateNumber string) bool {
	var slackMessage = notifier.GetNotificationMessage(detailsUrl, vehiclePlateNumber)
	notifier.SendMessage(slackMessage)

	return true
}

func (notifier SlackNotifier) GetNotificationMessage(detailsUrl string, vehiclePlateNumber string) SlackMessage {
	var slackMessage SlackMessage
	messageString := fmt.Sprintf(":warning: Alert, Your vehicle *%s* has just been impounded. Visit %s to get details :warning:", vehiclePlateNumber, detailsUrl)

	slackMessage.Text = messageString
	slackMessage.Channel = notifier.SlackChannel

	return slackMessage
}

func (notifier SlackNotifier) SendMessage(message SlackMessage) []byte {
	jsonBody, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	slackUrl := fmt.Sprintf("%s/%s", "https://slack.com/api/", "chat.postMessage")
	client := &http.Client{}
	req, err := http.NewRequest("POST", slackUrl, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", notifier.SlackToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
