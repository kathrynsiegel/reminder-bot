package main

import "models"

func CalculateMessageReply(messageData models.MessagingEvent) string {
	// TODO Categorize message into request for reminder, confirmation of action, and other
	// TODO Look up messageData.Sender.Id in DB
	return "Hi! I'm your friendly neighborhood reminder bot. Please ask me to remember something."
}
