package api_framework

import (
	"context"

	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
)

type contextKey struct {
	name string
}

func (self *contextKey) String() string {
	return "api_framework context value " + self.name
}

type RequestContext struct {
	context.Context
	*api_router.RequestContext
	appContext
	controller               *Controller
	serializerRequestContext serializers_mw.RequestContext
}

var requestContextCtxKey = &contextKey{"request_context"}

func (self *RequestContext) Value(key interface{}) interface{} {
	if key == requestContextCtxKey {
		return self
	}
	return self.Context.Value(key)
}

func (self *Controller) newRequestContextFromContext(ctx context.Context) *RequestContext {
	router_rctx := self.Router.RequestContext(ctx)
	ser_rctx := serializers_mw.RequestContextFromContext(ctx)
	return &RequestContext{
		Context:                  ctx,
		RequestContext:           router_rctx,
		appContext:               self.appContext,
		controller:               self,
		serializerRequestContext: ser_rctx,
	}
}

func (self *Controller) RequestContext(ctx context.Context) *RequestContext {
	if rctx, ok := ctx.Value(requestContextCtxKey).(*RequestContext); ok {
		return rctx
	}

	return self.newRequestContextFromContext(ctx)
}
