# Digo
A pluggable bot for your Discord server, written in Golang.

## Features
* a bot that runs easy-to-write plugins
* exposes a simple API for long-running services and daemons.
* cross platform (Linux, Mac, in theory, Windows?)

## Configuration
Currently, Digo is configured through environment variables. Might support configuration files (TOML?) if anyone cares.  
Your server ID (also known as Guild ID) can be [found here](https://support.discordapp.com/hc/en-us/articles/206346498)  
Here is an example wrapper script.


```sh
#!/bin/sh
export DIGO_DISCORD_EMAIL=foo@bar.com
export DIGO_DISCORD_PASS=*****
export DIGO_SERVER_ID=123456789
export DIGO_DISABLE_API=false
export DIGO_API_INTERFACE=127.0.0.1:8086
export DIGO_API_USERNAME=mydigo
export DIGO_API_PASSWORD=secretpass
export DIGO_TRIGGER=/bot
export DIGO_LOG_LEVEL=info  # set to debug at your own risk
./digo
```

setting             |   type    | default        |   description                                          | required
--------------------|-----------|-------------------------------------------------------------------------|---------
DIGO_DISCORD_EMAIL  |  string   | xxxxxxx        | Bot's discord login email                              | yes
DIGO_DISCORD_PASS   |  strirg   | xxxxxxx        | Bot's discord login pass                               | yes
DIGO_SERVER_ID      |  string   | xxxxxxx        | Discord server ID (guild)                              | yes
DIGO_DISABLE_API    |  boolean  | false          | Disable Digo API (default enabled)                     | No
DIGO_API_INTERFACE  |  string   | 127.0.0.1:8086 | Interface API listens on                               | No
DIGO_API_USERNAME   |  string   | <unset>        | Basic Auth username for the API                        | No
DIGO_API_PASSWORD   |  string   | <unset>        | Basic auth password for the API. If unset, no auth.    | No
DIGO_TRIGGER        |  string   | /bot           | Bot trigger for chat                                   | No
DIGO_LEAVE_TRIGGERS |  boolean  | false          | Leave triggers in chat. Otherwise, triggers are deleted| No
DIGO_LOG_LEVEL      |  string   | info           | Log level to show. (debug, notice, warn, error)        | No

Digo only needs the following permissions from Discord:
* Read Messages
* Send Messages
* Manage Messages (if DIGO_LEAVE_TRIGGERS is false, to delete trigger messages)

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
   "type": "simple",
   "name" : "youtuber"
}
```

Field       | type   | description
------------|--------|------------
triggers    | array  | Commands that will trigger the plugin from the chat channels
description | string | Plugin description shown for Digo's /bot plugins
name        | string | Required plugin name
type        | string | simple (default) or json. See description below for details

#### Simple type plugin
Once the plugin is registered, Digo will run the plugin whenever a trigger is mentioned in chat. Digo will pass every word after the trigger as an argument to the plugin.

##### Example, in # general, Billy writes  
>/yt cool cat videos  

The youtuber plugis arguments will be ["plugins/youtuber.py", "cool", "cat", "videos"]  

A simple plugins' stdout ends up in chat channel the trigger was called from after it exits.  
Simple plugin stderr ends up in Digo's stdout (sorry).

If the trigger is mentioned with no arguments, Digo assumes the user needs help, and will simply pass "help" to the arguments of the plugin. The plugin is free to ignore this in the case of plugins that have only one function.
Example, in #general, Billy writes  
>/yt  

The youtuber plugin will respond with  

>Usage: /yt search keywords  

#### Example Simple plugins included
* [Youtube KeyWord Search](examples/plugins/youtuber.py)  
* [Imgur Keyword Search](examples/plugins/imgurer.py)  
* [Urban Dictionary Search](examples/plugins/urbaner.py)  
* [Random Insult Generator](examples/plugins/insulter.py) - [Original Source](https://gist.github.com/quandyfactory/258915)

#### Example Services included
* [Twitter Stalker](examples/services)


## API
The API exposes a few routes. Basic Authentication can be enabled by exporting
DIGO_API_USERNAME and DIGO_API_PASSWORD.

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
* restart / resume disconnected sessions without restarting Digo (sorry)
* hot registering / reloading when new plugins are added (without restart)
* logging
* allow triggering plugins from content of messages (regexp, instead of just /commands)
* Godeps
