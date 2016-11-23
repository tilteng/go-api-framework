package rollbar

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

type TraceNotification struct {
	baseNotification
	notifierTraceBody `json:"body"`
}

type TraceChainNotification struct {
	baseNotification
	notifierTraceChainBody `json:"body"`
}

// 'body' container for a 'trace' notification
type notifierTraceBody struct {
	Trace NotifierTrace `json:"trace,omitempty"`
}

// 'body' container for a 'trace_chain' notification
type notifierTraceChainBody struct {
	TraceChain []*NotifierTrace `json:"trace_chain,omitempty"`
}

// 'trace' object used in 'trace' and 'trace_chain' notifications
type NotifierTrace struct {
	// Required array describing stack frames
	Frames []*NotifierFrame `json:"frames"`
	// Required object describing the exception
	Exception *NotifierException `json:"exception"`
}

// Exception info to use in NotifierTrace
type NotifierException struct {
	// Exception class (required)
	Class string `json:"class"`

	// Exception message as string (optional)
	Message string `json:"message,omitempty"`

	// Optional human-readable string describing the exception
	Description string `json:"description,omitempty"`
}

// Frame info to use in NotifierTrace
type NotifierFrame struct {
	// Required filename
	Filename string `json:"filename"`

	// Optional line number
	Line int `json:"lineno"`

	// Optional column number
	Column int `json:"colno,omitempty"`

	// Optional method/function name
	Method string `json:"method,omitempty"`

	// Optional line of code
	Code string `json:"code,omitempty"`

	// Optional additional code before and after the code line
	Context *NotifierCodeContext `json:"context,omitempty"`

	// Optional list of names of the arguments to method/function call
	ArgSpec []string `json:"argspec,omitempty"`

	// From API docs:
	// Optional: varargspec
	// If the function call takes an arbitrary number of unnamed pos    itional arguments,
	// the name of the argument that is the list containing those ar    guments.
	// For example, in Python, this would typically be "args" when "    *args" is used.
	// The actual list will be found in locals.
	VarargSpec string `json:"varargspec,omitempty"`

	// From API docs:
	// Optional: keywordspec
	// If the function call takes an arbitrary number of keyword arguments, the name
	// of the argument that is the object containing those arguments.
	// For example, in Python, this would typically be "kwargs" when "**kwargs" is used.
	// The actual object will be found in locals.
	KeywordSpec string `json:"keywordspec,omitempty"`

	// Optional: locals
	// Object of local variables for the method/function call.
	// The values of variables from argspec, vararspec and keywordspec
	// can be found in locals.
	Locales map[string]interface{} `json:"locales,omitempty"`
}

// Context info to use in NotifierFrame
type NotifierCodeContext struct {
	// Array of lines before code line
	Pre []string `json:"pre,omitempty"`

	// Array of lines after code line
	Post []string `json:"post,omitempty"`
}

func (self *NotifierTrace) AddExceptionFromError(err error) error {
	if self.Exception != nil {
		return errors.New("Already added an exception")
	}
	title := ""
	cls := "<nil>"
	if err != nil {
		title = err.Error()
		cls = reflect.TypeOf(err).String()
		if cls == "" {
			cls = "<unknown>"
		} else {
			cls = strings.TrimPrefix(cls, "*")
		}
	}

	self.Exception = &NotifierException{
		Class:   cls,
		Message: title,
	}

	return nil
}

func (self *NotifierTrace) AddRuntimeFrames(frames *runtime.Frames) error {
	if self.Frames != nil {
		return errors.New("Already added frames")
	}

	if frames == nil {
		pc := make([]uintptr, 100, 100)
		num := runtime.Callers(3, pc)
		frames = runtime.CallersFrames(pc[:num])
	}

	notif_frames := make([]*NotifierFrame, 0, 100)
	for {
		fr, more := frames.Next()
		notif_frames = append(
			notif_frames,
			&NotifierFrame{
				Filename: fr.File,
				Line:     fr.Line,
				Method:   fr.Function,
			},
		)

		if !more {
			break
		}
	}

	self.Frames = notif_frames

	return nil
}
