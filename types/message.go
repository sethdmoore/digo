package types

// Message is unmarshaled from plugins and the API
type Message struct {
	Prefix   string   `json:"prefix"`
	Channels []string `json:"channels"`
	Action   string   `json:"action"`
	Payload  []string `json:"payload"`
}
