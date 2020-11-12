package servers

import (
	"context"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"

	"github.com/legenove/cocore"
	"go.uber.org/zap"
)

type LogWriter interface {
	Write() error
	Put()
}

type BaseWriter struct {
	logName string
	msg     string
	ctx     context.Context
}

func (a *BaseWriter) getLogger() (*zap.Logger, error) {
	return cocore.LogPool.Instance(a.logName)
}

/**
access log
*/
var AccessWriterPool = sync.Pool{
	New: func() interface{} {
		return &AccessLogWriter{}
	},
}

func GetAccessLogWriter() *AccessLogWriter {
	return AccessWriterPool.Get().(*AccessLogWriter)
}

type AccessLogWriter struct {
	BaseWriter
	duration time.Duration
}

func (a *AccessLogWriter) Put() {
	AccessWriterPool.Put(a)
}

func (a *AccessLogWriter) Write() error {
	logger, err := a.getLogger()
	if err != nil {
		return err
	}
	md, ok := metadata.FromIncomingContext(a.ctx)
	if !ok {
		md = metadata.MD{}
	}
	raw := GetRequestRaw(a.ctx)
	logger.Info("access",
		zap.String("log_type", LOG_TYPE_APP_ACCESS),
		zap.String("event", LogEventAccess),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(a.ctx, raw)),
		zap.String("requestFunc", GetServerRequestFunc(a.ctx, raw)),
		zap.String("fromApp", GetServerName(a.ctx, md)),
		zap.String("fromProject", GetServerGroup(a.ctx, md)),
		zap.String("requestId", GetRequestId(a.ctx, md)),
		zap.String("clientIp", GetContextIP(a.ctx, md)),
		zap.Namespace("properties"),
		// TODO 增加query
		//zap.String("query", query),
		zap.String("user-agent", GetUserAgent(a.ctx, md)),
		zap.Duration("time", a.duration),
	)
	return err
}

/**
access log
*/
var ErrorWriterPool = sync.Pool{
	New: func() interface{} {
		return &ErrorLogWriter{}
	},
}

func GetErrorWriter() *ErrorLogWriter {
	return ErrorWriterPool.Get().(*ErrorLogWriter)
}

type ErrorLogWriter struct {
	BaseWriter
	duration  time.Duration
	reason    interface{}
	errorCode interface{}
}

func (a *ErrorLogWriter) Put() {
	AccessWriterPool.Put(a)
}

func (a *ErrorLogWriter) Write() error {
	logger, err := a.getLogger()
	if err != nil {
		return err
	}
	md, ok := metadata.FromIncomingContext(a.ctx)
	if !ok {
		md = metadata.MD{}
	}
	raw := GetRequestRaw(a.ctx)
	logger.Warn("warning",
		zap.String("log_type", LOG_TYPE_APP_WARN),
		zap.String("event", LogEventError),
		zap.String("logServer", Server.GetServerName()),
		zap.String("logServerGroup", Server.GetServerGroup()),
		zap.String("requestType", GetServerRequestType(a.ctx, raw)),
		zap.String("requestFunc", GetServerRequestFunc(a.ctx, raw)),
		zap.String("fromApp", GetServerName(a.ctx, md)),
		zap.String("fromProject", GetServerGroup(a.ctx, md)),
		zap.String("requestId", GetRequestId(a.ctx)),
		zap.String("clientIp", GetContextIP(a.ctx, md)),
		zap.Namespace("properties"),
		zap.Reflect("error_code", a.errorCode),
		zap.String("query", GetServerRequestInfo(a.ctx)),
		zap.String("user-agent", GetUserAgent(a.ctx, md)),
		zap.Duration("time", a.duration),
		zap.Reflect("reason", a.reason))
	return nil
}
