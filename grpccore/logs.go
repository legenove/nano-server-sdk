package grpccore

import (
	"context"
	"fmt"
	"github.com/legenove/nano-server-sdk/servers"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

func LoggerRecoveryHandler(funcName string, handler grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		// before
		start := time.Now()
		defer func() {
			var reason interface{}
			var error_code interface{}
			if err := recover(); err != nil {
				duration := time.Since(start)
				logDir := servers.LogDirError
				switch err.(type) {
				case *servers.ServerError:
					_err := err.(*servers.ServerError)
					reason = _err.Error()
					error_code = _err.Code
					// 定义的error  在warn日志中
					logDir = servers.LogDirWarn
				case error:
					reason = err.(error).Error()
					if errInfo, ok := servers.ServerErrorMap[reason.(string)]; ok {
						_err := errInfo
						reason = _err.Error()
						error_code = _err.Code
						// 定义的error 在warn日志中
						logDir = servers.LogDirWarn
					} else {
						_stack := stack(3)
						reason = fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n", err, _stack)
						error_code = "10001"
					}
				default:
					error_code = "10000"
					if _, ok := err.(string); ok {
						reason = err.(string)
					} else {
						reason = servers.ErrUnKnowRequest.Msg
					}
				}

				// 未定义的错误，在error中， 定义的错误在warn中
				servers.WarnLog(logDir, ctx, error_code, reason, duration)
			}
		}()
		ctx = servers.InitContext(ctx, funcName, req)
		res, err := handler(ctx, req)
		// after
		// Log only when path is not being skipped
		var accesslog bool
		if servers.OpenAccessLog <= 0 {
			accesslog = false
		} else if servers.OpenAccessLog >= 100 {
			accesslog = true
		} else if rand.Int()/100 < servers.OpenAccessLog {
			accesslog = true
		}
		if accesslog {
			duration := time.Since(start)
			servers.AccessLog(servers.LogDirAccess, ctx, duration)
		}
		return res, err
	}
}
