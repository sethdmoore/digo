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
