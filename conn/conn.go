package conn

/*
	Connection Manager for Digo. Should reconnect in the case of errors
*/

import (
	"github.com/bwmarrin/discordgo"
	//"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/handler"
	"time"
)

var log *logging.Logger
var session discordgo.Session

/* may be deprecated by session.Status in discordgo */
func pollConn(s *discordgo.Session) {
	c := config.Get()
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
			log.Warningf("Maybe I need a new invite? Using code %s", c.InviteID)
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

func acceptInvite(s *discordgo.Session) error {
	var err error
	c := config.Get()
	//time.Sleep(1 * time.Second)
	if c.InviteID != "" {
		log.Debugf("Attempting to accept invite: %s", c.InviteID)
		_, err = s.InviteAccept(c.InviteID)
	} else {
		log.Debug("No DIGO_INVITE_ID specified, no invite to accept.")
	}
	return err
}

func dgoListen(s *discordgo.Session) error {
	log.Notice("Digo listening for WS Events")
	// Listen blocks until it returns
	//go acceptInvite(s)
	err := s.Listen()
	return err
}

func doLogin(s *discordgo.Session) error {
	var err error
	var token string
	c := config.Get()
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

func loginFlow(s *discordgo.Session) {
	var err error

	for {
		// first we establish a login token
		err = doLogin(s)
		if err == nil {
			log.Info("Login successful")
		} else {
			log.Errorf("Login failed: %v", err)
			log.Error("Retrying login...")
			time.Sleep(3 * time.Second)
			// explicitly call continue so we don't try to handshake
			continue
		}
		// then we attempt a websocket connection
		// If anything fails along the way here, restart the process
		err = doWsHandshake(s)

		if err == nil {
			log.Debug("Websocket Handshake complete")
		} else {
			log.Errorf("Websocket handshake failed: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		err = acceptInvite(s)
		if err == nil {
			log.Notice("Accepted guild invite!")
		} else {
			log.Errorf("Could not accept invite: %v", err)
			log.Error("Retrying...")
			time.Sleep(3 * time.Second)
			continue
		}

		err = dgoListen(s)
		if err != nil {
			log.Errorf("Connection lost, restarting connection cycle.")
			time.Sleep(3 * time.Second)
			continue
		}

	}
}

func doWsHandshake(dg *discordgo.Session) error {
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

// Get returns a ref to a session
func Get() *discordgo.Session {
	return &session
}

// Init returns a Session reference
func Init(logger *logging.Logger) *discordgo.Session {
	//var err error
	// set the config reference
	log = logger

	log.Debug("Registering message handler")
	// type discordgo.Session (not a ref)
	session = discordgo.Session{
		OnMessageCreate: handler.MessageHandler,
	}

	go loginFlow(&session)

	time.Sleep(1 * time.Second)
	for {
		if session.Token != "" {
			break
		}
	}
	return &session
}
