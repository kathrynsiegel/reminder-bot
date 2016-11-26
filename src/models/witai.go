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

// GetAttribute determines whether a WitAiResponse has a specified
// attribute attr with a high confidence (0.75). It then returns
// a pointer to the value of that attr.
func (resp *WitAiResponse) GetAttribute(attr string) *string {
	attrs, ok := resp.Entities[attr]
	if !ok || len(attrs) == 0 {
		return nil
	}
	for _, attr := range attrs {
		if attr.Confidence >= 0.75 {
			return &attr.Value
		}
	}
	return nil
}
