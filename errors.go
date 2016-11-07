package controller

import (
	"context"
	"fmt"

	"github.com/comstud/go-api-controller"
	"github.com/comstud/go-api-controller/middleware/jsonschema"
	"github.com/comstud/go-api-controller/middleware/serializers"
)

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type BaseError struct {
	Id              string      `json:"id"`
	StatusCode      int         `json:"status"`
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	InternalError   string      `json:"-"`
	InternalDetails interface{} `json:"-"`
}

type Error struct {
	BaseError
	Code    string      `json:"code"`
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func NewError(status int, code string, err_str string) *Error {
	return &Error{
		BaseError: BaseError{
			Id:         "ERR" + GenUUIDHex(),
			StatusCode: status,
		},
		Code:  code,
		Error: err_str,
	}
}

func DefaultPanicHandler(ctx context.Context, v interface{}) {
	err_cont_obj, ok := v.(*ErrorResponse)
	if !ok {
		if _, ok := v.(BaseError); ok {
			err_cont_obj = &ErrorResponse{v}
		} else {
			err_obj := NewError(
				500,
				"ERR_ID_INTERNAL_SERVER_ERROR",
				"An unhandled exception has occurred",
			)

			err_obj.InternalDetails = v

			if s, ok := v.(fmt.Stringer); ok {
				err_obj.InternalError = s.String()
			} else if e, ok := v.(error); ok {
				err_obj.InternalError = e.Error()
			}

			err_cont_obj = &ErrorResponse{err_obj}
		}
	}

	if base_err, ok := err_cont_obj.Error.(BaseError); ok {
		HandleBaseError(ctx, &base_err)
	}

	// TODO: Do something with error?
	if serializers.WriteSerializedResponse(ctx, err_cont_obj) == nil {
		return
	}

	rctx := controller.RequestContextFromContext(ctx)
	rctx.WriteResponseString("An unknown error occurred")
}

func HandleBaseError(ctx context.Context, base_err *BaseError) {
	rctx := controller.RequestContextFromContext(ctx)

	if base_err.StatusCode <= 0 {
		base_err.StatusCode = 500
	}

	if base_err.Method == "" {
		base_err.Method = rctx.CurrentRoute().Method()
	}

	if base_err.Path == "" {
		base_err.Path = rctx.CurrentRoute().Path()
	}

	rctx.SetStatus(base_err.StatusCode)
}

func DefaultSerializerErrorHandler(ctx context.Context, err error) {
	rctx := controller.RequestContextFromContext(ctx)
	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

func DefaultJSONSchemaErrorHandler(ctx context.Context, result *jsonschema.JSONSchemaResult) bool {
	rctx := controller.RequestContextFromContext(ctx)
	rctx.SetStatus(400)
	details := []string{}
	for _, err := range result.Errors() {
		details = append(details, err.Description())
	}
	err := NewError(
		400,
		"ERR_ID_BAD_DATA",
		"Validation for json schema failed",
	)
	err.Details = details
	panic(err)
}
