package types

type PluginMessage struct {
	User      string   `json:"user"`      // user that called the trigger
	Channel   string   `json:"channel"`   // channel trigger was called from
	Arguments []string `json:"arguments"` // arguments sent from the trigger
}

type Plugin struct {
	Name        string   `json:"name"`         // pretty print name
	Filename    string   `json:"__do_not_set"` // set by the filename
	Triggers    []string `json:"triggers"`     // gets copied to top level
	Tokens      []string `json:"tokens"`       // regular expression tokens to scan each msg
	Description string   `json:"description"`  // optional
	Type        string   `json:"type"`         // simple or json (default:simple)
}

type Plugins struct {
	Plugins     map[string]*Plugin // top level key is filename, Plugin type
	AllTriggers map[string]string  // top level trigger cache to reduce cyclomatic complexity
	AllTokens   map[string]string  // top level token cache to reduce cyclomatic complexity
	Directory   string             // plugin directory
}
