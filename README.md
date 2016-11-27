# Reminder Bot

This project provides some basic infrastructure for writing a FB messenger bot in go.

To run your own messenger bot, you will need to create a config.go file that contains your FB messenger bot's Verify Token and Page Access Token as constants. The config.go file should also contain your wit.ai server access token.
```
package main

const VerifyToken = "[insert]"
const PageAccessToken = "[insert]"
const DB_USER = "[insert]"
const DB_PASSWORD = "[insert]"
const DB_NAME = "[insert]"
const WitAiServerAccessToken = "[insert]"
```

To connect to the postgres DB, run ```sh db.sh``` from the ```src``` directory. To run the app, run ```go run main.go config.go``` from the ```src``` directory.
