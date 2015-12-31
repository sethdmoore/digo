package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	//"github.com/sethdmoore/digo/errhandler"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/types"
	"strings"
	"time"
)

// need to package scope this
// as there's no obvious way to pass other params to MessageHandler
var c *types.Config

func MessageHandler(s *discordgo.Session, m discordgo.Message) {
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	if len(m.Content) > len(c.Trigger) && m.Content[:len(c.Trigger)] == c.Trigger {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command \"%s\" received from %s", m.Content[len(c.Trigger):], m.Author.Username))
	}
}

func Message(s *discordgo.Session, m *types.Message) (int, error) {
	message := strings.Join(m.Payload, "\n")
	var channels []string
	var dchannels []discordgo.Channel
	var err error
	spew.Dump(s.OnMessageCreate)

	if m.Prefix != "" {
		message = fmt.Sprintf("%s: %s", m.Prefix, message)
	}

	if m.Channels[0] == "*" {
		fmt.Println("wildcard!")
		dchannels, err = s.GuildChannels(c.Guild)

		if err != nil {
			return globals.ERROR, err
		}
		//errhandler.Handle(err)

		for _, chann := range dchannels {
			//spew.Dump(chann)
			channels = append(channels, chann.ID)
		}

	} else {
		channels = m.Channels
	}
	fmt.Printf("%s\n", len(channels))
	spew.Dump(channels)

	for _, channel := range channels {
		s.ChannelMessageSend(channel, message)
	}
	return globals.OK, nil
}

func Init(config *types.Config) {
	// set the config pointer
	c = config
}
