package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Error class used to define errors
type ErrorClass struct {
	// Name of the error class
	Name string `json:"name"`
	// ERR_ID* code
	Code string `json:"code"`
	// Default HTTP status code for this error
	Status int `json:"-"`
	// Title of error that shouldn't change between ocurrences
	Title string `json:"title"`
}

func (self *ErrorClass) newError(details string, forceFrames bool, skip int) *Error {
	err := &Error{
		ErrorClass: *self,
		ID:         Config.ErrIDGenerator.GenErrID(),
		Status:     self.Status,
		Details:    details,
	}
	if err.Status >= 500 || forceFrames {
		err.StackTrace = GetStackTrace(1 + skip)
	}
	return err
}

// Create an instance of an error
// Create an instance of an error. Takes an optional argument to use to
// set internal details
func (self *ErrorClass) New(details string) *Error {
	return self.newError(details, false, 1)
}

// Create an instance of an error, including a stacktrace. 'skip' is how
// many stack frames to skip.
func (self *ErrorClass) NewWithStack(details string, skip int) *Error {
	return self.newError(details, true, 1+skip)
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
	Status           int                    `json:"status"`
	ID               string                 `json:"id"`
	Details          string                 `json:"details,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	StackTrace       StackTrace             `json:"stack_trace,omitempty"`
	InternalError    string                 `json:"internal_error,omitempty"`
	InternalDetails  interface{}            `json:"internal_details,omitempty"`
	InternalMetadata map[string]interface{} `json:"internal_metadata,omitempty"`
}

func (self *Error) GetStatus() int {
	return self.Status
}

func (self *Error) SetStatus(status int) *Error {
	self.Status = status
	return self
}

func (self *Error) SetStackTrace(st StackTrace) *Error {
	self.StackTrace = st
	return self
}

func (self *Error) SetMetadata(v map[string]interface{}) *Error {
	self.Metadata = v
	return self
}

func (self *Error) SetInternal(v interface{}) *Error {
	self.InternalDetails = v
	if s, ok := v.(fmt.Stringer); ok {
		self.InternalError = s.String()
	} else if e, ok := v.(error); ok {
		self.InternalError = e.Error()
	} else if e, ok := v.(string); ok {
		self.InternalError = e
	}
	return self
}

func (self *Error) SetInternalMetadata(v map[string]interface{}) *Error {
	self.InternalMetadata = v
	return self
}

func (self *Error) AsJSONAPIError() *JSONAPIError {
	return &JSONAPIError{
		ID:     self.ID,
		Status: self.Status,
		Code:   self.Code,
		Title:  self.Title,
		Detail: self.Details,
		Meta:   self.Metadata,
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

func NewErrorClass(name string, code string, status int, title string) *ErrorClass {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Sprintf(
			"Couldn't determine caller defining error '%s'",
			name,
		))
	}
	if fn := runtime.FuncForPC(pc); fn != nil && len(fn.Name()) > 0 {
		fname := fn.Name()
		rindex := strings.LastIndex(fname, ".")
		if rindex == -1 {
			rindex = 0
		}
		if rindex > 0 {
			fullname := fname[:rindex] + "." + name
			if _, ok := RegisteredErrors[fullname]; ok {
				panic(fmt.Sprintf(
					"%s:%d: Error '%s' has already been defined",
					file,
					line,
					fullname,
				))
			}
			errcls := &ErrorClass{
				Name:   fullname,
				Code:   code,
				Status: status,
				Title:  title,
			}
			RegisteredErrors[fullname] = errcls
			return errcls
		}
	}
	panic(fmt.Sprintf(
		"%s:%d: Couldn't determine package for error '%s'",
		file,
		line,
		name,
	))
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
