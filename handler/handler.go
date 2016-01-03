package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	//"github.com/sethdmoore/digo/errhandler"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/plugins"
	"github.com/sethdmoore/digo/types"
	"strings"
	"time"
)

// need to package scope this
// as there's no obvious way to pass other params to MessageHandler
var c *types.Config
var p *types.Plugins

func print_help() string {
	s := "This would print help"
	return s
}

func print_plugins() string {
	var output []string
	output = append(output, fmt.Sprintf("%s plugins", globals.APP_NAME))
	output = append(output, fmt.Sprintf("%s plugins"), c.Trigger)
	for _, plugin := range p.Plugins {
		var s string
		t := strings.Join(plugin.Triggers, ",  ")
		s = fmt.Sprintf("%s (%s) - %s", plugin.Name, t, plugin.Description)
		output = append(output, s)
	}
	return strings.Join(output, "\n")
}

func check_triggers(triggers []string, message string) (status int, msg_split []string) {
	msg_split = strings.Split(message, " ")

	// filter down to just the trigger
	msg_trigger := msg_split[0]

	for _, trigger := range triggers {
		// if the length of the content is greater than the trigger
		// and the slice of the message to the length of the trigger is the trigger...
		if len(message) > len(trigger) && msg_trigger == trigger {
			status = globals.MATCH
			break
			// if the length of the content is the size of the trigger, and the message is the trigger
		} else if len(msg_split) == 1 && msg_trigger == trigger {
			status = globals.HELP
			break
		} else if len(msg_split) == 2 && msg_split[1] == "" && msg_trigger == trigger {
			// edge case where strings.Split returns 2 element slice
			// where the second element is the token used to split
			// but it's an empty string
			status = globals.HELP
			break
		}

	}
	if status == 0 {
		status = globals.NO_MATCH
	}
	return
}

func MessageHandler(s *discordgo.Session, m discordgo.Message) {
	var status int
	var command []string
	fmt.Printf("%20s %20s %20s > %s\n",
		m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	// prevent the bot from triggering itself
	if m.Author.ID == c.UserID {
		return
	}

	// the /bot (or whatever) trigger alwqys has precedence
	status, command = check_triggers([]string{c.Trigger}, m.Content)
	if status == globals.MATCH {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		switch {
		case command[1] == "help":
			s.ChannelMessageSend(m.ChannelID, print_help())
		case command[1] == "plugins":
			s.ChannelMessageSend(m.ChannelID, print_plugins())
		default:
			s.ChannelMessageSend(m.ChannelID, print_plugins())
		}
		return
	} else if status == globals.HELP {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, print_plugins())
		return
	}

	// clean up the command
	if status != globals.NO_MATCH {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	for plugin_file, plugin := range p.Plugins {
		status, command = check_triggers(plugin.Triggers, m.Content)
		if status == globals.MATCH {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			output, err := plugins.Exec(p.Directory, plugin_file, command[1:])
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, string(output))
			}
			break
		} else if status == globals.HELP {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			output, err := plugins.Exec(p.Directory, plugin_file, []string{"help"})
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, string(output))
			}
		}
	}

}

func Message(s *discordgo.Session, m *types.Message) (int, error) {
	message := strings.Join(m.Payload, "\n")
	var channels []string
	var dchannels []discordgo.Channel
	var err error
	//spew.Dump(s.OnMessageCreate)

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

func Init(config *types.Config, plugins *types.Plugins) {
	// set the config pointer
	c = config
	p = plugins
}
