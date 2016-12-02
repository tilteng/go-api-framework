package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Error ID generator type
type ErrIDGenerator interface {
	GenErrID() string
}

// Turn a single generator function into something that satisfies
// the ErrIDGenerator interface
type ErrIDGeneratorFn func() string

func (self ErrIDGeneratorFn) GenErrID() string {
	return self()
}

type NewErrorHandler interface {
	HandleNewError(context.Context, ErrorType)
}

type NewErrorHandlerFn func(context.Context, ErrorType)

func (self NewErrorHandlerFn) HandleNewError(ctx context.Context, err ErrorType) {
	self(ctx, err)
}

type ErrorManager struct {
	errorClasses    map[string]*ErrorClass
	errIDGenerator  ErrIDGenerator
	newErrorHandler NewErrorHandler
}

func (self *ErrorManager) ErrorClasses() []ErrorClass {
	classes := make([]ErrorClass, 0, len(self.errorClasses))
	for _, class := range self.errorClasses {
		classes = append(classes, *class)
	}
	return classes
}

func (self *ErrorManager) SetNewErrorHandler(handler NewErrorHandler) *ErrorManager {
	self.newErrorHandler = handler
	return self
}

func (self *ErrorManager) SetErrIDGenerator(generator ErrIDGenerator) *ErrorManager {
	self.errIDGenerator = generator
	return self
}

func (self *ErrorManager) NewClass(name string, code string, status int, title string) *ErrorClass {
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
			if existing, ok := self.errorClasses[fullname]; ok {
				panic(fmt.Sprintf(
					"%s:%d: Error '%s' was already defined at %s:%d",
					file,
					line,
					fullname,
					existing.sourceFile,
					existing.sourceLine,
				))
			}
			errcls := &ErrorClass{
				errorManager: self,
				sourceFile:   file,
				sourceLine:   line,
				Name:         fullname,
				Code:         code,
				Status:       status,
				Title:        title,
			}
			self.errorClasses[fullname] = errcls
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

// Error class used to define errors
type ErrorClass struct {
	// Pointer back to ErrorManager
	errorManager *ErrorManager
	// Filename where the error class was defined
	sourceFile string
	// Line number where the error class was defined
	sourceLine int

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
		ID:         self.errorManager.errIDGenerator.GenErrID(),
		Status:     self.Status,
		Details:    details,
	}
	if err.Status >= 500 || forceFrames {
		err.StackTrace = GetStackTrace(1 + skip)
	}
	return err
}

// Create an instance of an error. This automatically 'commits' such
// that the new error callback will be called, etc.
func (self *ErrorClass) New(ctx context.Context, details string) *Error {
	return self.newError(details, false, 1).Commit(ctx)
}

// Create an instance of an error. One must chain with Commit() to make
// sure that any new error callback is called.
func (self *ErrorClass) Start(details string) *Error {
	return self.newError(details, false, 1)
}

// Create an instance of an error, including a stacktrace. 'skip' is how
// many stack frames to skip. This automatically 'commits' such that
// the new error callback will be called, etc.
func (self *ErrorClass) NewWithStack(ctx context.Context, details string, skip int) *Error {
	return self.newError(details, true, 1+skip).Commit(ctx)
}

// Create an instance of an error, including a stacktrace. 'skip' is how
// many stack frames to skip.
func (self *ErrorClass) StartWithStack(details string, skip int) *Error {
	return self.newError(details, true, 1+skip)
}

// Interface that both Error and Errors satisfies
type ErrorType interface {
	GetName() string
	GetDetails() string
	GetStatus() int
	GetTitle() string
	GetInternalError() string
	AsJSON() (string, error)
	AsJSONAPIResponse() *JSONAPIErrorResponse
	GetStackTrace() StackTrace
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

func (self *Error) Commit(ctx context.Context) *Error {
	if hdlr := self.errorManager.newErrorHandler; hdlr != nil {
		hdlr.HandleNewError(ctx, self)
	}
	return self
}

func (self *Error) GetName() string {
	return self.Name
}

func (self *Error) GetDetails() string {
	return self.Details
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

func (self *Error) GetStackTrace() StackTrace {
	return self.StackTrace
}

func (self *Error) GetTitle() string {
	return self.Title
}

func (self *Error) GetInternalError() string {
	return self.InternalError
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

func (self Errors) GetName() string {
	if len(self) == 0 {
		return ""
	}
	return self[0].Name
}

func (self Errors) GetDetails() string {
	if len(self) == 0 {
		return ""
	}
	return self[0].Details
}

func (self Errors) GetTitle() string {
	if len(self) == 0 {
		return ""
	}
	return self[0].Title
}

func (self Errors) GetInternalError() string {
	if len(self) == 0 {
		return ""
	}
	return self[0].InternalError
}

func (self Errors) GetStackTrace() StackTrace {
	if len(self) == 0 {
		return nil
	}
	return self[0].GetStackTrace()
}

func (self Errors) GetStatus() int {
	if len(self) == 0 {
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

func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		errorClasses:   make(map[string]*ErrorClass),
		errIDGenerator: defaultErrorIDGenerator,
	}
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
