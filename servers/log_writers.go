package servers

import (
	"sync"
	"time"

	"github.com/legenove/cocore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogWriter interface {
	Write() error
	Put()
}

type LogTypeCase int

const (
	RequestAccessTypeCase LogTypeCase = iota + 1
	RequestErrorTypeCase
	RequestWarnTypeCase
)

type BaseWriter struct {
	logName        string
	logLevel       zapcore.Level
	msg            string
	logType        string
	event          string
	logServer      string
	logServerGroup string
	requestType    string
	requestFunc    string
	fromApp        string
	fromProject    string
	requestId      string
	clientIp       string
}

func (a *BaseWriter) getLogger() (*zap.Logger, error) {
	return cocore.LogPool.Instance(a.logName)
}

func (a *BaseWriter) getBaseField(logger *zap.Logger) *zap.Logger {
	return logger.With(
		zap.String("log_type", a.logType),
		zap.String("event", a.event),
		zap.String("logServer", a.logServer),
		zap.String("logServerGroup", a.logServerGroup),
		zap.String("requestType", a.requestType),
		zap.String("requestFunc", a.requestFunc),
		zap.String("fromApp", a.fromApp),
		zap.String("fromProject", a.fromProject),
		zap.String("requestId", a.requestId),
		zap.String("clientIp", a.clientIp),
		zap.Namespace("properties"),
	)
}

var RequestWriterPool = sync.Pool{
	New: newRequestWriter,
}

func GetRequestWriter() *RequestWriter {
	return RequestWriterPool.Get().(*RequestWriter)
}

func newRequestWriter() interface{} {
	return &RequestWriter{}
}

type RequestWriter struct {
	BaseWriter
	LogCase   LogTypeCase
	userAgent string
	errorCode interface{}
	query     string
	duration  time.Duration
	reason    interface{}
}

func (a *RequestWriter) Put() {
	RequestWriterPool.Put(a)
}

func (a *RequestWriter) Write() error {
	logger, err := a.getLogger()
	if err != nil {
		return err
	}
	logger.With(
		zap.String("user-agent", a.userAgent),
		zap.Duration("duration", a.duration),
	)
	if a.LogCase != RequestAccessTypeCase {
		logger.With(
			zap.Reflect("error_code", a.errorCode),
			zap.String("query", a.query),
			zap.Reflect("reason", a.reason),
		)
	}
	switch a.logLevel {
	case zap.DebugLevel:
		logger.Debug(a.msg)
	case zap.InfoLevel:
		logger.Info(a.msg)
	case zap.WarnLevel:
		logger.Warn(a.msg)
	case zap.ErrorLevel:
		logger.Error(a.msg)
		err = logger.Sync()
	default:
		logger.Error(a.msg)

	}
	return err
}
