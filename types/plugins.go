package types

type Plugin struct {
	Name        string   `json:"name"`        // pretty print name
	Handle      []string `json:"handle"`      // array of commands to trigger
	Description string   `json:"description"` // optional
	Callback    string   `json:"callback"`    // callback URL
}

type Plugins struct {
	Plugin map[string]*Plugin
}
