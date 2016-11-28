package request_tracing

import (
	"context"

	"github.com/tilteng/go-logger/logger"
)

func prependString(s string, v []interface{}) []interface{} {
	nv := make([]interface{}, 1+len(v), 1+len(v))
	nv[0] = s
	copy(nv[1:], v)
	return nv
}

// Implements logger.CtxLogger
type requestTraceLogger struct {
	requestTraceManager RequestTraceManager
	baseLogger          logger.Logger
}

func (self *requestTraceLogger) logPrefix(ctx context.Context) string {
	rt := self.requestTraceManager.RequestTraceFromContext(ctx)
	if rt == nil {
		return "- -"
	}
	return rt.GetSpanID() + " " + rt.GetTraceID()
}

func (self *requestTraceLogger) LogDebug(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogDebug(prependString(self.logPrefix(ctx), v)...)
}

func (self *requestTraceLogger) LogDebugf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogDebugf(self.logPrefix(ctx)+" "+fmt, v...)
}

func (self *requestTraceLogger) LogError(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogError(prependString(self.logPrefix(ctx), v)...)
}

func (self *requestTraceLogger) LogErrorf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogErrorf(self.logPrefix(ctx)+" "+fmt, v...)
}

func (self *requestTraceLogger) LogInfo(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogInfo(prependString(self.logPrefix(ctx), v)...)
}

func (self *requestTraceLogger) LogInfof(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogInfof(self.logPrefix(ctx)+" "+fmt, v...)
}

func (self *requestTraceLogger) LogWarn(ctx context.Context, v ...interface{}) {
	self.baseLogger.LogWarn(prependString(self.logPrefix(ctx), v)...)
}

func (self *requestTraceLogger) LogWarnf(ctx context.Context, fmt string, v ...interface{}) {
	self.baseLogger.LogWarnf(self.logPrefix(ctx)+" "+fmt, v...)
}

func (self *requestTraceLogger) BaseLogger() logger.Logger {
	return self.baseLogger
}
