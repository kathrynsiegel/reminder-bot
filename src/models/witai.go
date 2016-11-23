package models

// WitAiResponse structs are responses from wit.ai API queries.
// Contains a map of entities found in the sent string.
type WitAiResponse struct {
	MsgID    string                      `json:"msg_id"`
	Text     string                      `json:"_text"`
	Entities map[string][]WitAiAttribute `json:"entities"`
}

// WitAiAttribute structs contain information about a wit.ai attribute
// of a specific query string. These attributes are defined remotely
// using the wit.ai site.
type WitAiAttribute struct {
	Type       string  `json:"type"`
	Value      string  `json:"value"`
	Suggested  bool    `json:"suggested"`
	Confidence float64 `json:"confidence"`
}

// HasAttribute determines whether a WitAiResponse has a specified
// attribute attr with a high confidence (0.75).
func (resp *WitAiResponse) HasAttribute(attr string) bool {
	attrs, ok := resp.Entities[attr]
	if !ok || len(attrs) == 0 || attrs[0].Confidence < 0.75 {
		return false
	}
	return true
}
