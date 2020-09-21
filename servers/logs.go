package servers

import (
	"context"
	"strconv"
	"time"

	"github.com/legenove/cocore"
	"github.com/legenove/utils"
	"go.uber.org/zap"
)

var (
	LogDirError   string
	LogDirWarn    string
	LogDirAccess  string
	LogDirMysql   string
	LogDirRedis   string
	LogDirRequest string
	LogDirAsync   string
	LogDirSub     string
	LogDirOther   string

	// 开放access日志比例，默认全部开放
	OpenAccessLog int
)

var (
	LogEventAccess  string
	LogEventError   string
	LogEventMysql   string
	LogEventRedis   string
	LogEventRequest string
)

const (
	LOG_TYPE_APP_ACCESS = "access"
	LOG_TYPE_APP_ERROR  = "error"
	LOG_TYPE_APP_WARN   = "warning"
	LOG_TYPE_MYSQL      = "mysql"
	LOG_TYPE_REDIS      = "redis"
	LOG_TYPE_ASYNC      = "async"     // 异步任务
	LOG_TYPE_SUB        = "subscribe" // 订阅任务
	LOG_TYPE_REQUEST    = "request"   // 请求日志
	LOG_TYPE_OTHER      = "project"   // 业务日志
)

func InitServerLog() {
	// 开放access日志比例
	openAccessLog := cocore.App.GetStringConfig("OPEN_ACCESS_LOG", "100")
	i, err := strconv.Atoi(openAccessLog)
	if err != nil {
		i = 100
	}
	if i > 100 {
		OpenAccessLog = 100
	} else if i <= 0 {
		OpenAccessLog = 0
	} else {
		OpenAccessLog = i
	}
	LogDirError = cocore.App.GetStringConfig("ERROR_LOG_NAME", LOG_TYPE_APP_ERROR)
	LogDirAccess = cocore.App.GetStringConfig("ACCESS_LOG_NAME", LOG_TYPE_APP_ACCESS)
	LogDirWarn = cocore.App.GetStringConfig("WARN_LOG_NAME", LOG_TYPE_APP_WARN)
	LogDirMysql = cocore.App.GetStringConfig("MYSQL_LOG_NAME", LOG_TYPE_MYSQL)
	LogDirRedis = cocore.App.GetStringConfig("REDIS_LOG_NAME", LOG_TYPE_REDIS)
	LogDirAsync = cocore.App.GetStringConfig("ASYNC_LOG_NAME", LOG_TYPE_ASYNC)
	LogDirSub = cocore.App.GetStringConfig("SUBSCRIBE_LOG_NAME", LOG_TYPE_SUB)
	LogDirRequest = cocore.App.GetStringConfig("REQUEST_LOG_NAME", LOG_TYPE_REQUEST)
	LogDirOther = utils.ConcatenateStrings(Server.GetServerGroup(), "_", Server.GetServerName())
	LogEventAccess = utils.ConcatenateStrings(LogDirOther, "_", LOG_TYPE_APP_ACCESS)
	LogEventError = utils.ConcatenateStrings(LogDirOther, "_", LOG_TYPE_APP_ERROR)
	LogEventMysql = utils.ConcatenateStrings(LogDirOther, "_", LOG_TYPE_MYSQL)
	LogEventRedis = utils.ConcatenateStrings(LogDirOther, "_", LOG_TYPE_REDIS)
	LogEventRequest = utils.ConcatenateStrings(LogDirOther, "_", LOG_TYPE_REQUEST)
}

func AccessLog(logger *zap.Logger, ctx context.Context, duration time.Duration) {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	logger.Info("access",
		zap.String("log_type", LOG_TYPE_APP_ACCESS),
		zap.String("event", LogEventAccess),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(ctx)),
		zap.String("requestFunc", GetServerRequestFunc(ctx)),
		zap.String("fromApp", GetServerName(ctx)),
		zap.String("fromProject", GetServerGroup(ctx)),
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("clientIp", GetContextIp(ctx)),
		zap.Namespace("properties"),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.Duration("time", duration),
	)
}

func ErrorLog() {

}

func WarningLog() {

}

func LogKV() {

}
