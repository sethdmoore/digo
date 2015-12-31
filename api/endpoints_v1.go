package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/handler"
	"github.com/sethdmoore/digo/types"
)

func register_plugin_v1(c *gin.Context) {
	var p types.Plugin
	name := c.Param("plugin")
	fmt.Println(name)

	err := c.BindJSON(&p)
	if err == nil {
		c.JSON(200, gin.H{
			"info": fmt.Sprintf("%s is registered and enabled!", p.Name),
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Cannot parse JSON",
			"error":   fmt.Sprintf("%s", err),
		})
		fmt.Printf("%s\n", err)
	}
}

func message_v1(c *gin.Context) {
	var m *types.Message
	err := c.BindJSON(&m)
	if err == nil {
		if len(m.Channels) > 0 {
			status, err := handler.Message(session, m)
			if status == globals.OK {
				c.JSON(200, gin.H{
					"info": fmt.Sprintf("Sent message successfully"),
				})
			} else {
				c.JSON(500, gin.H{
					"error":   err,
					"message": "Server Error",
				})
			}
		} else {
			c.JSON(400, gin.H{
				"error":   "JSON validation error",
				"message": "need at least one channel! (ARRAY)",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"message": "Cannot parse JSON",
			"error":   fmt.Sprintf("%s", err),
		})
		fmt.Printf("%s\n", err)
	}
}

func version_v1(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": globals.VERSION,
	})
}
