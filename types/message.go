package types

type Message struct {
	Prefix   string   `json:"prefix"`
	Channels []string `json:"channels"`
	Action   string   `json:"action"`
	Payload  []string `json:"payload"`
}
