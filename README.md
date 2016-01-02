# Digo
A pluggable bot for your Discord server, written in Golang.

## Features
* currently does nothing out of the box. Plugins must be written or integrated.  
* exposes a simple API for long-running services and daemons.
* cross platform (should run on Windows, haven't tried it, though)

## Configuration
Currently, Digo is configured through environment variables. Might support configuration files (TOML?) if anyone cares.  
Your server ID (also known as Guild ID) can be [found here](https://support.discordapp.com/hc/en-us/articles/206346498)  
Here is an example wrapper script.

```sh
#!/bin/sh
export DIGO_USER=foo  # required
export DIGO_PASS=*****  # required
export DIGO_SERVER_ID=123456789  # required
export DIGO_INTERFACE=127.0.0.1:8081  # defaults to 127.0.0.1
export DIGO_TRIGGER=/cmd  # defaults to /bot
./digo
```

## Plugins
Plugins can be written in any language. If the shell can execute the program, Digo will be able to use execute it as well.

### Installing Plugins
Digo scans the following paths (in order) to determine the plugin directory. If anyone cares, this could be made configurable.

* /opt/digo/plugins
* /usr/local/digo/plugins
* ./plugins  # local to the binary

Simply create one of those directories, place your plugin executables there, and Digo will do the rest!

### Developing Plugins

When Digo first starts, it iterates over most files in the plugins directory. If the file begins with "_" or "." it is skipped. These prefixes are useful if your plugin requires an external configuration file.

Once it has the list of plugins, it runs the plugin with the argument "register". It expects the plugin to output JSON to stdout in the following format.

```json
{
   "triggers" : [
      "/youtube",
      "/yt"
   ],
   "description" : "Searches youtube for keywords",
   "name" : "youtuber"
}
```

Field       | type   | description
------------|--------|------------
triggers    | array  | Commands that will trigger the plugin from the chat channels
description | string | Plugin description shown for Digo's /plugins list
name        | string | Required plugin name

Once the plugin is registered, Digo will run the plugin whenever a trigger is mentioned in chat. Digo will pass every word after the trigger as an argument to the plugin.

Example, in # general, Billy writes  
/yt cool cat videos  
The youtuber plugin will have  
sys.stdin = ["plugins/youtuber.py", "cool", "cat", "videos"]  

The stdout of the plugin ends up in chat channel the trigger was called from after it exits.
Plugin stderr ends up in Digo's stdout (sorry).  

If the trigger is mentioned with no arguments, Digo assumes the user needs help, and will simply pass "help" to the arguments of the plugin. The plugin is free to ignore this in the case of plugins that have only one function.
Example, in #general, Billy writes  
/yt  
The youtuber plugin will respond with  
Usage: /yt search keywords  

Click here for a [fully working example script](examples/youtuber.py)


## API
The API exposes routes, but be careful as there is currently no security.

route         | method | description
--------------|--------|------------
/v1/version   | GET    | Returns version information
/v1/channels  | GET    | Returns JSON object containing all channels Digo is in
/v1/message   | POST   | POSTed JSON to this endpoint will send messages to the channels specified

#### /v1/message
This route expects JSON in the POST data. Required fields are bold.

field       | type   | description
------------|--------|------------
prefix      | string | optional prefix for the message
**channels**| array  | Channels to broadcast to. If set to ["*"], will broadcast to all channels
**payload** | array  | The message to send. Each element in the array will be on a new line

Example
```json
{
  "prefix": "Example",
  "channels": ["123456789", "987654321"],
  "payload": ["This message will go to 2 channels!",
              "SO EZ!"]
}

```

## Compiling
```sh
$ go get
$ go build
```

## Todo
* cleaning up debug output
* restart / resume disconnected sessions (sorry)
* logging