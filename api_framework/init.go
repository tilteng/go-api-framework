package api_framework

import (
	"context"
	"fmt"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-panichandler/panichandler_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-logger/apache_logger_mw"
	"github.com/tilteng/go-metrics/metrics_mw"
)

func (self *Controller) setupSchemaRoutes() error {
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

func (self *Controller) setupDefaultOptions() error {
	if self.options.BaseRouter == nil {
		self.options.BaseRouter = api_router.NewMuxRouter()
	}

	if mcli := self.options.MetricsClient; mcli != nil && self.MetricsMiddleware == nil {
		self.MetricsMiddleware = metrics_mw.NewMiddleware(mcli)
	}

	if self.PanicHandlerMiddleware == nil {
		self.PanicHandlerMiddleware = panichandler_mw.NewMiddleware(
			panichandler_mw.PanicHandlerFn(self.handlePanic),
		)
	}

	if self.JSONSchemaMiddleware == nil && self.options.JSONSchemaFilePath != "" {
		var route_prefix string
		if rp := self.options.JSONSchemaRoutePath; rp != "" {
			route_prefix = self.options.BaseAPIURL + rp
		}

		js_mw := jsonschema_mw.NewMiddlewareWithLinkPathPrefix(
			self.handleJSONSchemaError,
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
			self.handleSerializerError,
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

	return nil
}

// Wrap the original route. We want to achieve this order:
// metrics -> apache-logger -> serializer -> panic_handler -> jsonschema
// Ie, we want the logger to log exactly what is returned after
// serialization. We want the ability to serialize panic_handler
// responses. And json schema validation should just happen right
// before we call the real route handler.
func (self *Controller) wrapNewRoute(rt *api_router.Route, opts ...interface{}) {
	orig_fn := rt.RouteFn()

	// Create our Request right before calling middleware
	fn := func(ctx context.Context) {
		ctx, _ = self.getContexts(ctx)
		orig_fn(ctx)
	}

	// Wrap with our request last

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

	if metrics_mw := self.MetricsMiddleware; metrics_mw != nil {
		fn = metrics_mw.NewWrapper().Wrap(fn)
	}

	rt.SetRouteFn(fn)

	if logger := self.options.Logger; logger != nil {
		logger.Debugf("Registered route: %s %s", rt.Method(), rt.FullPath())
	}
}

// Final initialization that should be called before registering routes.
// If you wish to modify any defaults on Controller, you should do this
// before calling Init()
func (self *Controller) Init() error {
	if err := self.setupDefaultOptions(); err != nil {
		return err
	}

	self.Router = self.options.BaseRouter
	self.Logger = self.options.Logger
	self.Router.SetNewRouteNotifier(self.wrapNewRoute)

	if self.JSONSchemaMiddleware != nil && self.options.JSONSchemaRoutePath != "" {
		if err := self.setupSchemaRoutes(); err != nil {
			return err
		}
	}

	return nil
}
