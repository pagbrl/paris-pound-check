package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SlackNotifierEnvironment struct {
	SlackToken      string	`required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel 	string	`required:"true" envconfig:"SLACK_CHANNEL"`
}

type SlackNotifier struct {
	notify() bool
	slackToken
	slackChannel
}

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Token   string `json:"token"`
}


func makeSlackNotifier(e SlackNotifierEnvironment) SlackNotifier {
	err = envconfig.Process("poundcheck", &e)
	if err != nil {
		log.Fatalf("slackNotifier envconfig.Process: %w", err.error)
	}

	return SlackNotifier{slackToken: e.SlackToken, slackChannel: e.SlackChannel}
}

func getDefaultSlackMessage() SlackMessage {
	var slackMessage SlackMessage
	slackMessage.Channel = slackChannel

	return slackMessage
}

func getNotificationMessage(detailsUrl string, vehiclePlateNumber string) SlackMessage {
	messageString := fmt.Sprintf(":warning: Alert, Your vehicle *%s* has just been impounded. Visit %s to get details :warning:", vehiclePlateNumber, detailsUrl)

	slackMessage := GetDefaultSlackMessage()
	slackMessage.Text = messageString

	return slackMessage
}

func sendMessage(message SlackMessage) []byte {

	message.Token = slackToken
	jsonBody, err := json.Marshal(message)

	if err != nil {
		log.Fatal(err)
	}

	slackUrl := fmt.Sprintf("%s/%s", "https://slack.com/api/", "chat.postMessage")

	client := &http.Client{}
	req, err := http.NewRequest("POST", slackUrl, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", slackToken))
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
