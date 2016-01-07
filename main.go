package main

import (
	//"github.com/bwmarrin/discordgo"
	//"github.com/op/go-logging"
	//"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/digo/api"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/conn"
	"github.com/sethdmoore/digo/errhandler"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/handler"
	"github.com/sethdmoore/digo/logger"
	"github.com/sethdmoore/digo/plugins"
	"github.com/sethdmoore/digo/types"
)

func main() {
	var err error

	//var p *plugins.Plugins
	var c *types.Config

	//p = plugins.Init()
	lock := make(chan int)

	// set up the config struct
	c = config.Init()

	// set the log reference to pass around
	log := logger.Init(c)
	errhandler.Init(log)

	// set up the plugins struct
	p := plugins.Init(log)
	//log.Notice()

	// handler takes reference to config and plugins structs
	handler.Init(c, p, log)

	// login / websocket flow
	s := conn.Init(c, log)

	// determine the bot's userID
	user, err := s.User("@me")
	errhandler.Handle(err)

	c.UserID = user.ID

	// allow live plugins
	//go plugins.Poll(p)

	// listen for events on Discord
	// conn.Listen(s, c, log)

	// enable the API, if applicable
	if c.DisableApi {
		log.Notice("API explicitly disabled.")
	} else {
		go api.Listen(c.ApiInterface, s, c, log)
	}

	log.Noticef("Digo v%s Online", globals.VERSION)

	<-lock
}
