package api_framework

import (
	"context"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-errors/errors"
)

var ErrInternalServerError = errors.ErrInternalServerError
var ErrJSONSchemaValidationFailed = errors.ErrJSONSchemaValidationFailed

// Default ErrorFormatter
func (self *Controller) errorFormatter(errtype errors.ErrorType) interface{} {
	return errtype.AsJSONAPIResponse()
}

// Default PanicHandler.
func (self *Controller) handlePanic(ctx context.Context, v interface{}) {
	err_obj, ok := v.(*errors.Error)
	if !ok {
		err_obj = ErrInternalServerError.New()
		err_obj.SetInternal(v)
	}
	self.WriteResponse(ctx, err_obj)
}

// Default SerializerErrorHandler
func (self *Controller) handleSerializerError(ctx context.Context, err error) {
	rctx := api_router.RequestContextFromContext(ctx)
	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Default JSONSchemaErrorHandler
func (self *Controller) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	json_errors := result.Errors()
	api_errors := make(errors.Errors, 0, len(json_errors))
	for _, json_err := range json_errors {
		err := ErrJSONSchemaValidationFailed.New()
		err.Details = json_err.String()
		api_errors.AddError(err)
	}

	self.WriteResponse(ctx, api_errors)

	return false
}
