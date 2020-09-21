package gincore

import (
	"github.com/gin-gonic/gin"
	"github.com/legenove/nano-server-sdk/servers"
)

var router *gin.Engine

type HandlerDecorator func(handlerFunc gin.HandlerFunc) gin.HandlerFunc

func PingApi(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetRouter() *gin.Engine {
	if router == nil {
		//gin.Logger()
		r := gin.New()
		r.Use(LoggerRecovery())
		// regist error
		r.NoRoute(func(c *gin.Context) {
			c.JSON(404, servers.ErrPageNotFoundRequest)
		})
		r.NoMethod(func(c *gin.Context) {
			c.JSON(405, servers.ErrMethodNotAllowRequest)
		})
		router = r
	}
	return router
}

func GetCurrentGroup(relativePath string) (group *gin.RouterGroup) {
	r := GetRouter()
	group = r.Group(relativePath)
	// group add midderware
	return
}
