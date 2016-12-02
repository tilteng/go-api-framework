package api_framework

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-panichandler/panichandler_mw"
	"github.com/tilteng/go-api-request-logger/request_logger_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-app-context/app_context"
	"github.com/tilteng/go-errors/errors"
	"github.com/tilteng/go-logger/apache_logger_mw"
	"github.com/tilteng/go-logger/logger"
	"github.com/tilteng/go-metrics/metrics_mw"
	"github.com/tilteng/go-request-tracing/request_tracing"
)

type appContext app_context.AppContext

type ErrorFormatterFn func(context.Context, errors.ErrorType) interface{}

func (self ErrorFormatterFn) FormatErrors(ctx context.Context, err errors.ErrorType) interface{} {
	return self(ctx, err)
}

type ErrorFormatter interface {
	FormatErrors(context.Context, errors.ErrorType) interface{}
}

type ControllerOpts struct {
	// Used only for Link: header responses for json schema
	BaseAPIURL             string
	BaseRouter             *api_router.Router
	ConsumesContent        []string
	ProducesContent        []string
	JSONSchemaRoutePath    string
	JSONSchemaFilePath     string
	JSONSchemaErrorHandler jsonschema_mw.ErrorHandler
	PanicHandler           panichandler_mw.PanicHandler
	SerializerErrorHandler serializers_mw.ErrorHandler
	ApacheLogWriter        io.Writer
	ApacheLogCombined      bool
	ErrorFormatter         ErrorFormatter
	RequestTraceManager    request_tracing.RequestTraceManager
	RequestLoggerOpts      *request_logger_mw.RequestLoggerOpts

	// We pull metrics, rollbar, and logger from AppContext
	AppContext app_context.AppContext
}

type Controller struct {
	*api_router.Router
	appContext
	logger                  logger.CtxLogger
	options                 *ControllerOpts
	errorFormatter          ErrorFormatter
	requestTraceManager     request_tracing.RequestTraceManager
	JSONSchemaMiddleware    *jsonschema_mw.JSONSchemaMiddleware
	PanicHandlerMiddleware  *panichandler_mw.PanicHandlerMiddleware
	SerializerMiddleware    *serializers_mw.SerializerMiddleware
	ApacheLoggerMiddleware  *apache_logger_mw.ApacheLoggerMiddleware
	MetricsMiddleware       *metrics_mw.MetricsMiddleware
	RequestLoggerMiddleware *request_logger_mw.RequestLoggerMiddleware
}

func (self *Controller) GenUUID() *UUID {
	return GenUUID()
}

func (self *Controller) GenUUIDHex() string {
	return GenUUIDHex()
}

func (self *Controller) UUIDFromString(s string) *UUID {
	return UUIDFromString(s)
}

func (self *Controller) JSONSchemaOpts(name string) *jsonschema_mw.JSONSchemaOpts {
	return &jsonschema_mw.JSONSchemaOpts{Name: name}
}

func (self *Controller) ReadBody(ctx context.Context, v interface{}) error {
	rctx := self.RequestContext(ctx)
	return rctx.serializerRequestContext.ReadDeserializedBody(rctx, v)
}

func (self *Controller) WriteResponse(ctx context.Context, v interface{}) error {
	rctx := self.RequestContext(ctx)
	if tilterr, ok := v.(errors.ErrorType); ok {
		status := tilterr.GetStatus()
		rctx.SetStatus(status)
		v = self.errorFormatter.FormatErrors(ctx, tilterr)
	}

	rctx.SetResponseHeader(
		"X-Response-Time",
		fmt.Sprintf("%f ms",
			float64(rctx.TimeElapsed())/float64(time.Millisecond),
		),
	)

	return rctx.serializerRequestContext.WriteSerializedResponse(rctx, v)
}

func RequestContextFromContext(ctx context.Context) *RequestContext {
	controller_rctx, _ := ctx.Value(requestContextCtxKey).(*RequestContext)
	return controller_rctx
}
