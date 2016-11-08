package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/comstud/go-api-controller"
	"github.com/comstud/go-api-controller/middleware/jsonschema"
	"github.com/comstud/go-api-controller/middleware/panichandler"
	"github.com/comstud/go-api-controller/middleware/serializers"
)

type TiltControllerOpts struct {
	BaseAPIURL             string
	BaseRouter             *controller.Router
	ConsumesContent        []string
	ProducesContent        []string
	JSONSchemaRoutePath    string
	JSONSchemaFilePath     string
	JSONSchemaErrorHandler jsonschema.ErrorHandler
	PanicHandler           panichandler.PanicHandler
	SerializerErrorHandler serializers.ErrorHandler
	Logger                 controller.Logger
}

func NewTiltControllerOpts() *TiltControllerOpts {
	return &TiltControllerOpts{
		BaseAPIURL:             "http://localhost/",
		ConsumesContent:        []string{"application/json"},
		ProducesContent:        []string{"application/json"},
		PanicHandler:           DefaultPanicHandler,
		SerializerErrorHandler: DefaultSerializerErrorHandler,
		JSONSchemaErrorHandler: DefaultJSONSchemaErrorHandler,
		Logger:                 controller.NewDefaultLogger(os.Stdout, ""),
	}
}

type TiltController struct {
	*controller.Controller
	options                *TiltControllerOpts
	JSONSchemaMiddleware   *jsonschema.JSONSchemaMiddleware
	PanicHandlerMiddleware *panichandler.PanicHandlerMiddleware
	SerializerMiddleware   *serializers.SerializerMiddleware
}

func (self *TiltController) ReadBody(ctx context.Context, v interface{}) error {
	return serializers.DeserializedBody(ctx, v)
}

func (self *TiltController) WriteResponse(ctx context.Context, v interface{}) error {
	if err, ok := v.(ErrorType); ok {
		HandleErrorType(ctx, err)
		return serializers.WriteSerializedResponse(ctx, &ErrorResponse{v})
	}
	return serializers.WriteSerializedResponse(ctx, v)

}

func (self *TiltController) WriteError(ctx context.Context, err *Error) error {
	HandleErrorType(ctx, &err.BaseError)
	return serializers.WriteSerializedResponse(ctx, &ErrorResponse{err})
}

func (self *TiltController) NewError(status int, code string, err_str string) *Error {
	return NewError(status, code, err_str)
}

func (self *TiltController) setupSchemaRoutes() error {
	schemas := self.JSONSchemaMiddleware.GetSchemas()

	self.GET(self.options.JSONSchemaRoutePath, func(ctx context.Context) {
		rctx := controller.RequestContextFromContext(ctx)

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

	for k, v := range schemas {
		sr.GET("/"+k, func(ctx context.Context) {
			rctx := controller.RequestContextFromContext(ctx)
			rctx.SetStatus(200)
			rctx.SetResponseHeader("Content-Type", "application/json+schema")
			rctx.WriteResponseString(v.GetJSONString())
		})
	}

	return nil
}

func (self *TiltController) Init() error {
	if self.PanicHandlerMiddleware == nil && self.options.PanicHandler != nil {
		self.PanicHandlerMiddleware = panichandler.NewMiddleware(
			self.options.PanicHandler,
		)
	}

	if self.JSONSchemaMiddleware == nil && self.options.JSONSchemaFilePath != "" {
		var route_prefix string
		if rp := self.options.JSONSchemaRoutePath; rp != "" {
			route_prefix = self.options.BaseAPIURL + rp
		}

		js_mw := jsonschema.NewMiddlewareWithLinkPathPrefix(
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
		self.SerializerMiddleware = serializers.NewMiddleware(
			self.options.ConsumesContent,
			self.options.ProducesContent,
			self.options.SerializerErrorHandler,
		)
	}

	new_route_notifier := func(rt *controller.Route, opts ...interface{}) {
		fn := rt.ControllerFn()

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

		rt.SetControllerFn(fn)

		if logger := self.options.Logger; logger != nil {
			logger.Debugf("Registered route: %s %s", rt.Method(), rt.FullPath())
		}
	}

	if self.options.BaseRouter == nil {
		self.options.BaseRouter = controller.NewMuxRouter()
	}

	controller_config := &controller.Config{
		BaseRouter:         self.options.BaseRouter,
		Logger:             self.options.Logger,
		NewRouteNotifierFn: new_route_notifier,
	}

	base_controller, err := controller.NewController(controller_config)
	if err != nil {
		return err
	}

	self.Controller = base_controller

	if self.JSONSchemaMiddleware != nil && self.options.JSONSchemaRoutePath != "" {
		err = self.setupSchemaRoutes()
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *TiltController) GenUUID() UUID {
	return GenUUID()
}

func (self *TiltController) GenUUIDHex() string {
	return GenUUIDHex()
}

func NewTiltController(opts *TiltControllerOpts) *TiltController {
	return &TiltController{options: opts}
}
