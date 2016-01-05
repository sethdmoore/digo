package main

import (
	"github.com/bwmarrin/discordgo"
	//"github.com/op/go-logging"
	//"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/digo/api"
	"github.com/sethdmoore/digo/config"
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

	// set up the plugins struct
	p := plugins.Init(log)
	//log.Notice()

	// handler takes reference to config and plugins structs
	handler.Init(c, p, log)

	dg := discordgo.Session{
		OnMessageCreate: handler.MessageHandler,
	}

	dg.Token, err = dg.Login(c.Email, c.Password)
	errhandler.Handle(err)

	// determine the bot's userID
	user, err := dg.User("@me")
	errhandler.Handle(err)
	c.UserID = user.ID

	// open websocket...
	err = dg.Open()
	errhandler.Handle(err)

	// shouldn't this be abstracted...?
	err = dg.Handshake()
	errhandler.Handle(err)

	// allow live plugins
	//go plugins.Poll(p)

	// listen for events on Discord
	go dg.Listen()

	// enable the API, if applicable
	if c.DisableApi {
		log.Notice("API explicitly disabled.")
	} else {
		go api.Listen(c.ApiInterface, &dg, c, log)
	}

	log.Noticef("Digo v%s Online", globals.VERSION)

	<-lock
}
