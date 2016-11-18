package apache_logger_mw

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/tilteng/go-api-router/api_router"
)

type ApacheLoggerMiddleware struct {
	defaultCombined bool
	defaultWriter   io.Writer
}

func (self *ApacheLoggerMiddleware) NewWrapper() *ApacheLogWrapper {
	return &ApacheLogWrapper{
		combined: self.defaultCombined,
		writer:   self.defaultWriter,
	}
}

type ApacheLogWrapper struct {
	combined bool
	writer   io.Writer
}

func (self *ApacheLogWrapper) buildCommonLogLine(rctx *api_router.RequestContext, t time.Time, url *url.URL) []byte {
	http_req := rctx.HTTPRequest()
	resp_writer := rctx.ResponseWriter()
	return buildCommonLogLine(
		http_req,
		url,
		t,
		resp_writer.Status(),
		resp_writer.Size(),
	)
}

func (self *ApacheLogWrapper) writeCombinedLog(rctx *api_router.RequestContext, t time.Time, url *url.URL) {
	http_req := rctx.HTTPRequest()
	buf := self.buildCommonLogLine(rctx, t, url)
	// Stolen from github.com/gorilla/handlers/handlers.go
	buf = append(buf, ` "`...)
	buf = appendQuoted(buf, http_req.Referer())
	buf = append(buf, `" "`...)
	buf = appendQuoted(buf, http_req.UserAgent())
	buf = append(buf, '"', '\n')
	self.writer.Write(buf)
}

func (self *ApacheLogWrapper) writeCommonLog(rctx *api_router.RequestContext, t time.Time, url *url.URL) {
	buf := self.buildCommonLogLine(rctx, t, url)
	buf = append(buf, '\n')
	self.writer.Write(buf)
}

func (self *ApacheLogWrapper) DisableCombined() *ApacheLogWrapper {
	self.combined = false
	return self
}

func (self *ApacheLogWrapper) EnableCombined() *ApacheLogWrapper {
	self.combined = true
	return self
}

func (self *ApacheLogWrapper) SetWriter(writer io.Writer) *ApacheLogWrapper {
	self.writer = writer
	return self
}

func (self *ApacheLogWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		t := time.Now()
		rctx := api_router.RequestContextFromContext(ctx)
		url := *rctx.HTTPRequest().URL

		next(ctx)

		if self.combined {
			self.writeCombinedLog(rctx, t, &url)
		} else {
			self.writeCommonLog(rctx, t, &url)
		}
	}
}

func NewMiddleware(writer io.Writer, combined bool) *ApacheLoggerMiddleware {
	return &ApacheLoggerMiddleware{
		defaultCombined: combined,
		defaultWriter:   writer,
	}
}
