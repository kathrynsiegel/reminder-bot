package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// ReceivedMessage handles messages that are received by the messenger bot.
func (app *App) ReceivedMessage(messageData models.MessagingEvent) {
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
			// Only respond to a text message.
			app.SendMessage(messageData.Sender.ID, app.ProcessMessage(messageData))
		}
	}
}

// SendMessage sends a message to a FB user.
func (app *App) SendMessage(senderID string, text string) {
	var messageData models.MessagingEvent
	messageData.Recipient = &models.Messager{
		ID: senderID,
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
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

const misunderstoodResponse = "I'm sorry, I don't understand that message."

// ProcessMessage fetches a response to a message sent to the bot.
func (app *App) ProcessMessage(messageData models.MessagingEvent) string {
	if messageData.Message == nil || messageData.Message.Text == nil {
		return misunderstoodResponse
	}
	attributes, err := app.GetAttributesFromMessage(*messageData.Message.Text)
	if err != nil {
		return misunderstoodResponse
	}
	response, err := app.GenerateResponse(attributes, messageData.Sender.ID)
	if err != nil {
		return misunderstoodResponse
	}
	return response
}

// GetAttributesFromMessage uses the wit.ai API to get the attributes present
// in a text message.
func (app *App) GetAttributesFromMessage(messageText string) (*models.WitAiResponse, error) {
	reqURL, err := url.Parse("https://api.wit.ai/message?v=20161120")
	if err != nil {
		return nil, err
	}
	query := reqURL.Query()
	query.Set("q", messageText)
	reqURL.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", WitAiServerAccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response models.WitAiResponse
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GenerateResponse uses wit.ai attributes to generate a response to the message
// sent by the user.
func (app *App) GenerateResponse(message *models.WitAiResponse, senderID string) (string, error) {
	user := app.DB.FindUserOrCreate(senderID)
	log.Println(spew.Sdump(user))
	log.Println(spew.Sdump(message))
	if message.HasAttribute("greeting") {
		return "Hi! I'm your friendly neighborhood reminder bot. Please ask me to remember something.", nil
	}
	if !message.HasAttribute("reminder") {
		return "I'm sorry, you can only ask me for reminders. That doesn't seem to be a reminder. Try 'Remind me to...'", nil
	}
	reminderAttributes := message.Entities["reminder"]
	task := "do that"
	for _, reminderAttr := range reminderAttributes {
		rVal := reminderAttr.Value
		if strings.Index(strings.ToLower(rVal), "remind") == -1 && rVal != "" {
			task = fmt.Sprintf(" %s%s ", strings.ToLower(rVal[0:1]), rVal[1:len(rVal)])
			task = strings.Replace(task, " my ", " your ", -1)
			task = strings.Replace(task, " me ", " you ", -1)
		}
	}
	return fmt.Sprintf("Ok! I will remind you to %s.", task[1:len(task)-1]), nil
}
