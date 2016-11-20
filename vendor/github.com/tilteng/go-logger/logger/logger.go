package logger

import (
	"io"
	"log"
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(fmt string, v ...interface{})
	Error(v ...interface{})
	Errorf(fmt string, v ...interface{})
	Info(v ...interface{})
	Infof(fmt string, v ...interface{})
	Warn(v ...interface{})
	Warnf(fmt string, v ...interface{})
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

func (self *DefaultLogger) Debug(v ...interface{}) {
	self.logger.Println(prependString("[DEBUG]", v)...)
}

func (self *DefaultLogger) Debugf(fmt string, v ...interface{}) {
	self.logger.Printf("[DEBUG] "+fmt, v...)
}

func (self *DefaultLogger) Error(v ...interface{}) {
	self.logger.Println(prependString("[ERROR]", v)...)
}

func (self *DefaultLogger) Errorf(fmt string, v ...interface{}) {
	self.logger.Printf("[ERROR] "+fmt, v...)
}

func (self *DefaultLogger) Info(v ...interface{}) {
	self.logger.Println(prependString("[INFO]", v)...)
}

func (self *DefaultLogger) Infof(fmt string, v ...interface{}) {
	self.logger.Printf("[INFO] "+fmt, v...)
}

func (self *DefaultLogger) Warn(v ...interface{}) {
	self.logger.Println(prependString("[WARN]", v)...)
}

func (self *DefaultLogger) Warnf(fmt string, v ...interface{}) {
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
