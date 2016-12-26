package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ktsiegel/reminder-bot/config"
	"github.com/ktsiegel/reminder-bot/controllers"
	"github.com/ktsiegel/reminder-bot/helpers"
	"github.com/ktsiegel/reminder-bot/models"
)

// App contains all major components of the app.
type App struct {
	AppController *controllers.AppController
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DB_USER, config.DB_PASSWORD, config.DB_NAME)
	gormDB, err := gorm.Open("postgres", dbInfo)
	helpers.PanicIfError(err)
	defer gormDB.Close()
	app := &App{
		AppController: &controllers.AppController{DB: &models.Database{Gorm: gormDB}},
	}
	app.automigrate()
	http.HandleFunc("/webhook", app.AppController.WebhookHandler)
	fmt.Printf("web server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (app *App) automigrate() {
	app.AppController.DB.Gorm.AutoMigrate(
		&models.User{},
		&models.Reminder{},
		&models.MessageLog{},
	)
}
