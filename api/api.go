package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/types"
)

// package scoped so we don't have to pass them around
var session *discordgo.Session
var config *types.Config
var log *logging.Logger

func Listen(iface string, s *discordgo.Session, c *types.Config, logger *logging.Logger) {
	// set the refs to point to main
	var v1 *gin.RouterGroup
	session = s
	config = c
	log = logger

	if c.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if c.ApiPassword != "" {
		log.Notice("Basic Authentication enabled for API")
		v1 = r.Group("/v1", gin.BasicAuth(gin.Accounts{
			c.ApiUsername: c.ApiPassword,
		}))
	} else {
		log.Warning("DIGO_API_PASSWORD and DIGO_API_USERNAME are not set")
		log.Warning("The API is open to all requests")
		v1 = r.Group("/v1")
	}

	v1.GET("/version", version_v1)
	v1.GET("/channels", channels_v1)
	v1.POST("/register/:plugin", register_plugin_v1)
	v1.POST("/message", message_v1)

	go r.Run(iface)
	log.Noticef("Digo API is listening on %s", c.ApiInterface)
}
