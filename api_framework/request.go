package api_framework

import (
	"context"

	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
	"github.com/tilteng/go-errors/errors"
)

type contextKey struct {
	name string
}

func (self *contextKey) String() string {
	return "api_framework context value " + self.name
}

type RequestContext struct {
	*api_router.RequestContext
	appContext
	controller               *Controller
	serializerRequestContext serializers_mw.RequestContext
}

var requestContextCtxKey = &contextKey{"request_context"}

func (self *Controller) newRequestContextFromContext(ctx context.Context) *RequestContext {
	rctx := self.Router.RequestContext(ctx)
	ser_rctx := serializers_mw.RequestContextFromContext(ctx)
	return &RequestContext{
		RequestContext:           rctx,
		appContext:               self.appContext,
		controller:               self,
		serializerRequestContext: ser_rctx,
	}
}

func (self *RequestContext) ReadBody(ctx context.Context, v interface{}) error {
	return serializers_mw.DeserializedBody(ctx, v)
}

func (self *RequestContext) WriteResponse(ctx context.Context, v interface{}) error {
	if tilterr, ok := v.(errors.ErrorType); ok {
		status := tilterr.GetStatus()
		self.SetStatus(status)

		if status >= 500 {
			json, json_err := tilterr.AsJSON()
			if json_err != nil {
				self.LogErrorf("Returning exception: %+v", tilterr)
			} else {
				self.LogError("Returning exception: " + json)
			}
		}

		resp := self.controller.errorFormatter.FormatErrors(ctx, tilterr)

		return self.serializerRequestContext.WriteSerializedResponse(ctx, resp)
	}

	return self.serializerRequestContext.WriteSerializedResponse(ctx, v)
}
