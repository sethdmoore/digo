package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/sethdmoore/digo/types"
)

// package scoped so we don't have to pass them around
var session *discordgo.Session
var config *types.Config

func Listen(iface string, s *discordgo.Session, c *types.Config) {
	r := gin.Default()

	// set the refs to point to main
	session = s
	config = c

	v1 := r.Group("/v1")
	v1.GET("/version", version_v1)
	v1.GET("/channels", channels_v1)
	v1.POST("/register/:plugin", register_plugin_v1)
	v1.POST("/message", message_v1)

	r.Run(iface)
}
