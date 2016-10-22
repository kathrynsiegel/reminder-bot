package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"models"
	"net/http"
)

func (app App) ReceivedMessage(messageData models.MessagingEvent) {
	if messageData.Message.Text != nil {
		text := *messageData.Message.Text
		switch text {
		case "image":
			break
		case "button":
			break
		case "generic":
			break
		case "receipt":
			break
		default:
			app.SendMessage(messageData.Sender.Id, app.CalculateMessageReply(messageData))
		}
	}
}

func (app App) SendMessage(senderId string, text string) {
	var messageData models.MessagingEvent
	messageData.Recipient = &models.Messager{
		Id: senderId,
	}
	messageData.Message = &models.MessageLog{
		Text: &text,
	}
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + PageAccessToken
	jsonStr, err := json.Marshal(messageData)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}

func (app App) CalculateMessageReply(messageData models.MessagingEvent) string {
	// TODO Categorize message into request for reminder, confirmation of action, and other
	// TODO Look up messageData.Sender.Id in DB
	return "Hi! I'm your friendly neighborhood reminder bot. Please ask me to remember something."
}
