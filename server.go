package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"models"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", WebhookHandler)
	fmt.Printf("web server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func WebhookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		query := req.URL.Query()
		if query.Get("hub.mode") == "subscribe" && query.Get("hub.verify_token") == VerifyToken {
			fmt.Printf("Validating webhook\n")
			fmt.Fprintf(w, query.Get("hub.challenge"))
		} else {
			fmt.Printf("Failed validation. Make sure the validation tokens match.")
			http.Error(w, "Failed validation", http.StatusUnauthorized)
		}
		break
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var data models.WebhookData
		if err := decoder.Decode(&data); err != nil {
			panic(err)
		}
		defer req.Body.Close()
		if data.Object == "page" {
			for _, pageEntry := range data.Entry {
				for _, messagingEvent := range pageEntry.Messaging {
					if messagingEvent.Optin != nil {
						fmt.Printf("Received optin\n")
					} else if messagingEvent.Message != nil {
						fmt.Printf("Received message %s\n", *messagingEvent.Message)
						receivedMessage(messagingEvent)
					} else if messagingEvent.Delivery != nil {
						fmt.Printf("Received delivery %s\n", *messagingEvent.Delivery)
					} else if messagingEvent.Postback != nil {
						fmt.Printf("Received postback %s\n", *messagingEvent.Postback)
					}
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func receivedMessage(messageData models.MessagingEvent) {
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
			sendMessage(messageData.Sender.Id, CalculateMessageReply(messageData))
		}
	}
}

func sendMessage(senderId string, text string) {
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
