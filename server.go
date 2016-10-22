package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", WebhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type WebhookData struct {
	Object string      `json:"object"`
	Entry  []PageEntry `json:"entry"`
}

type PageEntry struct {
	Id        string           `json:"id"`
	Time      int64            `json:"time"`
	Messaging []MessagingEvent `json:"messaging"`
}

type MessagingEvent struct {
	Sender    *Messager   `json:"sender,omitempty"`
	Recipient *Messager   `json:"recipient,omitempty"`
	Timestamp *int64      `json:"timestamp,omitempty"`
	Optin     *MessageLog `json:"optin,omitempty"`
	Message   *MessageLog `json:"message,omitempty"`
	Delivery  *MessageLog `json:"delivery,omitempty"`
	Postback  *MessageLog `json:"postback,omitempty"`
}

type Messager struct {
	Id string `json:"id"`
}

type MessageLog struct {
	Mid  *string `json:"mid,omitempty"`
	Seq  *int64  `json:"seq,omitempty"`
	Text *string `json:"text,omitempty"`
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
		var data WebhookData
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

func receivedMessage(messageData MessagingEvent) {
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
			sendTextMessage(messageData.Sender.Id, CalculateMessageReply(messageData))
		}
	}
}

func sendTextMessage(senderId string, text string) {
	var messageData MessagingEvent
	messageData.Recipient = &Messager{
		Id: senderId,
	}
	messageData.Message = &MessageLog{
		Text: &text,
	}
	callSendAPI(messageData)
}

func callSendAPI(messageData MessagingEvent) {
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
