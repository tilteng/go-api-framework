package logger

import (
	"io"
	"log"
)

type Logger interface {
	LogDebug(v ...interface{})
	LogDebugf(fmt string, v ...interface{})
	LogError(v ...interface{})
	LogErrorf(fmt string, v ...interface{})
	LogInfo(v ...interface{})
	LogInfof(fmt string, v ...interface{})
	LogWarn(v ...interface{})
	LogWarnf(fmt string, v ...interface{})
}

type DefaultLogger struct {
	logger *log.Logger
}

func prependString(s string, v []interface{}) []interface{} {
	nv := make([]interface{}, 1+len(v), 1+len(v))
	nv[0] = s
	copy(nv[1:], v)
	return nv
}

func (self *DefaultLogger) LogDebug(v ...interface{}) {
	self.logger.Println(prependString("[DEBUG]", v)...)
}

func (self *DefaultLogger) LogDebugf(fmt string, v ...interface{}) {
	self.logger.Printf("[DEBUG] "+fmt, v...)
}

func (self *DefaultLogger) LogError(v ...interface{}) {
	self.logger.Println(prependString("[ERROR]", v)...)
}

func (self *DefaultLogger) LogErrorf(fmt string, v ...interface{}) {
	self.logger.Printf("[ERROR] "+fmt, v...)
}

func (self *DefaultLogger) LogInfo(v ...interface{}) {
	self.logger.Println(prependString("[INFO]", v)...)
}

func (self *DefaultLogger) LogInfof(fmt string, v ...interface{}) {
	self.logger.Printf("[INFO] "+fmt, v...)
}

func (self *DefaultLogger) LogWarn(v ...interface{}) {
	self.logger.Println(prependString("[WARN]", v)...)
}

func (self *DefaultLogger) LogWarnf(fmt string, v ...interface{}) {
	self.logger.Printf("[WARN] "+fmt, v...)
}

func (self *DefaultLogger) SetWriter(out io.Writer) {
	self.logger.SetOutput(out)
}

func NewDefaultLogger(out io.Writer, prefix string) *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(out, prefix, log.LstdFlags|log.Lmicroseconds),
	}
}
