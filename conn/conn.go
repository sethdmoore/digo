package conn

/*
	Connection Manager for Digo. Should reconnect in the case of errors
*/

import (
	"github.com/bwmarrin/discordgo"
	//"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/handler"
	"github.com/sethdmoore/digo/types"
	"time"
)

var log *logging.Logger
var c *types.Config
var session discordgo.Session

func poll_conn(s *discordgo.Session) {
	for {
		time.Sleep(10 * time.Second)
		found := false
		guilds, err := s.UserGuilds()
		for _, g := range guilds {
			if g.ID == c.Guild {
				found = true
				break
			}
		}

		if !found {
			log.Warningf("Could not find membership matching guild ID. %s", c.Guild)
			log.Warningf("Maybe I need a new invite? Using code %s", c.InviteId)
			if s.Token != "" {
				s.Close()
			}
			return
		}

		if err != nil {
			log.Warningf("Could not fetch guild info %v", err)
			if s.Token != "" {
				s.Close()
			}
			return
		}
	}
}

func accept_invite(s *discordgo.Session) error {
	var err error
	//time.Sleep(1 * time.Second)
	if c.InviteId != "" {
		log.Debugf("Attempting to accept invite: %s", c.InviteId)
		_, err = s.InviteAccept(c.InviteId)
	} else {
		log.Debug("No DIGO_INVITE_ID specified, no invite to accept.")
	}
	return err
}

func dgo_listen(s *discordgo.Session) error {
	log.Notice("Digo listening for WS Events")
	// Listen blocks until it returns
	//go accept_invite(s)
	err := s.Listen()
	return err
}

func DoLogin(s *discordgo.Session) error {
	var err error
	var token string
	log.Debug("Logging in")
	token, err = s.Login(c.Email, c.Password)
	if err == nil {
		if token != "" {
			s.Token = token
		}
	} else {
		log.Errorf("Can't log in: %s", err)
		log.Error("Maybe your credentials are invalid?")
	}

	// since we're dealing with a ref, only return the error
	return err
}

func LoginFlow(s *discordgo.Session) {
	var err error

	for {
		// first we establish a login token
		err = DoLogin(s)
		if err == nil {
			log.Info("Login successful")
		} else {
			log.Error("Login failed: %s", err)
			log.Error("Retrying login...")
			time.Sleep(3 * time.Second)
			// explicitly call continue so we don't try to handshake
			continue
		}
		// then we attempt a websocket connection
		// If anything fails along the way here, restart the process
		err = DoWsHandshake(s)

		if err == nil {
			log.Debug("Websocket Handshake complete")
		} else {
			log.Error("Websocket handshake failed: %s", err)
			time.Sleep(3 * time.Second)
			continue
		}

		err = accept_invite(s)
		if err == nil {
			log.Notice("Accepted guild invite!")
		} else {
			log.Errorf("Could not accept invite: %s", err)
			log.Error("Retrying...")
			time.Sleep(3 * time.Second)
			continue
		}

		err = dgo_listen(s)
		if err != nil {
			log.Errorf("Connection lost, restarting connection cycle.")
			time.Sleep(3 * time.Second)
			continue
		}

	}
}

func DoWsHandshake(dg *discordgo.Session) error {
	var err error

	// open websocket...
	err = dg.Open()
	if err != nil {
		log.Errorf("Problem opening WS: %s", err)
		return err
	}

	err = dg.Handshake()
	if err != nil {
		log.Errorf("Problem handshaking WS: %s", err)
	}
	return err
}

// init could probably handle this. Don't know if there's a use-case for this
func FetchSession() *discordgo.Session {
	return &session
}

func Init(conf *types.Config, logger *logging.Logger) *discordgo.Session {
	//var err error
	// set the config reference
	c = conf
	log = logger

	log.Debug("Registering message handler")
	// type discordgo.Session (not a ref)
	session = discordgo.Session{
		OnMessageCreate: handler.MessageHandler,
	}

	go LoginFlow(&session)

	time.Sleep(1 * time.Second)
	for {
		if session.Token != "" {
			break
		}
	}
	return &session
}
