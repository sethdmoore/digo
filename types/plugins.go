package types

// PluginMessage is the source object that gets sent to JSON-type plugins
type PluginMessage struct {
	User      string   `json:"user"`      // user that called the trigger
	Channel   string   `json:"channel"`   // channel trigger was called from
	Arguments []string `json:"arguments"` // arguments sent from the trigger
}

// Plugin stores all the information for a plugin file
type Plugin struct {
	Name        string   `json:"name"`         // pretty print name
	Filename    string   `json:"__do_not_set"` // set by the filename
	Triggers    []string `json:"triggers"`     // gets copied to top level
	Tokens      []string `json:"tokens"`       // regular expression tokens to scan each msg
	Description string   `json:"description"`  // optional
	Type        string   `json:"type"`         // simple or json (default:simple)
}

//Plugins contains a map of all Plugin types, plus a top level key containing all triggers
type Plugins struct {
	Plugins     map[string]*Plugin // top level key is filename, Plugin type
	AllTriggers map[string]string  // top level trigger cache to reduce cyclomatic complexity
	AllTokens   map[string]string  // top level token cache to reduce cyclomatic complexity
	Directory   string             // plugin directory
}
