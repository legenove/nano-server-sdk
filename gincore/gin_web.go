package gincore

import "github.com/gin-gonic/gin"

const (
	SERVER_TYPE_REST   = "rest"
	SERVER_TYPE_RPC    = "rpc"
	SERVER_TYPE_TCP    = "tcp"
	SERVER_TYPE_SERVER = "server" // rest and rpc and tcp
	SERVER_TYPE_ASYNC  = "async"
)

var router *gin.Engine

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
			c.JSON(404, ErrPageNotFoundRequest)
		})
		r.NoMethod(func(c *gin.Context) {
			c.JSON(405, ErrMethodNotAllowRequest)
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
