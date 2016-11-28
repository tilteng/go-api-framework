package serializers_mw

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tilteng/go-api-router/api_router"

	"bitbucket.org/ww/goautoneg"
)

type contextKey struct {
	name string
}

func (self *contextKey) String() string {
	return "serializers_mw context value " + self.name
}

type ErrorHandler func(context.Context, error)

func (self ErrorHandler) Error(ctx context.Context, err error) {
	self(ctx, err)
}

type SerializerMiddleware struct {
	consumes     []string
	produces     []string
	errorHandler ErrorHandler
}

func (self *SerializerMiddleware) NewWrapper() *SerializerWrapper {
	return &SerializerWrapper{
		consumes:     self.consumes,
		produces:     self.produces,
		errorHandler: self.errorHandler,
	}
}

type SerializerWrapper struct {
	consumes     []string
	produces     []string
	errorHandler ErrorHandler
}

var requestContextCtxKey = &contextKey{"requestContext"}

func (self *SerializerWrapper) newRequestContext(ctx context.Context, rctx *api_router.RequestContext) (*requestContext, error) {
	ctype := strings.Split(rctx.Header("Content-Type"), ";")[0]
	ctype = strings.Trim(ctype, " ")

	if ctype == "" {
		ctype = self.consumes[0]
	} else {
		var matched bool
		for _, consume_type := range self.consumes {
			if ctype == consume_type {
				matched = true
			}
			break
		}
		if !matched {
			return nil, fmt.Errorf(
				"No support for Content-Type media type: %s",
				ctype,
			)
		}
	}

	accept := rctx.Header("Accept")
	var atype string

	if accept == "" {
		atype = self.produces[0]
	} else {
		atype = goautoneg.Negotiate(accept, self.produces)
		if atype == "" {
			return nil, fmt.Errorf(
				"No support for Accept media type(s): %s",
				atype,
			)
		}
	}

	return &requestContext{
		rctx:         rctx,
		deserializer: serializers[ctype],
		serializer:   serializers[atype],
	}, nil
}

func (self *SerializerWrapper) SetErrorHandler(error_handler ErrorHandler) *SerializerWrapper {
	self.errorHandler = error_handler
	return self
}

func (self *SerializerWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		rctx := api_router.RequestContextFromContext(ctx)
		mw_rctx, err := self.newRequestContext(ctx, rctx)
		if err != nil {
			if self.errorHandler == nil {
				panic(err)
			}
			self.errorHandler.Error(ctx, err)
			return
		}
		ctx = context.WithValue(ctx, requestContextCtxKey, mw_rctx)
		rctx = rctx.WithContext(ctx)
		next(ctx)
	}
}

func NewMiddleware(consumes, produces []string, error_handler ErrorHandler) *SerializerMiddleware {
	return &SerializerMiddleware{
		consumes:     consumes,
		produces:     produces,
		errorHandler: error_handler,
	}
}

func RequestContextFromContext(ctx context.Context) RequestContext {
	mw_rctx, _ := ctx.Value(requestContextCtxKey).(*requestContext)
	return mw_rctx
}

func WriteSerializedResponse(ctx context.Context, v interface{}) error {
	mw_rctx, ok := ctx.Value(requestContextCtxKey).(*requestContext)
	if !ok {
		return errors.New("Request did not pass through serializers middleware")
	}
	return mw_rctx.WriteSerializedResponse(ctx, v)
}

func ReadDeserializedBody(ctx context.Context, v interface{}) error {
	mw_rctx, ok := ctx.Value(requestContextCtxKey).(*requestContext)
	if !ok {
		return errors.New("Request did not pass through serializers middleware")
	}
	return mw_rctx.ReadDeserializedBody(ctx, v)
}
