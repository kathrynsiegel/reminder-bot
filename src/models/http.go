package models

// Models used for http requests and responses

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
