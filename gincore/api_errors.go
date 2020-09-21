package gincore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/legenove/cocore"
	"github.com/legenove/nano-server-sdk/servers"
	"io"
	"os"
)

var DefaultWriter io.Writer = os.Stdout

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func LoggerRecovery() gin.HandlerFunc {
	return LoggerWithWriter("/metric")
}

// LoggerWithWriter instance a Logger middleware with the specified writer buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(notlogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		defer func() {
			var reason interface{}
			if err := recover(); err != nil {
				switch err.(type) {
				case *servers.ServerError:
					_err := err.(*servers.ServerError)
					reason = _err.Error()
					c.JSON(_err.StatusCode(), _err)
				case error:
					reason = err.(error).Error()
					if _err, ok := servers.ServerErrorMap[reason.(string)]; ok {
						reason = _err.Error()
						c.JSON(_err.StatusCode(), _err)
					} else {
						_stack := stack(3)
						reason = fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n", err, _stack)
						if cocore.App.DEBUG {
							c.JSON(400, servers.NewServerError(reason.(string), "10001", 400))
						} else {
							c.JSON(400, servers.ErrUnKnowRequest)
						}
					}
				default:
					if _, ok := err.(string); ok {
						reason = err.(string)
					} else {
						reason = servers.ErrUnKnowRequest.Msg
					}
					c.JSON(400, servers.ErrUnKnowRequest.New([]string{reason.(string)}))

				}
			}
		}()

		c.Next()
	}
}
