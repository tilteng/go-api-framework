package api_router

import "net/http"

type FrameworkRoute interface {
	RouteVars(*http.Request) map[string]string
}

type FrameworkRouter interface {
	http.Handler
	NewRoute(method string, path string, fn http.HandlerFunc) FrameworkRoute
	SubRouterForPath(path string) FrameworkRouter
	Set404Handler(fn http.HandlerFunc)
}

type Framework interface {
	NewRouter() FrameworkRouter
}
