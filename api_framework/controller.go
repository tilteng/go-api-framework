package api_framework

import (
	"context"
	"io"
	"os"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-panichandler/panichandler_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-logger/apache_logger_mw"
	"github.com/tilteng/go-logger/logger"
	"github.com/tilteng/go-metrics/metrics"
	"github.com/tilteng/go-metrics/metrics_mw"
)

type TiltControllerOpts struct {
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
	MetricsClient          metrics.MetricsClient
	Logger                 logger.Logger
}

type TiltController struct {
	*api_router.Router
	logger.Logger
	options                *TiltControllerOpts
	JSONSchemaMiddleware   *jsonschema_mw.JSONSchemaMiddleware
	PanicHandlerMiddleware *panichandler_mw.PanicHandlerMiddleware
	SerializerMiddleware   *serializers_mw.SerializerMiddleware
	ApacheLoggerMiddleware *apache_logger_mw.ApacheLoggerMiddleware
	MetricsMiddleware      *metrics_mw.MetricsMiddleware
}

func (self *TiltController) ReadBody(ctx context.Context, v interface{}) error {
	return serializers_mw.DeserializedBody(ctx, v)
}

func (self *TiltController) WriteResponse(ctx context.Context, v interface{}) error {
	if err, ok := v.(ErrorType); ok {
		return self.handleError(ctx, err)
	} else {
		return serializers_mw.WriteSerializedResponse(ctx, v)
	}

}

func (self *TiltController) WriteError(ctx context.Context, err *Error) error {
	return self.handleError(ctx, err)
}

func (self *TiltController) GenUUID() *UUID {
	return GenUUID()
}

func (self *TiltController) GenUUIDHex() string {
	return GenUUIDHex()
}

func (self *TiltController) UUIDFromString(s string) *UUID {
	return UUIDFromString(s)
}

func NewTiltControllerOpts() *TiltControllerOpts {
	return &TiltControllerOpts{
		BaseAPIURL:      "http://localhost/",
		ConsumesContent: []string{"application/json"},
		ProducesContent: []string{"application/json"},
		Logger:          logger.NewDefaultLogger(os.Stdout, ""),
	}
}

func NewTiltController(opts *TiltControllerOpts) *TiltController {
	return &TiltController{options: opts}
}
