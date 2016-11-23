package main

import (
	"encoding/json"
	"fmt"
	"helpers"
	"log"
	"models"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App contains all major components of the app.
type App struct {
	DB models.Database
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	gormDB, err := gorm.Open("postgres", dbInfo)
	helpers.PanicIfError(err)
	defer gormDB.Close()
	app := &App{
		DB: models.Database{Gorm: gormDB},
	}
	app.automigrate()
	http.HandleFunc("/webhook", app.WebhookHandler)
	fmt.Printf("web server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (app *App) automigrate() {
	app.DB.Gorm.AutoMigrate(
		&models.User{},
		&models.Reminder{},
	)
}

// WebhookHandler handles the endpoint that FB messenger bots must support
// to verify its identity and to receive messages from users.
func (app *App) WebhookHandler(w http.ResponseWriter, req *http.Request) {
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
						app.ReceivedMessage(messagingEvent)
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
