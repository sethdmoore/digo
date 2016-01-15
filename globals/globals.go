package globals

// Version of the app
const Version = "0.5.0"

// AppName of bot
const AppName = "Digo"

// ConfigName is unused
const ConfigName = "config.toml"

// ConfigPaths are not used
var ConfigPaths = [...]string{"/etc/digo", "/opt/digo"}

// PluginPaths to search for plugin dir
var PluginPaths = [...]string{"/opt/digo/plugins", "/usr/local/digo/plugins", "plugins"}
