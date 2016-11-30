package api_router

import "net/http"

type Route struct {
	fwRoute       FrameworkRoute
	router        *Router
	method        string
	path          string
	fullPath      string
	defaultStatus int
	routeFn       RouteFn
	virtual       bool
}

func (self *Route) RouteFn() RouteFn {
	return self.routeFn
}

func (self *Route) FullPath() string {
	return self.fullPath
}

func (self *Route) Method() string {
	return self.method
}

func (self *Route) Path() string {
	return self.path
}

func (self *Route) RouteVars(r *http.Request) map[string]string {
	if self.fwRoute == nil {
		return make(map[string]string)
	}
	return self.fwRoute.RouteVars(r)
}

func (self *Route) SetRouteFn(fn RouteFn) *Route {
	self.routeFn = fn
	return self
}

func (self *Route) SetDefaultStatus(status int) *Route {
	self.defaultStatus = status
	return self
}

func (self *Route) register() *Route {
	if self.defaultStatus == 0 {
		if self.method == "POST" {
			self.defaultStatus = 201
		} else {
			self.defaultStatus = 200
		}
	}

	self.fwRoute = self.router.fwRouter.NewRoute(
		self.method,
		self.path,
		self.handleRequest,
	)

	return self
}

func (self *Route) handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := NewContextForRequest(
		newResponseWriter(w, self.defaultStatus),
		r,
		self,
	)
	self.routeFn(ctx)
	// Ensure we've set status, even if no body was written
	ctx.writer.WriteStatusHeader()
}
