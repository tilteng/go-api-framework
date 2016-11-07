package controller

type Config struct {
	BaseRouter         *Router
	Logger             Logger
	NewRouteNotifierFn NewRouteNotifier
}
