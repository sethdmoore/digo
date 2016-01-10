package globals

const VERSION = "0.4.1"
const APP_NAME = "Digo"
const CONFIG_NAME = "config.toml"

var CONFIG_PATHS = [...]string{"/etc/digo", "/opt/digo"}

var PLUGIN_PATHS = [...]string{"/opt/digo/plugins", "/usr/local/digo/plugins", "plugins"}
