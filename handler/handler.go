package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/sethdmoore/digo/errhandler"
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/plugins"
	"github.com/sethdmoore/digo/types"
	"strings"
)

// need to package scope this
// as there's no obvious way to pass other params to MessageHandler
var p *types.Plugins
var log *logging.Logger

func printHelp() string {
	var output []string
	output = append(output, fmt.Sprintf("**%s plugins** - /bot plugins", globals.AppName))
	output = append(output, "---------------------------------------")
	for _, plugin := range p.Plugins {
		var s string
		t := strings.Join(plugin.Triggers, ",  ")
		s = fmt.Sprintf("%s (%s) - %s", plugin.Name, t, plugin.Description)
		output = append(output, s)
	}
	return strings.Join(output, "\n")
}

func messageDelete(s *discordgo.Session, chanID string, mID string) {
	c := config.Get()
	if !c.KeepTriggers {
		s.ChannelMessageDelete(chanID, mID)
	}
}

func handleJSONPlugin(pluginName string, stdout []byte, s *discordgo.Session) {
	var pluginResponse types.Message
	var err error
	err = json.Unmarshal(stdout, &pluginResponse)
	if err != nil {
		log.Warning("Plugin \"%s\" did not reply with properly formatted json: %s", pluginName, err)
		return
	}
	err = Message(s, &pluginResponse)
	if err != nil {
		log.Warning("Discord API returned an error: %s", err)
	} else {
		log.Debug("JSON plugin was successful")
	}
}

func checkTrigger(trigger string, message string) (status int, msgSplit []string) {
	msgSplit = strings.Split(message, " ")

	// filter down to just the trigger
	msgTrigger := msgSplit[0]

	// if the length of the content is greater than the trigger
	// and the slice of the message to the length of the trigger is the trigger...
	if len(message) > len(trigger) && msgTrigger == trigger {
		status = globals.MATCH
		// if the length of the content is the size of the trigger, and the message is the trigger
	} else if len(msgSplit) == 1 && msgTrigger == trigger {
		status = globals.HELP
	} else if len(msgSplit) == 2 && msgSplit[1] == "" && msgTrigger == trigger {
		// edge case where strings.Split returns 2 element slice
		// where the second element is the token used to split
		// but it's an empty string
		status = globals.HELP
	}

	// 0 is globals status uninitialized value
	if status == 0 {
		status = globals.NOMATCH
	}

	return
}

func handleInternalCommand(status int, command []string, s *discordgo.Session, m *discordgo.Message) (handled bool) {
	if status == globals.MATCH || status == globals.HELP {
		messageDelete(s, m.ChannelID, m.ID)
		handled = true
	}
	if status == globals.MATCH {
		switch {
		case command[1] == "plugins":
			s.ChannelMessageSend(m.ChannelID, printHelp())
		case command[1] == "reload":
			plugins.Register()
			s.ChannelMessageSend(m.ChannelID, "Digo Reloaded")
		default:
			s.ChannelMessageSend(m.ChannelID, printHelp())
		}
	} else if status == globals.HELP {
		s.ChannelMessageSend(m.ChannelID, printHelp())
	}
	return
}

func handleExternalCommand(status int, command []string, pluginDir string, plugin *types.Plugin, s *discordgo.Session, m *discordgo.Message) (handled bool) {
	//s := conn.Get()
	// always send the event to messageDelete
	if status == globals.MATCH || status == globals.HELP {
		messageDelete(s, m.ChannelID, m.ID)
		handled = true
	}

	if status == globals.MATCH && plugin.Type == "simple" {
		output, err := plugins.Exec(pluginDir, plugin.Filename, command[1:])
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, string(output))
		}

	} else if (status == globals.MATCH || status == globals.HELP) && plugin.Type == "json" {
		var message types.PluginMessage

		if status == globals.MATCH {
			message.Arguments = command[1:]
		} else if status == globals.HELP {
			message.Arguments = []string{}
		}

		message.User = m.Author.Username

		message.Channel = m.ChannelID
		output, err := plugins.ExecJSON(pluginDir, plugin.Filename, &message)
		if err == nil {
			handleJSONPlugin(plugin.Name, output, s)
		} else {
			log.Warning("Could not exec json plugin %s", plugin.Name)
		}

	} else if status == globals.HELP {
		output, err := plugins.Exec(pluginDir, plugin.Filename, []string{})
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, string(output))
		}

	}
	return
}

// MessageHandler is the callback for the discordgo session
func MessageHandler(s *discordgo.Session, m *discordgo.Message) {
	var status int
	var command []string
	var handled bool
	c := config.Get()

	log.Infof("%s %s > %s", m.ChannelID, m.Author.Username, m.Content)

	// prevent the bot from triggering itself
	if m.Author.ID == c.UserID {
		return
	}

	// the /bot (or whatever) trigger always has precedence
	status, command = checkTrigger(c.Trigger, m.Content)

	handled = handleInternalCommand(status, command, s, m)
	if handled {
		return
	}

	// clean up the command
	// if status != globals.NOMATCH {
	// }

	for trigger, pluginFile := range p.AllTriggers {
		// grab the correct plugin from the Plugins struct
		plugin := p.Plugins[pluginFile]

		status, command = checkTrigger(trigger, m.Content)
		handled = handleExternalCommand(status, command, p.Directory, plugin, s, m)
		if handled {
			return
		}

	}
}

// Message receives a Message struct and sends it to appropriate channels
func Message(s *discordgo.Session, m *types.Message) error {
	message := strings.Join(m.Payload, "\n")
	var channels []string
	var dchannels []*discordgo.Channel
	var err error
	c := config.Get()

	if m.Prefix != "" {
		message = fmt.Sprintf("%s: %s", m.Prefix, message)
	}

	if m.Channels[0] == "*" {
		dchannels, err = s.GuildChannels(c.Guild)

		if err != nil {
			return err
		}
		//errhandler.Handle(err)

		for _, chann := range dchannels {
			channels = append(channels, chann.ID)
		}

	} else {
		channels = m.Channels
	}
	log.Debugf("%s\n", len(channels))

	for _, channel := range channels {
		s.ChannelMessageSend(channel, message)
	}
	return nil
}

// Init sets up the plugins and logger struct
// need to change plugins to a .Get() method
func Init(plugins *types.Plugins, logger *logging.Logger) {
	// set the config pointer
	p = plugins
	log = logger
}
