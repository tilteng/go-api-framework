package logger

import "os"

var defaultStdoutLogger = NewDefaultLogger(os.Stdout, "")
var defaultStdoutCtxLogger = NewDefaultCtxLogger(defaultStdoutLogger)

func DefaultStdoutLogger() Logger {
	return defaultStdoutLogger
}

func DefaultStdoutCtxLogger() CtxLogger {
	return defaultStdoutCtxLogger
}
