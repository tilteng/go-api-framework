package api_framework

import (
	"context"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-errors/errors"
)

var ErrInternalServerError = errors.ErrInternalServerError
var ErrJSONSchemaValidationFailed = errors.ErrJSONSchemaValidationFailed

// Called to format an error or errors. Pass to custom callback, if set.
func (self *Controller) formatErrors(ctx context.Context, errtype errors.ErrorType) interface{} {
	ctx, _ = self.getContexts(ctx)
	if self.options.ErrorFormatter != nil {
		return self.options.ErrorFormatter.FormatErrors(ctx, errtype)
	}
	return errtype.AsJSONAPIResponse()
}

// Called when a panic occurs. Pass to custom callback, if set.
func (self *Controller) handlePanic(ctx context.Context, v interface{}) {
	var rctx *RequestContext
	ctx, rctx = self.getContexts(ctx)

	if self.options.PanicHandler != nil {
		self.options.PanicHandler.Panic(ctx, v)
		return
	}

	err_obj, ok := v.(*errors.Error)
	if !ok {
		err_obj = ErrInternalServerError.New()
		err_obj.SetInternal(v)
	}

	rctx.WriteResponse(ctx, err_obj)
}

// Called when a serializer erorr occurs. Pass to custom callback, if set.
func (self *Controller) handleSerializerError(ctx context.Context, err error) {
	var rctx *RequestContext
	ctx, rctx = self.getContexts(ctx)

	if self.options.SerializerErrorHandler != nil {
		self.options.SerializerErrorHandler.Error(ctx, err)
		return
	}

	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Called when a json schema validation failure or error occurs.
// Pass to custom callback, if set.
func (self *Controller) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	var rctx *RequestContext
	ctx, rctx = self.getContexts(ctx)

	if self.options.JSONSchemaErrorHandler != nil {
		return self.options.JSONSchemaErrorHandler.Error(ctx, result)
	}

	json_errors := result.Errors()
	api_errors := make(errors.Errors, 0, len(json_errors))
	for _, json_err := range json_errors {
		err := ErrJSONSchemaValidationFailed.New()
		err.Details = json_err.String()
		api_errors.AddError(err)
	}

	rctx.WriteResponse(ctx, api_errors)

	return false
}

// Ensure we use a context that has our Request in callbacks
func (self *Controller) getContexts(ctx context.Context) (context.Context, *RequestContext) {
	rctx := self.RequestContext(ctx)
	if rctx == nil {
		rctx = self.newRequestContextFromContext(ctx)
		ctx = context.WithValue(ctx, requestContextCtxKey, rctx)
	}
	return ctx, rctx
}
