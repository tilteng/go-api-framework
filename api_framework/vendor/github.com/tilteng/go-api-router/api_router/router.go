package api_router

import (
	"context"
	"net/http"
)

type RouteFn func(context.Context)
type NewRouteNotifier func(*Route, ...interface{})

type Router struct {
	basePath         string
	fwRouter         FrameworkRouter
	topRouter        *Router
	newRouteNotifier NewRouteNotifier
	routes           []*Route
}

func (self *Router) SetNewRouteNotifier(route_notifier NewRouteNotifier) *Router {
	self.newRouteNotifier = route_notifier
	return self
}

func (self *Router) NewRoute(method string, path string, fn RouteFn, opts ...interface{}) *Route {
	rt := &Route{
		router:   self,
		method:   method,
		path:     path,
		fullPath: combinePaths(self.basePath, path),
		routeFn:  fn,
	}
	rt.register(fn)
	self.topRouter.routes = append(self.topRouter.routes, rt)
	if self.newRouteNotifier != nil {
		self.newRouteNotifier(rt, opts...)
	}
	return rt
}

func (self *Router) SubRouterForPath(path string) *Router {
	return &Router{
		basePath:         combinePaths(self.basePath, path),
		fwRouter:         self.fwRouter.SubRouterForPath(path),
		topRouter:        self.topRouter,
		newRouteNotifier: self.newRouteNotifier,
	}
}

// Implements http.Handler interface
func (self *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.fwRouter.ServeHTTP(w, r)
}

func (self *Router) DELETE(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("DELETE", path, fn, opts...)
}

func (self *Router) GET(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("GET", path, fn, opts...)
}

func (self *Router) HEAD(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("HEAD", path, fn, opts...)
}

func (self *Router) PATCH(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("PATCH", path, fn, opts...)
}

func (self *Router) POST(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("POST", path, fn, opts...)
}

func (self *Router) PUT(path string, fn RouteFn, opts ...interface{}) *Route {
	return self.NewRoute("PUT", path, fn, opts...)
}

func (self *Router) RequestContext(ctx context.Context) *RequestContext {
	return RequestContextFromContext(ctx)
}

func NewRouter(framework Framework) *Router {
	r := &Router{
		basePath: "/",
		fwRouter: framework.NewRouter(),
		routes:   []*Route{},
	}
	r.topRouter = r
	return r
}

func combinePaths(basepath string, path string) string {
	fullpath := basepath
	basepath_len := len(basepath)
	if basepath_len == 0 || basepath[basepath_len-1:] != "/" {
		fullpath += "/"
	}
	if len(path) != 0 && path[:1] == "/" {
		fullpath += path[1:]
	} else {
		fullpath += path
	}
	return fullpath
}
