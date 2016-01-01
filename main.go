package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	//"github.com/jzelinskie/geddit"
	//"github.com/franela/goreq"
	"github.com/sethdmoore/digo/api"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/errhandler"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/types"
	//"github.com/sethdmoore/digo/plugins"
	"github.com/sethdmoore/digo/handler"
)

// q=Dota+2+Update+-+MAIN+CLIENT+-++author%3Asirbelvedere&amp=&restrict_sr=on&t=hour&sort=new

func main() {
	var err error

	//var p *plugins.Plugins
	var c *types.Config

	//p = plugins.Init()
	lock := make(chan int)

	c = config.Init()
	spew.Dump(c) //DEBUG

	// handler takes a config struct
	handler.Init(c)

	dg := discordgo.Session{
		OnMessageCreate: handler.MessageHandler,
	}

	dg.Token, err = dg.Login(c.Email, c.Password)
	errhandler.Handle(err)

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

	spew.Dump(&dg.OnMessageCreate)
	spew.Dump(&dg.Token)

	// enable the API
	go api.Listen(c.Interface, &dg, c)

	fmt.Printf("Digo version %s\n", globals.VERSION)

	<-lock
}
