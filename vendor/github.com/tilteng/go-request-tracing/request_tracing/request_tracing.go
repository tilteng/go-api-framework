package request_tracing

import (
	"context"
	"net/http"

	"github.com/tilteng/go-logger/logger"
)

var traceIDHeaders = []string{
	"X-Trace-Id",
	"X-Request-Id",
	"X-Crowdtilt-Requestid",
}

type RequestTraceManager interface {
	NewRequestTraceFromHTTPRequest(*http.Request) RequestTrace
	NewEmptyRequestTrace() RequestTrace
	ContextWithRequestTrace(context.Context, RequestTrace) context.Context
	RequestTraceFromContext(context.Context) RequestTrace
	SetBaseLogger(logger.Logger) RequestTraceManager
	Logger() logger.CtxLogger
}

type RequestTrace interface {
	logger.Logger
	// Returns the span id. This ID should be unique per request.
	GetSpanID() string
	// Returns the trace id.
	GetTraceID() string
	// Returns the original SpanID or the current SpanID if there was no
	// original
	GetOriginalSpanID() string
}

type SpanIDGeneratorFn func() string

func (self SpanIDGeneratorFn) GenID() string {
	return self()
}

type SpanIDGenerator interface {
	GenID() string
}

type requestTraceManager struct {
	requestTraceCtxKey *contextKey
	baseLogger         logger.Logger
	spanIDGenerator    SpanIDGenerator
}

func (self *requestTraceManager) SetBaseLogger(logger logger.Logger) RequestTraceManager {
	self.baseLogger = logger
	return self
}

func (self *requestTraceManager) Logger() logger.CtxLogger {
	return &requestTraceLogger{
		baseLogger:          self.baseLogger,
		requestTraceManager: self,
	}
}

// Implements log.Logger
type requestTrace struct {
	baseLogger     logger.Logger
	traceID        string
	spanID         string
	originalSpanID string
}

func (self *requestTrace) logPrefix() string {
	return self.GetSpanID() + " " + self.GetTraceID()
}

func (self *requestTrace) LogDebug(v ...interface{}) {
	self.baseLogger.LogDebug(prependString(self.logPrefix(), v)...)
}

func (self *requestTrace) LogDebugf(fmt string, v ...interface{}) {
	self.baseLogger.LogDebugf(self.logPrefix()+" "+fmt, v...)
}

func (self *requestTrace) LogError(v ...interface{}) {
	self.baseLogger.LogError(prependString(self.logPrefix(), v)...)
}

func (self *requestTrace) LogErrorf(fmt string, v ...interface{}) {
	self.baseLogger.LogErrorf(self.logPrefix()+" "+fmt, v...)
}

func (self *requestTrace) LogInfo(v ...interface{}) {
	self.baseLogger.LogInfo(prependString(self.logPrefix(), v)...)
}

func (self *requestTrace) LogInfof(fmt string, v ...interface{}) {
	self.baseLogger.LogInfof(self.logPrefix()+" "+fmt, v...)
}

func (self *requestTrace) LogWarn(v ...interface{}) {
	self.baseLogger.LogWarn(prependString(self.logPrefix(), v)...)
}

func (self *requestTrace) LogWarnf(fmt string, v ...interface{}) {
	self.baseLogger.LogWarnf(self.logPrefix()+" "+fmt, v...)
}

func (self *requestTrace) GetTraceID() string {
	return self.traceID
}

func (self *requestTrace) GetSpanID() string {
	return self.spanID
}

func (self *requestTrace) GetOriginalSpanID() string {
	return self.originalSpanID
}

func (self *requestTraceManager) NewEmptyRequestTrace() RequestTrace {
	return &requestTrace{
		spanID:         "-",
		traceID:        "-",
		originalSpanID: "-",
	}
}

func (self *requestTraceManager) NewRequestTraceFromHTTPRequest(req *http.Request) RequestTrace {
	hdrs := req.Header

	span_id := self.spanIDGenerator.GenID()

	orig_span_id := hdrs.Get("X-Span-Id")
	if len(orig_span_id) == 0 {
		orig_span_id = span_id
	}

	var trace_id string

	for _, hdr := range traceIDHeaders {
		trace_id = hdrs.Get(hdr)
		if len(trace_id) != 0 {
			break
		}
	}

	if len(trace_id) == 0 {
		trace_id = span_id
	}

	return &requestTrace{
		baseLogger:     self.baseLogger,
		traceID:        trace_id,
		spanID:         span_id,
		originalSpanID: orig_span_id,
	}
}

func NewRequestTraceManager() RequestTraceManager {
	return &requestTraceManager{
		requestTraceCtxKey: &contextKey{"requestTrace"},
		baseLogger:         logger.DefaultStdoutLogger(),
		spanIDGenerator:    SpanIDGeneratorFn(defaultIDGenerator),
	}
}

func NewRequestTraceFromHTTPRequest(req *http.Request) RequestTrace {
	return defaultTraceManager.NewRequestTraceFromHTTPRequest(req)
}
