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

func NewError(status int, code string, err_str string) *Error {
	return &Error{
		BaseError: BaseError{
			ID:         "ERR" + GenUUIDHex(),
			StatusCode: status,
		},
		Code:  code,
		Error: err_str,
	}
}

func DefaultPanicHandler(ctx context.Context, v interface{}) {
	err_cont_obj, ok := v.(*ErrorResponse)
	if !ok {
		if _, ok := v.(ErrorType); ok {
			err_cont_obj = &ErrorResponse{v}
		} else {
			log.Printf("Something busted: %+v", v)
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

	if err_type, ok := err_cont_obj.Error.(ErrorType); ok {
		HandleErrorType(ctx, err_type)
	}

	// TODO: Do something with error?
	if serializers_mw.WriteSerializedResponse(ctx, err_cont_obj) == nil {
		return
	}

	rctx := api_router.RequestContextFromContext(ctx)
	rctx.WriteResponseString("An unknown error occurred")
}

func HandleErrorType(ctx context.Context, err ErrorType) {
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
}

func DefaultSerializerErrorHandler(ctx context.Context, err error) {
	rctx := api_router.RequestContextFromContext(ctx)
	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

func DefaultJSONSchemaErrorHandler(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	rctx := api_router.RequestContextFromContext(ctx)
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
