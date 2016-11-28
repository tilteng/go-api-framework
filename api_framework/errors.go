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
	rctx := self.RequestContext(ctx)
	if self.options.ErrorFormatter != nil {
		return self.options.ErrorFormatter.FormatErrors(rctx, errtype)
	}
	return errtype.AsJSONAPIResponse()
}

// Called when a panic occurs. Pass to custom callback, if set.
func (self *Controller) handlePanic(ctx context.Context, v interface{}) {
	rctx := self.RequestContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			self.logger.LogErrorf(
				rctx,
				"Received panic while processing another panic: %+v",
				r,
			)
			rctx.SetStatus(500)
			self.WriteResponse(rctx, nil)
		}
	}()

	if self.options.PanicHandler != nil {
		self.options.PanicHandler.Panic(rctx, v)
		return
	}

	err_obj, ok := v.(*errors.Error)
	if !ok {
		err_obj = ErrInternalServerError.New("")
		err_obj.SetInternal(v)
	}

	self.WriteResponse(rctx, err_obj)
}

// Called when a serializer erorr occurs. Pass to custom callback, if set.
func (self *Controller) handleSerializerError(ctx context.Context, err error) {
	rctx := self.RequestContext(ctx)

	if self.options.SerializerErrorHandler != nil {
		self.options.SerializerErrorHandler.Error(rctx, err)
		return
	}

	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Called when a json schema validation failure or error occurs.
// Pass to custom callback, if set.
func (self *Controller) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	rctx := self.RequestContext(ctx)

	if self.options.JSONSchemaErrorHandler != nil {
		return self.options.JSONSchemaErrorHandler.Error(rctx, result)
	}

	json_errors := result.Errors()
	api_errors := make(errors.Errors, 0, len(json_errors))
	for _, json_err := range json_errors {
		err := ErrJSONSchemaValidationFailed.New("")
		err.Details = json_err.String()
		api_errors.AddError(err)
	}

	self.WriteResponse(rctx, api_errors)

	return false
}
