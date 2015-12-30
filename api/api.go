package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sethdmoore/digo/globals"
)

func Listen(iface string) {
	r := gin.Default()
	r.Listen()
}
