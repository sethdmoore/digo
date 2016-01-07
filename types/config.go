package types

/*
	Email           digo login email
	Password        digo login password
	Trigger         /trigger for internal bot commands
	ApiInterface    interface API listens on, like "127.0.0.1:8080", ":80", and "0.0.0.0:1009"
	ApiUsername     Basic Auth username for API
	ApiPassword     Basic Auth password for API
	RemoveTriggers  Remove all plugin / bot triggers from chat
	Guild           server to connect to
	LogDir          log directory
	LogFile         log file
	LogStreams      log streams, comma delim. Valid options: "stdout,file"
	UserID          Non configurable, set in main.go by discord API
*/

type Config struct {
	Email        string `envconfig:"discord_email" json:"discord_email" toml:"discord_email"`
	Password     string `envconfig:"discord_pass" json:"discord_pass" toml:"discord_pass"`
	InviteId     string `envconfig:"discord_invite_id" json:"discord_invite_id" toml:"discord_invite_id"`
	Trigger      string `envconfig:"trigger" json:"trigger" toml:"trigger"`
	KeepTriggers bool   `envconfig:"keep_triggers" json:"keep_triggers" toml:"keep_triggers"`
	DisableApi   bool   `envconfig:"disable_api" json:"disable_api" toml:"disable_api"`
	ApiInterface string `envconfig:"api_interface" json:"api_interface" toml:"api_interface"`
	ApiUsername  string `envconfig:"api_username" json:"api_username" toml:"api_username"`
	ApiPassword  string `envconfig:"api_password" json:"api_password" toml:"api_password"`
	Guild        string `envconfig:"server_id" json:"server_id" toml:"server_id"`
	LogDir       string `envconfig:"log_directory" json:"log_directory" toml:"log_directory"`
	LogFile      string `envconfig:"log_file" json:"log_file" toml:"log_file"`
	LogStreams   string `envconfig:"log_streams" json:"log_streams" toml:"log_streams"`
	LogLevel     string `envconfig:"log_level" json:"log_level" toml:"log_level"`
	UserID       string // set in main.go by discord API
}
