package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// Error class used to define errors
type ErrorClass struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      int    `json:"-"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (self *ErrorClass) newError(forceFrames bool, skip int) *Error {
	err := &Error{
		ErrorClass: *self,
		ID:         Config.ErrIDGenerator.GenErrID(),
		Status:     self.Status,
	}
	if err.Status >= 500 || forceFrames {
		err.StackTrace = GetStackTrace(1 + skip)
	}
	return err
}

// Create an instance of an error
func (self *ErrorClass) New() *Error {
	return self.newError(false, 1)
}

// Create an instance of an error, including a stacktrace
func (self *ErrorClass) NewWithStack(skip int) *Error {
	return self.newError(true, 1+skip)
}

// Interface that both Error and Errors satisfies
type ErrorType interface {
	GetStatus() int
	AsJSON() (string, error)
	AsJSONAPIResponse() *JSONAPIErrorResponse
}

// An instance of an error
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

func (self *Error) GetStatus() int {
	return self.Status
}

func (self *Error) SetInternal(v interface{}) {
	self.InternalDetails = v
	if s, ok := v.(fmt.Stringer); ok {
		self.InternalError = s.String()
	} else if e, ok := v.(error); ok {
		self.InternalError = e.Error()
	} else if e, ok := v.(string); ok {
		self.InternalError = e
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
	if *self == nil {
		*self = make([]*Error, 1, 10)
		(*self)[0] = err
		return
	}
	*self = append(*self, err)
}

func (self Errors) AsJSON() (string, error) {
	byt, err := json.Marshal(self)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func (self Errors) GetStatus() int {
	if self == nil || len(self) == 0 {
		return 0
	}
	return self[0].GetStatus()
}

func (self Errors) AsJSONAPIResponse() *JSONAPIErrorResponse {
	if self == nil {
		return nil
	}
	jsonapi_errors := make(JSONAPIErrors, len(self), len(self))
	for i, err := range self {
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
