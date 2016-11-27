package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kathrynsiegel/reminder-bot/config"
	"github.com/kathrynsiegel/reminder-bot/helpers"
	"github.com/kathrynsiegel/reminder-bot/models"
)

// ReceivedMessage handles messages that are received by the messenger bot.
func (ac *AppController) ReceivedMessage(messageData models.MessagingEvent) {
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
			ac.SendMessage(messageData.Sender.ID, ac.ProcessMessage(messageData))
		}
	}
}

// SendMessage sends a message to a FB user.
func (ac *AppController) SendMessage(senderID string, text string) {
	var messageData models.MessagingEvent
	messageData.Recipient = &models.Messager{
		ID: senderID,
	}
	messageData.Message = &models.MessageRecord{
		Text: &text,
	}
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + config.PageAccessToken
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
func (ac *AppController) ProcessMessage(messageData models.MessagingEvent) string {
	if messageData.Message == nil || messageData.Message.Text == nil {
		return misunderstoodResponse
	}
	if messageData.Sender == nil {
		return ""
	}
	attributes, err := ac.GetAttributesFromMessage(*messageData.Message.Text)
	if err != nil {
		return misunderstoodResponse
	}
	response, success := ac.GenerateResponse(attributes, messageData.Sender.ID)
	ac.DB.Gorm.Create(&models.MessageLog{
		SenderID:     messageData.Sender.ID,
		Text:         *messageData.Message.Text,
		ReplySuccess: success,
	})
	return response
}

// GetAttributesFromMessage uses the wit.ai API to get the attributes present
// in a text message.
func (ac *AppController) GetAttributesFromMessage(messageText string) (*models.WitAiResponse, error) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.WitAiServerAccessToken))
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
func (ac *AppController) GenerateResponse(message *models.WitAiResponse, senderID string) (string, bool) {
	user := ac.DB.FindUserOrCreate(senderID)
	log.Println(spew.Sdump(message))
	if message.GetAttribute("greeting") != nil {
		return "Hi! I'm your friendly neighborhood reminder bot. Please ask me to remember something.", true
	}
	if message.GetAttribute("reminder") == nil {
		return "I'm sorry, you can only ask me for reminders. That doesn't seem to be a reminder. Try 'Remind me to...'", false
	}
	newReminder := &models.Reminder{
		UserID:            user.ID,
		Recurring:         false,
		RepeatInterval:    models.RepeatDaily,
		RepeatTimeOfDayMs: 0,
		RepeatEvery:       1,
	}

	// Desription of reminder
	reminderAttributes := message.Entities["reminder"]
	task := ""
	for _, reminderAttr := range reminderAttributes {
		rVal := reminderAttr.Value
		if strings.Index(strings.ToLower(rVal), "remind") == -1 && rVal != "" {
			task = fmt.Sprintf(" %s%s ", strings.ToLower(rVal[0:1]), rVal[1:len(rVal)])
			task = strings.Replace(task, " my ", " your ", -1)
			task = strings.Replace(task, " me ", " you ", -1)
		}
	}
	if task == "" {
		return misunderstoodResponse, false
	}
	newReminder.Description = task

	// Recurrence
	if message.GetAttribute("recurrence") != nil {
		newReminder.Recurring = true
	}

	// TODO(Katie) RepeatInterval
	// TODO(Katie) RepeatDay
	// TODO(Katie) RepeatDayOfMonth

	// RepeatTimeOfDayMs and Timezone
	timeOfDay := message.GetAttribute("datetime")
	if timeOfDay != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000-07:00", *timeOfDay)
		log.Println(*timeOfDay)
		if err != nil {
			return misunderstoodResponse, false
		}
		newReminder.Timezone = t.Location().String()
		newReminder.RepeatTimeOfDayMs = helpers.TimeToMs(t)
	}

	// TODO(Katie) RepeatEvery

	ac.DB.Gorm.Create(newReminder)
	return fmt.Sprintf("Ok! I will remind you to %s.", task[1:len(task)-1]), true
}
