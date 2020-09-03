package gincore

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/legenove/cocore"
	"github.com/legenove/utils"
)

type APIError struct {
	Status     string   `json:"status"`
	statusCode int      `json:"-"`
	Code       string   `json:"code"`
	Msg        string   `json:"msg"`
	Details    []string `json:"details"`
}

func (a *APIError) Error() string {
	if len(a.Details) == 0 {
		return a.Msg
	}
	out := make([]string, len(a.Details)+2)
	out[0] = a.Msg
	out[1] = " : "
	copy(out[2:], a.Details[:])
	return utils.ConcatenateStrings(out...)
}

var appErrorMap = make(map[string]*APIError)

func (apiError *APIError) New(details []string, error_code ...string) *APIError {
	var code string
	if len(error_code) > 0 {
		code = error_code[0]
	} else {
		code = apiError.Code
	}
	return &APIError{
		Status:     apiError.Status,
		statusCode: apiError.statusCode,
		Code:       code,
		Msg:        apiError.Msg,
		Details:    details,
	}
}

func (apiError *APIError) StatusCode() int {
	if apiError.statusCode > 0 {
		return apiError.statusCode
	}
	return 200
}

func (apiError *APIError) SetStatusCode(code int) *APIError {
	apiError.statusCode = code
	return apiError
}

func NewAPIError(args ...string) *APIError {
	msg := args[0]
	code := args[1]
	status := args[2]
	apiErr := &APIError{
		Status:  status,
		Code:    code,
		Msg:     msg,
		Details: []string{},
	}
	appErrorMap[msg] = apiErr
	return apiErr
}

var (
	ErrUnKnowRequest         = NewAPIError("unknow_error", "10000", "400")
	ErrProjectValidator      = NewAPIError("project_validator_error", "10002", "400")
	ErrProjectMatch          = NewAPIError("project_match_error", "10003", "400")
	ErrPageNotFoundRequest   = NewAPIError("not_found", "10004", "404")
	ErrMethodNotAllowRequest = NewAPIError("no_method", "10005", "405")
	ErrSchemaOptionNotFound  = NewAPIError("unknow_error", "10006", "400")
	ErrUnDefineRequest       = NewAPIError("undefined_error", "10007", "400")
	ErrRequestErr            = NewAPIError("requests_error", "10008", "400")
	ErrGetRequestHost        = NewAPIError("get_request_host_error", "10009", "400")
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}

var DefaultWriter io.Writer = os.Stdout

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func LoggerRecovery() gin.HandlerFunc {
	return LoggerWithWriter("/metric")
}

// TODO 先把日志删掉，看看grpc层如何做装饰器把日志加上，这样是最全面的。

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
		// Start timer
		//start := time.Now()
		//path := c.Request.URL.Path
		defer func() {
			var reason interface{}
			//var error_code interface{}
			if err := recover(); err != nil {
				//duration := time.Since(start)
				//logDir := LogDirError
				switch err.(type) {
				case *APIError:
					_err := err.(*APIError)
					reason = _err.Error()
					//error_code = _err.Code
					c.JSON(_err.StatusCode(), _err)
					// 定义的error 在access日志中
					//logDir = LogDirWarn
				case error:
					reason = err.(error).Error()
					if _err, ok := appErrorMap[reason.(string)]; ok {
						reason = _err.Error()
						//error_code = _err.Code
						c.JSON(_err.StatusCode(), _err)
						// 定义的error 在access日志中
						//logDir = LogDirWarn
					} else {
						_stack := stack(3)
						reason = fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n", err, _stack)
						//error_code = "10001"
						if cocore.App.DEBUG {
							c.JSON(400, NewAPIError(reason.(string), "10001", "400"))
						} else {
							c.JSON(400, ErrUnKnowRequest)
						}
					}
				default:
					//error_code = "10000"
					if _, ok := err.(string); ok {
						reason = err.(string)
					} else {
						reason = ErrUnKnowRequest.Msg
					}
					c.JSON(400, ErrUnKnowRequest.New([]string{reason.(string)}))

				}

				// 未定义的错误，在error中， 定义的错误在Access中
				//elog, _ := LogPool.Instance(logDir)
				//WarnLog(elog, c, error_code, reason, duration)
			}
		}()

		// Process request
		c.Next()

		//// Log only when path is not being skipped
		//if OpenAccessLog {
		//	if _, ok := skip[path]; !ok {
		//		// Stop timer
		//		//duration := time.Since(start)
		//
		//		statusCode := c.Writer.Status()
		//		if statusCode > 0 && statusCode < 400 {
		//			//LogContextChan <- LogContext{c.Copy(), duration}
		//			//log, _ := LogPool.Instance(LogDirAccess)
		//			//AccessLog(log, c, duration)
		//		}
		//	}
		//}
	}
}
