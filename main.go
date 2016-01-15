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
)

func main() {
	var err error

	//p = plugins.Init()
	lock := make(chan int)

	// set up the config struct
	config.Init()

	// set the log reference to pass around
	log := logger.Init()
	errhandler.Init(log)

	// set up the plugins struct
	p := plugins.Init(log)
	//log.Notice()

	// handler takes reference to config and plugins structs
	handler.Init(p, log)

	// login / websocket flow
	s := conn.Init(log)

	// determine the bot's userID
	user, err := s.User("@me")
	errhandler.Handle(err)

	c := config.Get()

	c.UserID = user.ID

	// listen for events on Discord
	// conn.Listen(s, c, log)

	// enable the API, if applicable
	if c.DisableAPI {
		log.Notice("API explicitly disabled.")
	} else {
		go api.Listen(c.APIInterface, s, log)
	}

	log.Noticef("Digo v%s Online", globals.Version)

	<-lock
}
