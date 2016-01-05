package globals

const VERSION = "0.2.0"
const APP_NAME = "Digo"
const CONFIG_NAME = "config.toml"

var CONFIG_PATHS = [...]string{"/etc/digo", "/opt/digo"}

var PLUGIN_PATHS = [...]string{"/opt/digo/plugins", "/usr/local/digo/plugins", "plugins"}
