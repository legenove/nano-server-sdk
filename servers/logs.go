package servers

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	raw := GetRequestRaw(ctx)
	logger.Info("access",
		zap.String("log_type", LOG_TYPE_APP_ACCESS),
		zap.String("event", LogEventAccess),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(ctx, raw)),
		zap.String("requestFunc", GetServerRequestFunc(ctx, raw)),
		zap.String("fromApp", GetServerName(ctx, md)),
		zap.String("fromProject", GetServerGroup(ctx, md)),
		zap.String("requestId", GetRequestId(ctx, md)),
		zap.String("clientIp", GetContextIP(ctx, md)),
		zap.Namespace("properties"),
		// TODO 增加query
		//zap.String("query", query),
		zap.String("user-agent", GetUserAgent(ctx, md)),
		zap.Duration("time", duration),
	)
}

func ErrorLog(logger *zap.Logger, ctx context.Context, error_code, reason interface{}, duration time.Duration) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	raw := GetRequestRaw(ctx)
	logger.Error("error",
		zap.String("log_type", LOG_TYPE_APP_ERROR),
		zap.String("event", LogEventError),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(ctx, raw)),
		zap.String("requestFunc", GetServerRequestFunc(ctx, raw)),
		zap.String("fromApp", GetServerName(ctx, md)),
		zap.String("fromProject", GetServerGroup(ctx, md)),
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("clientIp", GetContextIP(ctx, md)),
		zap.Namespace("properties"),
		zap.Reflect("error_code", error_code),
		zap.String("query", GetServerRequestInfo(ctx)),
		zap.String("user-agent", GetUserAgent(ctx, md)),
		zap.Duration("time", duration),
		zap.Reflect("reason", reason))
}

func WarnLog(logger *zap.Logger, ctx context.Context, error_code, reason interface{}, duration time.Duration) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	raw := GetRequestRaw(ctx)
	logger.Warn("warning",
		zap.String("log_type", LOG_TYPE_APP_WARN),
		zap.String("event", LogEventError),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(ctx, raw)),
		zap.String("requestFunc", GetServerRequestFunc(ctx, raw)),
		zap.String("fromApp", GetServerName(ctx, md)),
		zap.String("fromProject", GetServerGroup(ctx, md)),
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("clientIp", GetContextIP(ctx, md)),
		zap.Namespace("properties"),
		zap.Reflect("error_code", error_code),
		zap.String("query", GetServerRequestInfo(ctx)),
		zap.String("user-agent", GetUserAgent(ctx, md)),
		zap.Duration("time", duration),
		zap.Reflect("reason", reason))
}

func AddRequestLog(logger *zap.Logger, ctx context.Context) *zap.Logger {
	if ctx != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		raw := GetRequestRaw(ctx)
		return logger.With([]zapcore.Field{
			zap.String("logServer", Server.GetServerName()),
			zap.String("logServerGroup", Server.GetServerGroup()),
			zap.String("requestType", GetServerRequestType(ctx, raw)),
			zap.String("requestFunc", GetServerRequestFunc(ctx, raw)),
			zap.String("fromApp", GetServerName(ctx, md)),
			zap.String("fromProject", GetServerGroup(ctx, md)),
			zap.String("requestId", GetRequestId(ctx)),
			zap.String("clientIp", GetContextIP(ctx, md)),
		}...)
	} else {
		return logger.With([]zapcore.Field{
			zap.String("logServer", Server.GetServerName()),
			zap.String("logServerGroup", Server.GetServerGroup()),
			zap.String("requestType", ""),
			zap.String("requestFunc", ""),
			zap.String("fromApp", ""),
			zap.String("fromProject", ""),
			zap.String("requestId", ""),
			zap.String("clientIp", ""),
		}...)
	}
}

func InterleavedKVToFields(log_act string, keyValues ...interface{}) ([]zapcore.Field, error) {
	if len(keyValues)%2 != 0 {
		return nil, fmt.Errorf("non-even keyValues len: %d", len(keyValues))
	}
	fields := make([]zapcore.Field, len(keyValues)/2+1)
	fields[0] = zap.Namespace("properties")
	for i := 0; i*2 < len(keyValues); i++ {
		key, ok := keyValues[i*2].(string)
		if key != "log_act" {
			key = utils.ConcatenateStrings(log_act, "_", key)
		}
		if !ok {
			return nil, fmt.Errorf(
				"non-string key (pair #%d): %T",
				i, keyValues[i*2])
		}
		switch typedVal := keyValues[i*2+1].(type) {
		case bool:
			fields[i+1] = zap.Bool(key, typedVal)
		case string:
			fields[i+1] = zap.String(key, typedVal)
		case int:
			fields[i+1] = zap.Int(key, typedVal)
		case int8:
			fields[i+1] = zap.Int8(key, typedVal)
		case int16:
			fields[i+1] = zap.Int16(key, typedVal)
		case int32:
			fields[i+1] = zap.Int32(key, typedVal)
		case int64:
			fields[i+1] = zap.Int64(key, typedVal)
		case uint:
			fields[i+1] = zap.Uint(key, typedVal)
		case uint64:
			fields[i+1] = zap.Uint64(key, typedVal)
		case uint8:
			fields[i+1] = zap.Uint32(key, uint32(typedVal))
		case uint16:
			fields[i+1] = zap.Uint32(key, uint32(typedVal))
		case uint32:
			fields[i+1] = zap.Uint32(key, typedVal)
		case float32:
			fields[i+1] = zap.Float32(key, typedVal)
		case float64:
			fields[i+1] = zap.Float64(key, typedVal)
		default:
			// When in doubt, coerce to a string
			fields[i+1] = zap.Reflect(key, typedVal)
		}
	}
	return fields, nil
}

func getEventString(event string) string {
	return utils.ConcatenateStrings(LogDirOther, "_", event)
}

/*
 * LogKV is a concise, readable way to record key:value logging data about
 * a Span, though unfortunately this also makes it less efficient and less
 * type-safe than ZapField(). Here's an example:
 *    servers.LogKV( core.LOG_LEVEL_INFO,"message", "finan_add", ctx, key1, val1, key2, val2, key3, val3, ...)
 */
func LogKV(level, message string, logAct string, ctx context.Context, options ...interface{}) error {
	if !utils.CheckNormalKey(logAct) {
		return fmt.Errorf("log_act:%s, only support a-zA-Z0-9 '-_' and  must startwith a-zA-Z", logAct)
	}
	eventString := getEventString(logAct)
	var err error
	logger, err := cocore.LogPool.Instance(eventString)
	if err != nil {
		return err
	}
	logger = logger.With([]zapcore.Field{
		zap.String("log_type", LOG_TYPE_OTHER),
		zap.String("event", eventString)}...)
	logger = AddRequestLog(logger, ctx)
	options = append(options, "log_act", logAct)
	var fields []zapcore.Field
	fields, err = InterleavedKVToFields(logAct, options...)
	if err != nil {
		return err
	}
	switch level {
	case cocore.LOG_LEVEL_DEBUG:
		logger.With(fields...).Debug(message)
	case cocore.LOG_LEVEL_INFO:
		logger.With(fields...).Info(message)
	case cocore.LOG_LEVEL_WARN:
		logger.With(fields...).Warn(message)
	case cocore.LOG_LEVEL_ERROR:
		logger.With(fields...).Error(message)
		err = logger.Sync()
	}
	return err
}
