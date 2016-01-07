package types

type PluginResponse struct {
	Plugin   string   `json:"plugin"`   // name of the plugin
	Action   string   `json:"action"`   // currently only "ChannelMessageSend"
	Channels []string `json:"channels"` // channel(s) to post the output to
	Payload  []string `json:"payload"`  // message from the plugin
	//Users    []string `json:"users"`	// users
}

type PluginRequest struct {
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
