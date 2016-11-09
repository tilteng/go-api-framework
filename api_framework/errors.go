package api_framework

import (
	"context"
	"fmt"
	"log"

	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-api-router/api_router"
	"github.com/tilteng/go-api-serializers/serializers_mw"
)

type ErrorType interface {
	GetID() string
	GetStatusCode() int
	SetStatusCode(int)
	GetMethod() string
	SetMethod(string)
	GetPath() string
	SetPath(string)
	GetInternalError() string
	GetInternalDetails() interface{}
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type BaseError struct {
	ID              string      `json:"id"`
	StatusCode      int         `json:"status"`
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	InternalError   string      `json:"-"`
	InternalDetails interface{} `json:"-"`
}

func (self *BaseError) GetID() string {
	return self.ID
}

func (self *BaseError) GetStatusCode() int {
	return self.StatusCode
}

func (self *BaseError) SetStatusCode(code int) {
	self.StatusCode = code
}

func (self *BaseError) GetMethod() string {
	return self.Method
}

func (self *BaseError) SetMethod(method string) {
	self.Method = method
}

func (self *BaseError) GetPath() string {
	return self.Path
}

func (self *BaseError) SetPath(path string) {
	self.Path = path
}

func (self *BaseError) GetInternalError() string {
	return self.InternalError
}

func (self *BaseError) GetInternalDetails() interface{} {
	return self.InternalDetails
}

type Error struct {
	BaseError
	Code    string      `json:"code"`
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// Default PanicHandler.
func (self *TiltController) handlePanic(ctx context.Context, v interface{}) {
	err_cont_obj, ok := v.(*ErrorResponse)
	if !ok {
		if _, ok := v.(ErrorType); ok {
			err_cont_obj = &ErrorResponse{v}
		} else {
			log.Printf("Something busted: %+v", v)
			err_obj := self.NewError(
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

	self.WriteResponse(ctx, err_cont_obj)
}

// Default SerializerErrorHandler
func (self *TiltController) handleSerializerError(ctx context.Context, err error) {
	rctx := api_router.RequestContextFromContext(ctx)
	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Default JSONSchemaErrorHandler
func (self *TiltController) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	rctx := api_router.RequestContextFromContext(ctx)
	rctx.SetStatus(400)
	errors := result.Errors()
	details := make([]string, len(errors), len(errors))
	for i, err := range result.Errors() {
		details[i] = err.String()
	}
	err := self.NewError(
		400,
		"ERR_ID_BAD_DATA",
		"Validation for json schema failed",
	)
	err.Details = details
	panic(err)
}

func (self *TiltController) handleError(ctx context.Context, err ErrorType) error {
	rctx := api_router.RequestContextFromContext(ctx)

	if err.GetStatusCode() <= 0 {
		err.SetStatusCode(500)
	}

	if err.GetMethod() == "" {
		err.SetMethod(rctx.CurrentRoute().Method())
	}

	if err.GetPath() == "" {
		err.SetPath(rctx.CurrentRoute().Path())
	}

	rctx.SetStatus(err.GetStatusCode())

	return serializers_mw.WriteSerializedResponse(ctx, &ErrorResponse{err})
}

func (self *TiltController) NewError(status int, code string, err_str string) *Error {
	return &Error{
		BaseError: BaseError{
			ID:         "ERR" + GenUUIDHex(),
			StatusCode: status,
		},
		Code:  code,
		Error: err_str,
	}
}
