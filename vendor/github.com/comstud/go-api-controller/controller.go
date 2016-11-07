package controller

import "context"

type Controller struct {
	*Router
	Logger
	Config *Config
}

func (self *Controller) SetLogger(logger Logger) *Controller {
	self.Logger = logger
	return self
}

func (self *Controller) RequestContext(ctx context.Context) *RequestContext {
	return RequestContextFromContext(ctx)
}

func NewController(config *Config) (*Controller, error) {
	bc := &Controller{
		Config: config,
		Router: config.BaseRouter,
		Logger: config.Logger,
	}

	if rn := config.NewRouteNotifierFn; rn != nil {
		bc.Router.SetNewRouteNotifier(rn)
	}

	return bc, nil
}
