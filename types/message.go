package types

type Message struct {
	Prefix   string   `json:"prefix"`
	Channels []string `json:"channels"`
	Payload  []string `json:"payload"`
}
