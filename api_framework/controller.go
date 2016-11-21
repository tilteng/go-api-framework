package api_framework

import (
	"context"
	"io"
	"os"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-panichandler/panichandler_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-errors/errors"
	"github.com/tilteng/go-logger/apache_logger_mw"
	"github.com/tilteng/go-logger/logger"
	"github.com/tilteng/go-metrics/metrics"
	"github.com/tilteng/go-metrics/metrics_mw"
)

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
	MetricsClient          metrics.MetricsClient
	Logger                 logger.Logger
}

type Controller struct {
	*api_router.Router
	logger.Logger
	options                *ControllerOpts
	errorFormatter         ErrorFormatter
	JSONSchemaMiddleware   *jsonschema_mw.JSONSchemaMiddleware
	PanicHandlerMiddleware *panichandler_mw.PanicHandlerMiddleware
	SerializerMiddleware   *serializers_mw.SerializerMiddleware
	ApacheLoggerMiddleware *apache_logger_mw.ApacheLoggerMiddleware
	MetricsMiddleware      *metrics_mw.MetricsMiddleware
}

func (self *Controller) RequestContext(ctx context.Context) *RequestContext {
	return RequestContextFromContext(ctx)
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

func NewControllerOpts() *ControllerOpts {
	return &ControllerOpts{
		BaseAPIURL:      "http://localhost/",
		ConsumesContent: []string{"application/json"},
		ProducesContent: []string{"application/json"},
		Logger:          logger.NewDefaultLogger(os.Stdout, ""),
	}
}

func NewController(opts *ControllerOpts) *Controller {
	c := &Controller{
		options: opts,
	}
	c.errorFormatter = ErrorFormatterFn(c.formatErrors)
	return c
}

func RequestContextFromContext(ctx context.Context) *RequestContext {
	controller_rctx, _ := ctx.Value(requestContextCtxKey).(*RequestContext)
	return controller_rctx
}
