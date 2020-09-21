package grpccore

import (
	"context"
	"fmt"
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
				logDir := LogDirError
				switch err.(type) {
				case *APIError:
					_err := err.(*APIError)
					reason = _err.Error()
					error_code = _err.Code
					c.JSON(_err.StatusCode(), _err)
					// 定义的error 在access日志中
					logDir = LogDirWarn
				case error:
					reason = err.(error).Error()
					if errInfo, ok := appErrorMap[reason.(string)]; ok {
						_err := NewAPIError(errInfo...)
						reason = _err.Error()
						error_code = _err.Code
						c.JSON(_err.StatusCode(), _err)
						// 定义的error 在access日志中
						logDir = LogDirWarn
					} else {
						_stack := stack(3)
						reason = fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n", err, _stack)
						error_code = "10001"
						if App.DEBUG {
							c.JSON(400, NewAPIError(reason.(string), "10001", "400"))
						} else {
							c.JSON(400, ErrUnKnowRequest)
						}
					}
				default:
					error_code = "10000"
					if _, ok := err.(string); ok {
						reason = err.(string)
					} else {
						reason = ErrUnKnowRequest.Msg
					}
					c.JSON(400, ErrUnKnowRequest.New([]string{reason.(string)}))

				}

				// 未定义的错误，在error中， 定义的错误在Access中
				elog, _ := LogPool.Instance(logDir)
				WarnLog(elog, c, error_code, reason, duration)
			}
		}()
		res, err := handler(ctx, req)
		// after
		return res, err
	}
}
