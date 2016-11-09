package api_framework

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-panichandler/panichandler_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-logger/apache_logger_mw"
	"github.com/tilteng/go-logger/logger"
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
	Logger                 logger.Logger
}

func NewTiltControllerOpts() *TiltControllerOpts {
	return &TiltControllerOpts{
		BaseAPIURL:             "http://localhost/",
		ConsumesContent:        []string{"application/json"},
		ProducesContent:        []string{"application/json"},
		PanicHandler:           DefaultPanicHandler,
		SerializerErrorHandler: DefaultSerializerErrorHandler,
		JSONSchemaErrorHandler: DefaultJSONSchemaErrorHandler,
		Logger:                 logger.NewDefaultLogger(os.Stdout, ""),
	}
}

type TiltController struct {
	*api_router.Router
	logger.Logger
	options                *TiltControllerOpts
	JSONSchemaMiddleware   *jsonschema_mw.JSONSchemaMiddleware
	PanicHandlerMiddleware *panichandler_mw.PanicHandlerMiddleware
	SerializerMiddleware   *serializers_mw.SerializerMiddleware
	ApacheLoggerMiddleware *apache_logger_mw.ApacheLoggerMiddleware
}

func (self *TiltController) ReadBody(ctx context.Context, v interface{}) error {
	return serializers_mw.DeserializedBody(ctx, v)
}

func (self *TiltController) WriteResponse(ctx context.Context, v interface{}) error {
	if err, ok := v.(ErrorType); ok {
		HandleErrorType(ctx, err)
		return serializers_mw.WriteSerializedResponse(ctx, &ErrorResponse{v})
	}
	return serializers_mw.WriteSerializedResponse(ctx, v)

}

func (self *TiltController) WriteError(ctx context.Context, err *Error) error {
	HandleErrorType(ctx, &err.BaseError)
	return serializers_mw.WriteSerializedResponse(ctx, &ErrorResponse{err})
}

func (self *TiltController) NewError(status int, code string, err_str string) *Error {
	return NewError(status, code, err_str)
}

func (self *TiltController) setupSchemaRoutes() error {
	if self.JSONSchemaMiddleware == nil || self.options.JSONSchemaRoutePath == "" {
		panic("setupSchemaRoutes() called with no middleware or route path")
	}

	schemas := self.JSONSchemaMiddleware.GetSchemas()

	self.GET(self.options.JSONSchemaRoutePath, func(ctx context.Context) {
		rctx := self.RequestContext(ctx)

		rctx.SetStatus(200)
		rctx.SetResponseHeader("Content-Type", "application/json")
		rctx.WriteResponseString("[")

		prefix := ""

		for k, v := range schemas {
			rctx.WriteResponseString(prefix + fmt.Sprintf(
				`{ "%s": %s }`, k, v.GetJSONString(),
			))
			prefix = ", "
		}
		rctx.WriteResponseString("]")
	})

	sr := self.SubRouterForPath(self.options.JSONSchemaRoutePath)

	for k := range schemas {
		sr.GET("/"+k, func(ctx context.Context) {
			// We need the closure to have its own value, so we don't
			// grab the key from the 'range' above.
			v := schemas[k]
			rctx := self.RequestContext(ctx)
			rctx.SetStatus(200)
			rctx.SetResponseHeader("Content-Type", "application/json+schema")
			rctx.WriteResponseString(v.GetJSONString())
		})
	}

	return nil
}

func (self *TiltController) Init() error {
	if self.PanicHandlerMiddleware == nil && self.options.PanicHandler != nil {
		self.PanicHandlerMiddleware = panichandler_mw.NewMiddleware(
			self.options.PanicHandler,
		)
	}

	if self.JSONSchemaMiddleware == nil && self.options.JSONSchemaFilePath != "" {
		var route_prefix string
		if rp := self.options.JSONSchemaRoutePath; rp != "" {
			route_prefix = self.options.BaseAPIURL + rp
		}

		js_mw := jsonschema_mw.NewMiddlewareWithLinkPathPrefix(
			self.options.JSONSchemaErrorHandler,
			route_prefix,
		).SetLogger(self.options.Logger)

		err := js_mw.LoadFromPath(self.options.JSONSchemaFilePath)
		if err != nil {
			return err
		}

		self.JSONSchemaMiddleware = js_mw
	}

	if self.SerializerMiddleware == nil {
		self.SerializerMiddleware = serializers_mw.NewMiddleware(
			self.options.ConsumesContent,
			self.options.ProducesContent,
			self.options.SerializerErrorHandler,
		)
	}

	if self.ApacheLoggerMiddleware == nil {
		if self.options.ApacheLogWriter != nil {
			self.ApacheLoggerMiddleware = apache_logger_mw.NewMiddleware(
				self.options.ApacheLogWriter,
				self.options.ApacheLogCombined,
			)
		}
	}

	// Callback for new route being added.
	new_route_notifier := func(rt *api_router.Route, opts ...interface{}) {
		// Wrap the original route. We want to achieve this order:
		// apache-logger -> serializer -> panic_handler -> jsonschema
		// Ie, we want the logger to log exactly what is returned after
		// serialization. We want the ability to serialize panic_handler
		// responses. And json schema validation should just happen right
		// before we call the real route handler.

		fn := rt.RouteFn()

		if js_mw := self.JSONSchemaMiddleware; js_mw != nil {
			if wrapper := js_mw.NewWrapperFromRouteOptions(opts...); wrapper != nil {
				fn = wrapper.Wrap(fn)
			}
		}

		if ph_mw := self.PanicHandlerMiddleware; ph_mw != nil {
			fn = ph_mw.NewWrapper().Wrap(fn)
		}

		if ser_mw := self.SerializerMiddleware; ser_mw != nil {
			fn = ser_mw.NewWrapper().Wrap(fn)
		}

		if log_mw := self.ApacheLoggerMiddleware; log_mw != nil {
			fn = log_mw.NewWrapper().Wrap(fn)
		}

		rt.SetRouteFn(fn)

		if logger := self.options.Logger; logger != nil {
			logger.Debugf("Registered route: %s %s", rt.Method(), rt.FullPath())
		}
	}

	if self.options.BaseRouter == nil {
		self.options.BaseRouter = api_router.NewMuxRouter()
	}

	self.Router = self.options.BaseRouter
	self.Logger = self.options.Logger
	self.Router.SetNewRouteNotifier(new_route_notifier)

	if self.JSONSchemaMiddleware != nil && self.options.JSONSchemaRoutePath != "" {
		if err := self.setupSchemaRoutes(); err != nil {
			return err
		}
	}

	return nil
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

func NewTiltController(opts *TiltControllerOpts) *TiltController {
	return &TiltController{options: opts}
}
