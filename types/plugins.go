package types

type PluginMessage struct {
	User      string   `json:"user"`      // user that called the trigger
	Channel   string   `json:"channel"`   // channel trigger was called from
	Arguments []string `json:"arguments"` // arguments sent from the trigger
}

type Plugin struct {
	Name        string   `json:"name"`        // pretty print name
	Triggers    []string `json:"triggers"`    // array of commands to trigger
	Description string   `json:"description"` // optional
	Type        string   `json:"type"`        // simple or json (default:simple)
}

type Plugins struct {
	Plugins   map[string]*Plugin
	Directory string
}
