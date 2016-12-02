package metrics_mw

import (
	"context"
	"fmt"
	"time"

	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-metrics/metrics"
)

type MetricsMiddleware struct {
	client     metrics.MetricsClient
	timingName string
}

func (self *MetricsMiddleware) NewWrapper() *MetricsWrapper {
	return &MetricsWrapper{
		client:     self.client,
		timingName: self.timingName,
	}
}

func (self *MetricsMiddleware) SetTimingName(name string) {
	self.timingName = name
}

type MetricsWrapper struct {
	client     metrics.MetricsClient
	timingName string
}

func (self *MetricsWrapper) formatValue(s string) string {
	bytes := []byte(s)
	for i, byt := range bytes {
		if byt == '{' || byt == '}' {
			bytes[i] = ':'
		}
	}
	return string(bytes)
}

func (self *MetricsWrapper) sendMetrics(ctx context.Context, duration time.Duration) {
	// Catch and ignore any errors
	defer func() {
		recover()
	}()

	rctx := api_router.RequestContextFromContext(ctx)
	writer := rctx.ResponseWriter()
	http_req := rctx.HTTPRequest()
	cur_route := rctx.CurrentRoute()

	self.client.Timing(
		self.timingName,
		duration,
		1,
		map[string]string{
			"route":  self.formatValue(cur_route.Path()),
			"path":   http_req.URL.String(),
			"method": http_req.Method,
			"status": fmt.Sprintf("%d", writer.Status()),
			"size":   fmt.Sprintf("%d", writer.Size()),
		},
	)
}

func (self *MetricsWrapper) SetTimingName(name string) {
	self.timingName = name
}

func (self *MetricsWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		next(ctx)
		// Grab duration before the go routine
		rctx := api_router.RequestContextFromContext(ctx)
		duration := rctx.TimeElapsed()
		go self.sendMetrics(ctx, duration)
	}
}

func NewMiddleware(client metrics.MetricsClient) *MetricsMiddleware {
	return &MetricsMiddleware{
		client:     client,
		timingName: "route.timing",
	}
}
