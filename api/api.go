package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/config"
	"time"
)

// package scoped so we don't have to pass them around
var session *discordgo.Session
var log *logging.Logger

func apiLog(c *gin.Context) {
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
	return apiLog
}

// Listen Tells Gin API to start
func Listen(iface string, s *discordgo.Session, logger *logging.Logger) {
	// set the refs to point to main
	var v1 *gin.RouterGroup
	session = s
	c := config.Get()
	log = logger

	if c.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	//r := gin.Default()
	r := gin.New()

	r.Use(loggerino())
	r.Use(gin.Recovery())

	if c.APIPassword != "" {
		log.Info("Basic Authentication enabled for API")
		v1 = r.Group("/v1", gin.BasicAuth(gin.Accounts{
			c.APIUsername: c.APIPassword,
		}))
	} else {
		log.Warning("DIGO_API_PASSWORD and DIGO_API_USERNAME are not set")
		log.Warning("The API is open to all requests")
		v1 = r.Group("/v1")
	}

	v1.GET("/version", versionV1)
	v1.GET("/channels", channelsV1)
	v1.POST("/message", messageV1)

	go r.Run(iface)
	log.Noticef("Digo API is listening on %s", c.APIInterface)
}
