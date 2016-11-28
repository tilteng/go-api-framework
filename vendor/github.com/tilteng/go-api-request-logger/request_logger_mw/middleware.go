package request_logger_mw

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-logger/logger"
)

type LogBodyFilterFn func(context.Context, []byte) []byte

func (self LogBodyFilterFn) FilterBody(ctx context.Context, bytes []byte) []byte {
	return self(ctx, bytes)
}

type LogBodyFilter interface {
	FilterBody(context.Context, []byte) []byte
}

type LogHeadersFilterFn func(context.Context, http.Header) http.Header

func (self LogHeadersFilterFn) FilterHeaders(ctx context.Context, hdrs http.Header) http.Header {
	return self(ctx, hdrs)
}

type LogHeadersFilter interface {
	FilterHeaders(context.Context, http.Header) http.Header
}

type RequestLoggerOpts struct {
	LogBodyFilter    LogBodyFilter
	LogHeadersFilter LogHeadersFilter
	Logger           logger.CtxLogger
	Disable          bool
}

type RequestLoggerMiddleware struct {
	opts *RequestLoggerOpts
}

func (self *RequestLoggerMiddleware) NewWrapper(ctx context.Context, opts ...interface{}) *RequestLoggerWrapper {
	var opt *RequestLoggerOpts

	for _, opt_map_i := range opts {
		var ok bool
		opt, ok = opt_map_i.(*RequestLoggerOpts)
		if ok {
			break
		}
	}

	if opt == nil {
		opt = &RequestLoggerOpts{
			Logger: self.opts.Logger,
		}
	}

	if opt.Disable {
		return nil
	}

	return &RequestLoggerWrapper{
		base_opts: self.opts,
		opts:      opt,
	}
}

type RequestLoggerWrapper struct {
	opts      *RequestLoggerOpts
	base_opts *RequestLoggerOpts
}

func (self *RequestLoggerMiddleware) SetLogger(logger logger.CtxLogger) *RequestLoggerMiddleware {
	self.opts.Logger = logger
	return self
}

func (self *RequestLoggerWrapper) formatBody(ctx context.Context, body []byte) interface{} {
	if self.opts.LogBodyFilter != nil {
		body = self.opts.LogBodyFilter.FilterBody(ctx, body)
	}
	if self.base_opts.LogBodyFilter != nil {
		body = self.opts.LogBodyFilter.FilterBody(ctx, body)
	}
	// This is all hacky... just to get something logged for now
	// in a reasonable format.
	var body_obj interface{}

	if err := json.Unmarshal(body, &body_obj); err == nil {
		return body
	}

	body_str := string(body)

	if json_body, err := json.Marshal(body_str); err == nil {
		return json_body
	}

	return `"` + body_str + `"`
}

func (self *RequestLoggerWrapper) formatHeaders(ctx context.Context, hdrs http.Header) interface{} {
	if self.opts.LogHeadersFilter != nil {
		hdrs = self.opts.LogHeadersFilter.FilterHeaders(ctx, hdrs)
	}

	if self.base_opts.LogHeadersFilter != nil {
		hdrs = self.opts.LogHeadersFilter.FilterHeaders(ctx, hdrs)
	}

	json_hdrs, _ := json.Marshal(hdrs)
	return json_hdrs
}

func (self *RequestLoggerWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		rctx := api_router.RequestContextFromContext(ctx)
		rt := rctx.CurrentRoute()
		body, err := rctx.BodyCopy()
		if err != nil {
			panic(fmt.Sprintf("Couldn't read body: %+v", err))
		}
		if self.opts.Logger != nil && !self.opts.Disable {
			self.opts.Logger.LogDebugf(
				ctx,
				`Received request: route:{"method":"%s","route":"%s","path":"%s"} headers:%s body:%s`,
				rt.Method(),
				rt.FullPath(),
				rctx.HTTPRequest().URL.String(),
				self.formatHeaders(ctx, rctx.HTTPRequest().Header),
				self.formatBody(ctx, body),
			)
		}

		// Do something with body
		next(ctx)

		writer := rctx.ResponseWriter()
		body = writer.ResponseCopy()

		self.opts.Logger.LogDebugf(
			ctx,
			`Sent response: route:{"status":%d,"method":"%s","route":"%s","path":"%s"} headers:%s body:%s`,
			writer.Status(),
			rt.Method(),
			rt.FullPath(),
			rctx.HTTPRequest().URL.String(),
			self.formatHeaders(ctx, rctx.ResponseWriter().Header()),
			self.formatBody(ctx, body),
		)
	}
}

func NewMiddleware(opts *RequestLoggerOpts) *RequestLoggerMiddleware {
	if opts == nil {
		opts = &RequestLoggerOpts{}
	}
	if opts.Logger == nil {
		opts.Logger = logger.DefaultStdoutCtxLogger()
	}
	return &RequestLoggerMiddleware{
		opts: opts,
	}
}
