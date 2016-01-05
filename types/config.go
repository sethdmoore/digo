package types

/*
	Email           digo login email
	Password        digo login password
	Trigger         /trigger for internal bot commands
	Interface       API interface, options like "127.0.0.1:8080", ":80", and "0.0.0.0:1009"
	RemoveTriggers  Remove all plugin / bot triggers from chat
	Guild           server to connect to
	LogDir          log directory
	LogFile         log file
	LogStreams      log streams, comma delim. Valid options: "stdout,file"
	UserID          Non configurable, set in main.go by discord API
*/

type Config struct {
	Email          string `envconfig:"user" json:"email" toml:"email"`
	Password       string `envconfig:"pass" json:"password" toml:"password"`
	Trigger        string `envconfig:"trigger" json:"trigger" toml:"trigger"`
	RemoveTriggers bool   `envconfig:"remove_triggers" json:"remove_triggers" toml:"remove_triggers"`
	DisableApi     bool   `envconfig:"disable_api" json:"disable_api" toml:"disable_api"`
	ApiInterface   string `envconfig:"interface" json:"interface" toml:"interface"`
	Guild          string `envconfig:"server_id" json:"server_id" toml:"server_id"`
	LogDir         string `envconfig:"log_directory" json:"log_directory" toml:"log_directory"`
	LogFile        string `envconfig:"log_file" json:"log_file" toml:"log_file"`
	LogStreams     string `envconfig:"log_streams" json:"log_streams" toml:"log_streams"`
	LogLevel       string `envconfig:"log_level" json:"log_level" toml:"log_level"`
	UserID         string // set in main.go by discord API
}
