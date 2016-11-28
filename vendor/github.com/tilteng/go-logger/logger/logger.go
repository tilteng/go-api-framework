package logger

import (
	"context"
	"io"
	"log"
)

type Logger interface {
	LogDebug(v ...interface{})
	LogDebugf(fmt string, v ...interface{})
	LogError(v ...interface{})
	LogErrorf(fmt string, v ...interface{})
	LogInfo(v ...interface{})
	LogInfof(fmt string, v ...interface{})
	LogWarn(v ...interface{})
	LogWarnf(fmt string, v ...interface{})
}

type CtxLogger interface {
	LogDebug(ctx context.Context, v ...interface{})
	LogDebugf(ctx context.Context, fmt string, v ...interface{})
	LogError(ctx context.Context, v ...interface{})
	LogErrorf(ctx context.Context, fmt string, v ...interface{})
	LogInfo(ctx context.Context, v ...interface{})
	LogInfof(ctx context.Context, fmt string, v ...interface{})
	LogWarn(ctx context.Context, v ...interface{})
	LogWarnf(ctx context.Context, fmt string, v ...interface{})
	BaseLogger() Logger
}

type defaultLogger struct {
	logger *log.Logger
}

type defaultCtxLogger struct {
	baseLogger Logger
}

func prependString(s string, v []interface{}) []interface{} {
	nv := make([]interface{}, 1+len(v), 1+len(v))
	nv[0] = s
	copy(nv[1:], v)
	return nv
}

func (self *defaultLogger) LogDebug(v ...interface{}) {
	self.logger.Println(prependString("[DEBUG]", v)...)
}

func (self *defaultLogger) LogDebugf(fmt string, v ...interface{}) {
	self.logger.Printf("[DEBUG] "+fmt, v...)
}

func (self *defaultLogger) LogError(v ...interface{}) {
	self.logger.Println(prependString("[ERROR]", v)...)
}

func (self *defaultLogger) LogErrorf(fmt string, v ...interface{}) {
	self.logger.Printf("[ERROR] "+fmt, v...)
}

func (self *defaultLogger) LogInfo(v ...interface{}) {
	self.logger.Println(prependString("[INFO]", v)...)
}

func (self *defaultLogger) LogInfof(fmt string, v ...interface{}) {
	self.logger.Printf("[INFO] "+fmt, v...)
}

func (self *defaultLogger) LogWarn(v ...interface{}) {
	self.logger.Println(prependString("[WARN]", v)...)
}

func (self *defaultLogger) LogWarnf(fmt string, v ...interface{}) {
	self.logger.Printf("[WARN] "+fmt, v...)
}

func (self *defaultLogger) SetWriter(out io.Writer) {
	self.logger.SetOutput(out)
}

func (self *defaultCtxLogger) LogDebug(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogDebug(v...)
}

func (self *defaultCtxLogger) LogDebugf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogDebugf(fmt, v...)
}

func (self *defaultCtxLogger) LogError(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogError(v...)
}

func (self *defaultCtxLogger) LogErrorf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogErrorf(fmt, v...)
}

func (self *defaultCtxLogger) LogInfo(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogInfo(v...)
}

func (self *defaultCtxLogger) LogInfof(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogInfof(fmt, v...)
}

func (self *defaultCtxLogger) LogWarn(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogWarn(v...)
}

func (self *defaultCtxLogger) LogWarnf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogWarnf(fmt, v...)
}

func (self *defaultCtxLogger) BaseLogger() Logger {
	return self.baseLogger
}

func NewDefaultLogger(out io.Writer, prefix string) *defaultLogger {
	return &defaultLogger{
		logger: log.New(out, prefix, log.LstdFlags|log.Lmicroseconds),
	}
}

func NewDefaultCtxLogger(base_logger Logger) *defaultCtxLogger {
	if base_logger == nil {
		base_logger = defaultStdoutLogger
	}
	return &defaultCtxLogger{
		baseLogger: base_logger,
	}
}
