package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	//"github.com/jzelinskie/geddit"
	"os"
	"time"
	//"github.com/franela/goreq"
	"github.com/sethdmoore/digo/api"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/globals"
)

// q=Dota+2+Update+-+MAIN+CLIENT+-++author%3Asirbelvedere&amp=&restrict_sr=on&t=hour&sort=new

// need to package scope this
// as there's no obvious way to pass other params to messageCreate
var c *config.Config

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(2)
	}
}

func messageCreate(s *discordgo.Session, m discordgo.Message) {
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	if len(m.Content) > len(c.Trigger) && m.Content[:len(c.Trigger)] == c.Trigger {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		//spew.Dump(m)
		//spew.Dump(s)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command \"%s\" received from %s", m.Content[len(c.Trigger):], m.Author.Username))
	}
}

func main() {
	var err error
	lock := make(chan int)
	fmt.Println("nada")
	//Uri := "https://discordapp.com/api/auth/login"
	c = config.Init()
	spew.Dump(c) //DEBUG
	/*
		res, err := goreq.Request{
			Method:      "POST",
			Uri:         Uri,
			ContentType: "application/json",
			Body:        c,
		}.Do()
	*/

	dg := discordgo.Session{
		OnMessageCreate: messageCreate,
	}

	dg.Token, err = dg.Login(c.Email, c.Password)
	handleError(err)

	// open websocket...
	err = dg.Open()
	handleError(err)

	// shouldn't this be abstracted...?
	err = dg.Handshake()
	handleError(err)

	go dg.Listen()
	fmt.Printf("Digo version %s\n", globals.Version)

	<-lock
}
