package api_router

import (
	"net/http"

	"github.com/gorilla/mux"
)

var _muxFramework = &muxFramework{}

type muxFramework struct{}

func (self *muxFramework) NewRouter() FrameworkRouter {
	return &muxRouter{Router: mux.NewRouter()}
}

type muxRouter struct {
	*mux.Router
}

func (self *muxRouter) NewRoute(method string, path string, fn http.HandlerFunc) FrameworkRoute {
	return &muxRoute{Route: self.HandleFunc(path, fn).Methods(method)}
}

func (self *muxRouter) SubRouterForPath(path string) FrameworkRouter {
	return &muxRouter{
		Router: self.PathPrefix(path).Subrouter(),
	}
}

type muxRoute struct {
	*mux.Route
}

func (self *muxRoute) RouteVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func MuxFramework() Framework {
	return _muxFramework
}

func NewMuxRouter() *Router {
	return NewRouter(_muxFramework)
}
