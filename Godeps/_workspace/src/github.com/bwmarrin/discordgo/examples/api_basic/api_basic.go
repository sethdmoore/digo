// This file provides a basic "quick start" example of using the Discordgo
// package to connect to Discord using the low level API functions.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sethdmoore/digo/Godeps/_workspace/src/github.com/bwmarrin/discordgo"
)

func main() {

	var err error

	// Check for Username and Password CLI arguments.
	if len(os.Args) != 3 {
		fmt.Println("You must provide username and password as arguments. See below example.")
		fmt.Println(os.Args[0], " [username] [password]")
		return
	}

	// Create a new Discord Session interface and set a handler for the
	// OnMessageCreate event that happens for every new message on any channel
	dg := discordgo.Session{
		OnMessageCreate: messageCreate,
	}

	// Login to the Discord server and store the authentication token
	dg.Token, err = dg.Login(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		fmt.Println(err)
	}

	// Do websocket handshake.
	err = dg.Handshake()
	if err != nil {
		fmt.Println(err)
	}

	// Listen for events.
	go dg.Listen()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}

// This function will be called (due to above assignment) every time a new
// message is created on any channel that the autenticated user has access to.
func messageCreate(s *discordgo.Session, m *discordgo.Message) {

	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
