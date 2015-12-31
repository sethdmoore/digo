package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

var session *discordgo.Session

func Listen(iface string, s *discordgo.Session) {
	r := gin.Default()
	session = s

	v1 := r.Group("/v1")
	v1.GET("/version", version_v1)
	v1.GET("/channels", channels_v1)
	v1.POST("/register/:plugin", register_plugin_v1)
	v1.POST("/message", message_v1)

	r.Run(":8080")
}
