package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ktsiegel/reminder-bot/config"
	"github.com/ktsiegel/reminder-bot/models"
)

// AppController is the app's main controller, used for receiving
// and processing messages.
type AppController struct {
	DB *models.Database
}

// WebhookHandler handles the endpoint that FB messenger bots must support
// to verify its identity and to receive messages from users.
func (ac *AppController) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		query := req.URL.Query()
		if query.Get("hub.mode") == "subscribe" && query.Get("hub.verify_token") == config.VerifyToken {
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
						ac.ReceivedMessage(messagingEvent)
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
