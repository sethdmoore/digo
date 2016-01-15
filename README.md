# Digo
[![Go report](http://goreportcard.com/badge/sethdmoore/digo)](http://goreportcard.com/report/sethdmoore/digo)  
A pluggable bot for your Discord server, written in Golang.

## Notice
Digo is in very active development. At this point, things shouldn't  
drastically change, but they may.

## Features
* (Mostly) stable bot, with automatic rejoin / reconnect!
* Runs any number of plugins, written in any language
* Exposes a simple API for long-running services and daemons.
* Cross platform (Linux, Mac, in theory, Windows?)

## Configuration
Currently, Digo is configured through environment variables. Might support configuration files (TOML?) in the future.  
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
export DIGO_DISCORD_INVITE_ID=xuja8ije23
export DIGO_TRIGGER=/bot
export DIGO_LOG_LEVEL=info  # set to debug at your own risk
export DIGO_LOG_FILE=/var/log/digo.log  # set to debug at your own risk
export DIGO_LOG_STREAMS=/var/log/digo.log  # set to debug at your own risk
./digo
```

setting               |   type    | default           |   description                                          | required
----------------------|-----------|-------------------|--------------------------------------------------------|---------
DIGO_DISCORD_EMAIL    |  string   | xxxxxxx           | Bot's discord login email                              | yes
DIGO_DISCORD_PASS     |  strirg   | xxxxxxx           | Bot's discord login pass                               | yes
DIGO_DISCORD_INVITE_ID|  string   | xxxxxxx           | Bot's invite code                                      | no
DIGO_SERVER_ID        |  string   | xxxxxxx           | Discord server ID (guild)                              | yes
DIGO_DISABLE_API      |  boolean  | false             | Disable Digo API (default enabled)                     | No
DIGO_API_INTERFACE    |  string   | 127.0.0.1:8086    | Interface API listens on                               | No
DIGO_API_USERNAME     |  string   | unset             | Basic Auth username for the API                        | No
DIGO_API_PASSWORD     |  string   | unset             | Basic auth password for the API. If unset, no auth.    | No
DIGO_TRIGGER          |  string   | /bot              | Bot trigger for chat                                   | No
DIGO_LEAVE_TRIGGERS   |  boolean  | false             | Leave triggers in chat. Otherwise, triggers are deleted| No
DIGO_LOG_LEVEL        |  string   | info              | Log level to show. (debug, notice, warn, error)        | No
DIGO_LOG_FILE         |  string   | /var/log/digo.log | Path to the Digo log file. Can be local, EG: "digo.log"| No
DIGO_LOG_STREAMS      |  string   | stdout            | Bot outputs to "stdout" or "stdout,file" or "file"     | No

Digo only needs the following permissions from Discord:
* Read Messages
* Send Messages
* Manage Messages (if DIGO_LEAVE_TRIGGERS is false, to delete trigger messages)

## Bot built-in commands
This assumes you have not changed your DIGO_TRIGGER from /bot  
>/bot reload    # registers and unregisters plugins added / removed from the plugin directory  
>/bot plugins   # lists all plugins installed and their triggers  

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

If the trigger is mentioned with no arguments, Digo simply runs the plugin with no arguments
Example, in #general, Billy writes  
>/yt  

The youtuber plugin could respond with  

>Usage: /yt search keywords  

#### JSON type plugin

If a plugin registers its type as "json", Digo will treat the plugin differently in some regards. It still expects the plugin to dump its config when passed the argument "register", but this is where the similarity to simple-type plugins ends.

When digo detects a trigger for a JSON plugin, it runs the plugin with the first argument as "json" and the second argument as stringified json.

For instance, Billy writes
/roll 6d12 3d6

The diceroller plugin will receive arguments
["./diceroller.py", "json", '{"user": "Billy", "channel": "012938521123", "arguments": ["6d12", "3d6"]}']


Here is a table for the request JSON

Field       | type     | description
------------|----------|------------
user        | string   | Username of person calling the trigger
channel     | string   | Channel the trigger originated from
arguments   | []string | All of the arguments passed to the trigger (/foo bar baz => ["bar", "baz"])


The plugin should respond with a JSON blob to stdout in the same format the API route /v1/message takes

#### Example Simple plugins included
* [Youtube KeyWord Search](examples/plugins/youtuber.py)  
* [Imgur Keyword Search](examples/plugins/imgurer.py)  
* [Urban Dictionary Search](examples/plugins/urbaner.py)  
* [Random Insult Generator](examples/plugins/insulter.py) - [Original Source](https://gist.github.com/quandyfactory/258915)

#### Example JSON plugins included
* [Dice Roller](examples/plugins/diceroller.py)

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
- [x] Restart / resume disconnected sessions without restarting Digo
- [x] Hot registering / reloading when new plugins are added (without restart)
- [x] Logging
- [x] Godeps
- [x] Upgrade bwmarrin/discordgo Develop branch
- [ ] Support Windows plugins by invoking interpreter directly (os/exec)
- [ ] Allow triggering plugins from content of messages (regexp, instead of just /commands)
- [ ] Configurable plugin paths
