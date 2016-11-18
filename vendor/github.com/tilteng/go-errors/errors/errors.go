package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type ErrorClass struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      int    `json:"-"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (self *ErrorClass) newError(status int, forceFrames bool, skip int) *Error {
	if status == 0 {
		status = self.Status
	}
	err := &Error{
		ErrorClass: *self,
		ID:         Config.ErrIDGenerator.GenErrID(),
		Status:     status,
	}
	if status >= 500 || forceFrames {
		err.StackTrace = GetStackTrace(1 + skip)
	}
	return err
}

func (self *ErrorClass) New(status int) *Error {
	return self.newError(status, false, 0)
}

func (self *ErrorClass) NewWithStack(status int, skip int) *Error {
	return self.newError(status, true, 1+skip)
}

type Error struct {
	ErrorClass `json:"class"`
	// Default is ErrorClass.Status
	Status          int         `json:"status"`
	ID              string      `json:"id"`
	Details         string      `json:"details"`
	StackTrace      StackTrace  `json:"stack_trace"`
	InternalError   string      `json:"internal_error"`
	InternalDetails interface{} `json:"internal_details"`
}

func (self *Error) SetInternal(v interface{}) {
	self.InternalDetails = v
	if s, ok := v.(fmt.Stringer); ok {
		self.InternalError = s.String()
	} else if e, ok := v.(error); ok {
		self.InternalError = e.Error()
	}
}

func (self *Error) AsJSONAPIError() *JSONAPIError {
	return &JSONAPIError{
		ID:     self.ID,
		Status: self.Status,
		Code:   self.Code,
		Title:  self.Message,
		Detail: self.Details,
	}
}

func (self *Error) AsJSONAPIResponse() *JSONAPIErrorResponse {
	return &JSONAPIErrorResponse{
		Errors: JSONAPIErrors{
			self.AsJSONAPIError(),
		},
	}
}

func (self *Error) AsJSON() (string, error) {
	byt, err := json.Marshal(self)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

type Errors []*Error

func (self *Errors) AddError(err *Error) {
	*self = append(*self, err)
}

func (self *Errors) AsJSON() (string, error) {
	byt, err := json.Marshal(self)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func (self *Errors) AsJSONAPIResponse() *JSONAPIErrorResponse {
	errors := []*Error(*self)
	jsonapi_errors := make(JSONAPIErrors, len(errors), len(errors))
	for i, err := range errors {
		jsonapi_errors[i] = err.AsJSONAPIError()
	}
	return &JSONAPIErrorResponse{
		Errors: jsonapi_errors,
	}
}

type StackTrace []*StackFrame
type StackFrame struct {
	Function string `json:"function"`
	Filename string `json:"filename"`
	LineNo   int    `json:"lineno"`
}

func NewErrors() Errors {
	return make([]*Error, 0, 10)
}

func GetRuntimeFrames(skip int) *runtime.Frames {
	pc := make([]uintptr, 100, 100)
	num := runtime.Callers(2+skip, pc)
	return runtime.CallersFrames(pc[:num])
}

func GetStackTrace(skip int) StackTrace {
	pc := make([]uintptr, 100, 100)
	num := runtime.Callers(2+skip, pc)

	trace := make(StackTrace, 0, num)
	frames := runtime.CallersFrames(pc[:num])
	for {
		frame, more := frames.Next()

		trace = append(trace, &StackFrame{
			Function: frame.Function,
			Filename: frame.File,
			LineNo:   frame.Line,
		})

		if !more {
			break
		}
	}

	return trace
}
