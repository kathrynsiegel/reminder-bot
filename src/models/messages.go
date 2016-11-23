package models

// Models used for http requests and responses

// WebhookData structs capture PageEntry objects sent via FB messenger
// to the bot.
type WebhookData struct {
	Object string      `json:"object"`
	Entry  []PageEntry `json:"entry"`
}

// PageEntry structs wrap multiple messages sent via FB messenger to the bot.
type PageEntry struct {
	ID        string           `json:"id"`
	Time      int64            `json:"time"`
	Messaging []MessagingEvent `json:"messaging"`
}

// MessagingEvent structs are single messages sent via FB messenger to the bot.
type MessagingEvent struct {
	Sender    *Messager   `json:"sender,omitempty"`
	Recipient *Messager   `json:"recipient,omitempty"`
	Timestamp *int64      `json:"timestamp,omitempty"`
	Optin     *MessageLog `json:"optin,omitempty"`
	Message   *MessageLog `json:"message,omitempty"`
	Delivery  *MessageLog `json:"delivery,omitempty"`
	Postback  *MessageLog `json:"postback,omitempty"`
}

// Messager structs represent entities that send/receive messages.
type Messager struct {
	ID string `json:"id"`
}

// MessageLog structs contain information about messages sent to the bot.
type MessageLog struct {
	Mid  *string `json:"mid,omitempty"`
	Seq  *int64  `json:"seq,omitempty"`
	Text *string `json:"text,omitempty"`
}
