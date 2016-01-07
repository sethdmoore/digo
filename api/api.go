package api

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/types"
	"github.com/sethdmoore/discordgo"
	"time"
)

// package scoped so we don't have to pass them around
var session *discordgo.Session
var config *types.Config
var log *logging.Logger

func api_log(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	// before request
	c.Next()
	// after request

	end := time.Now()
	latency := end.Sub(start)

	status := c.Writer.Status()
	client := c.ClientIP()
	method := c.Request.Method

	log.Infof("%v %v %v %v %v",
		latency, client, method, path, status)

}

func loggerino() gin.HandlerFunc {
	return api_log
}

func Listen(iface string, s *discordgo.Session, c *types.Config, logger *logging.Logger) {
	// set the refs to point to main
	var v1 *gin.RouterGroup
	session = s
	config = c
	log = logger

	if c.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	//r := gin.Default()
	r := gin.New()

	r.Use(loggerino())
	r.Use(gin.Recovery())

	if c.ApiPassword != "" {
		log.Info("Basic Authentication enabled for API")
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
	//v1.POST("/register/:plugin", register_plugin_v1)
	v1.POST("/message", message_v1)

	go r.Run(iface)
	log.Noticef("Digo API is listening on %s", c.ApiInterface)
}
