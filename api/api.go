package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/plugins"
)

func Listen(iface string) {
	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": globals.Version,
		})
	})

	r.POST("/register/:plugin", func(c *gin.Context) {
		var p plugins.Plugin
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
	})
	r.Run(":8080")
}
