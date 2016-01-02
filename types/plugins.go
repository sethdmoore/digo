package types

type Plugin struct {
	Name        string   `json:"name"`        // pretty print name
	Triggers    []string `json:"triggers"`    // array of commands to trigger
	Description string   `json:"description"` // optional
	//Callback    string   `json:"callback"`    // callback URL
}

type Plugins struct {
	Plugins   map[string]*Plugin
	Directory string
}
