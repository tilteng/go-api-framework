package logger

import "os"

var defaultStdoutLogger = NewDefaultLogger(os.Stdout, "")

func DefaultStdoutLogger() Logger {
	return defaultStdoutLogger
}
