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

func print_help() string {
	var output []string
	output = append(output, fmt.Sprintf("**%s plugins** - /bot plugins", globals.APP_NAME))
	output = append(output, "---------------------------------------")
	for _, plugin := range p.Plugins {
		var s string
		t := strings.Join(plugin.Triggers, ",  ")
		s = fmt.Sprintf("%s (%s) - %s", plugin.Name, t, plugin.Description)
		output = append(output, s)
	}
	return strings.Join(output, "\n")
}

func message_delete(s *discordgo.Session, chan_id string, m_id string) {
	c := config.Get()
	if !c.KeepTriggers {
		s.ChannelMessageDelete(chan_id, m_id)
	}
}

func handle_json_plugin(plugin_name string, stdout []byte, s *discordgo.Session) {
	var plugin_response types.Message
	var err error
	err = json.Unmarshal(stdout, &plugin_response)
	if err != nil {
		log.Warning("Plugin \"%s\" did not reply with properly formatted json: %s", plugin_name, err)
		return
	}
	err = Message(s, &plugin_response)
	if err != nil {
		log.Warning("Discord API returned an error: %s", err)
		return
	} else {
		log.Debug("JSON plugin was successful")
	}
}

func check_trigger(trigger string, message string) (status int, msg_split []string) {
	msg_split = strings.Split(message, " ")

	// filter down to just the trigger
	msg_trigger := msg_split[0]

	// if the length of the content is greater than the trigger
	// and the slice of the message to the length of the trigger is the trigger...
	if len(message) > len(trigger) && msg_trigger == trigger {
		status = globals.MATCH
		// if the length of the content is the size of the trigger, and the message is the trigger
	} else if len(msg_split) == 1 && msg_trigger == trigger {
		status = globals.HELP
	} else if len(msg_split) == 2 && msg_split[1] == "" && msg_trigger == trigger {
		// edge case where strings.Split returns 2 element slice
		// where the second element is the token used to split
		// but it's an empty string
		status = globals.HELP
	}

	// 0 is globals status uninitialized value
	if status == 0 {
		status = globals.NO_MATCH
	}

	return
}

func handleInternalCommand(status int, command []string, s *discordgo.Session, m *discordgo.Message) (handled bool) {
	if status == globals.MATCH || status == globals.HELP {
		message_delete(s, m.ChannelID, m.ID)
		handled = true
	}
	if status == globals.MATCH {
		switch {
		case command[1] == "plugins":
			s.ChannelMessageSend(m.ChannelID, print_help())
		case command[1] == "reload":
			plugins.Register()
			s.ChannelMessageSend(m.ChannelID, "Digo Reloaded")
		default:
			s.ChannelMessageSend(m.ChannelID, print_help())
		}
	} else if status == globals.HELP {
		s.ChannelMessageSend(m.ChannelID, print_help())
	}
	return
}

func handleExternalCommand(status int, command []string, plugin_dir string, plugin *types.Plugin, s *discordgo.Session, m *discordgo.Message) (handled bool) {
	//s := conn.Get()
	// always send the event to message_delete
	if status == globals.MATCH || status == globals.HELP {
		message_delete(s, m.ChannelID, m.ID)
		handled = true
	}

	if status == globals.MATCH && plugin.Type == "simple" {
		output, err := plugins.Exec(plugin_dir, plugin.Filename, command[1:])
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
		output, err := plugins.ExecJson(plugin_dir, plugin.Filename, &message)
		if err == nil {
			handle_json_plugin(plugin.Name, output, s)
		} else {
			log.Warning("Could not exec json plugin %s", plugin.Name)
		}

	} else if status == globals.HELP {
		output, err := plugins.Exec(plugin_dir, plugin.Filename, []string{})
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, string(output))
		}

	}
	return
}

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
	status, command = check_trigger(c.Trigger, m.Content)

	handled = handleInternalCommand(status, command, s, m)
	if handled {
		return
	}

	// clean up the command
	// if status != globals.NO_MATCH {
	// }

	for trigger, plugin_file := range p.AllTriggers {
		// grab the correct plugin from the Plugins struct
		plugin := p.Plugins[plugin_file]

		status, command = check_trigger(trigger, m.Content)
		handled = handleExternalCommand(status, command, p.Directory, plugin, s, m)
		if handled {
			return
		}

	}
}

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

func Init(plugins *types.Plugins, logger *logging.Logger) {
	// set the config pointer
	p = plugins
	log = logger
}
