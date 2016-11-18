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
func (self *TiltController) errorFormatter(errs errors.Errors) interface{} {
	return errs.AsJSONAPIResponse()
}

// Default PanicHandler.
func (self *TiltController) handlePanic(ctx context.Context, v interface{}) {
	err_obj, ok := v.(*errors.Error)
	if !ok {
		err_obj = ErrInternalServerError.New(0)
		err_obj.SetInternal(v)
	}
	self.WriteResponse(ctx, err_obj)
}

// Default SerializerErrorHandler
func (self *TiltController) handleSerializerError(ctx context.Context, err error) {
	rctx := api_router.RequestContextFromContext(ctx)
	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Default JSONSchemaErrorHandler
func (self *TiltController) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	json_errors := result.Errors()

	api_errors := make(errors.Errors, len(json_errors), len(json_errors))
	for i, json_err := range json_errors {
		err := ErrJSONSchemaValidationFailed.New(0)
		err.Details = json_err.String()
		api_errors[i] = err
	}

	self.WriteResponse(ctx, api_errors)

	return false
}
