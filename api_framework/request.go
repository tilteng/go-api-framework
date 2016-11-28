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

type privateContext interface {
	context.Context
}

type requestContext struct {
	*api_router.RequestContext
}

type RequestContext struct {
	privateContext
	requestContext
	appContext
	serializerRequestContext serializers_mw.RequestContext
}

var requestContextCtxKey = &contextKey{"request_context"}

func (self *RequestContext) Value(key interface{}) interface{} {
	if key == requestContextCtxKey {
		return self
	}
	return self.privateContext.Value(key)
}

func (self *Controller) newRequestContextFromContext(ctx context.Context) *RequestContext {
	router_rctx := self.Router.RequestContext(ctx)
	ser_rctx := serializers_mw.RequestContextFromContext(ctx)
	rctx := &RequestContext{
		privateContext:           ctx,
		appContext:               self.appContext,
		serializerRequestContext: ser_rctx,
	}
	rctx.requestContext.RequestContext = router_rctx
	return rctx
}

func (self *Controller) RequestContext(ctx context.Context) *RequestContext {
	if rctx, ok := ctx.Value(requestContextCtxKey).(*RequestContext); ok {
		return rctx
	}

	return self.newRequestContextFromContext(ctx)
}
