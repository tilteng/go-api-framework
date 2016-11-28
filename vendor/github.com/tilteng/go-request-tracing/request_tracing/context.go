package request_tracing

import "context"

type contextKey struct {
	name string
}

func (self *contextKey) String() string {
	return "request_tracing context value " + self.name
}

var requestTraceCtxKey = &contextKey{"requestTrace"}

func (self *requestTraceManager) RequestTraceFromContext(ctx context.Context) RequestTrace {
	rt, _ := ctx.Value(self.requestTraceCtxKey).(RequestTrace)
	return rt
}

func (self *requestTraceManager) ContextWithRequestTrace(ctx context.Context, rt RequestTrace) context.Context {
	return context.WithValue(ctx, self.requestTraceCtxKey, rt)
}
