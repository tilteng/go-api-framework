package panichandler_mw

import (
	"context"

	"github.com/tilteng/go-api-router/api_router"
)

type PanicHandler func(context.Context, interface{})

func (self PanicHandler) Panic(ctx context.Context, i interface{}) {
	self(ctx, i)
}

type PanicHandlerMiddleware struct {
	panicHandler PanicHandler
}

func (self *PanicHandlerMiddleware) NewWrapper() *PanicHandlerWrapper {
	return &PanicHandlerWrapper{
		panicHandler: self.panicHandler,
	}
}

type PanicHandlerWrapper struct {
	panicHandler PanicHandler
}

func (self *PanicHandlerWrapper) SetPanicHandler(panic_handler PanicHandler) *PanicHandlerWrapper {
	self.panicHandler = panic_handler
	return self
}

func (self *PanicHandlerWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		if ph := self.panicHandler; ph != nil {
			defer func() {
				if r := recover(); r != nil {
					ph.Panic(ctx, r)
				}
			}()
		}
		next(ctx)
	}
}

func NewMiddleware(panic_handler PanicHandler) *PanicHandlerMiddleware {
	return &PanicHandlerMiddleware{
		panicHandler: panic_handler,
	}
}
