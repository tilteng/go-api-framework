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

func (self *DefaultLogger) Debug(v ...interface{}) {
	v = append([]interface{}{"[DEBUG]"}, v...)
	self.logger.Print(v)
}

func (self *DefaultLogger) Debugf(fmt string, v ...interface{}) {
	self.logger.Printf("[DEBUG] "+fmt, v...)
}

func (self *DefaultLogger) Error(v ...interface{}) {
	v = append([]interface{}{"[ERROR]"}, v...)
	self.logger.Print(v...)
}

func (self *DefaultLogger) Errorf(fmt string, v ...interface{}) {
	self.logger.Printf("[ERROR] "+fmt, v...)
}

func (self *DefaultLogger) Info(v ...interface{}) {
	v = append([]interface{}{"[INFO]"}, v...)
	self.logger.Print(v...)
}

func (self *DefaultLogger) Infof(fmt string, v ...interface{}) {
	self.logger.Printf("[INFO] "+fmt, v...)
}

func (self *DefaultLogger) Warn(v ...interface{}) {
	v = append([]interface{}{"[WARN]"}, v...)
	self.logger.Print(v...)
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
