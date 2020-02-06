package tarantula

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Debug(msg string)
	Debugf(msg string, args ...interface{})
	Error(msg string)
	Errorf(msg string, args ...interface{})
}

func NewLogger(roomName string) Logger {
	return &defaultLogger{
		prefix: fmt.Sprintf("[%q]", roomName),
	}
}

type defaultLogger struct {
	prefix string
}

func (l *defaultLogger) Info(msg string) {
	log.Printf("INFO  %v %v", l.prefix, msg)
}

func (l *defaultLogger) Infof(msg string, args ...interface{}) {
	log.Printf(fmt.Sprintf("INFO  %v %v", l.prefix, msg), args...)
}

func (l *defaultLogger) Debug(msg string) {
	log.Printf("DEBUG %v %v", l.prefix, msg)
}

func (l *defaultLogger) Debugf(msg string, args ...interface{}) {
	log.Printf(fmt.Sprintf("DEBUG %v %v", l.prefix, msg), args...)
}

func (l *defaultLogger) Error(msg string) {
	log.Printf("ERROR %v %v", l.prefix, msg)
}

func (l *defaultLogger) Errorf(msg string, args ...interface{}) {
	log.Printf(fmt.Sprintf("ERROR %v %v", l.prefix, msg), args...)
}

type dummyLogger struct{}

func (l dummyLogger) Info(msg string)                        {}
func (l dummyLogger) Infof(msg string, args ...interface{})  {}
func (l dummyLogger) Debug(msg string)                       {}
func (l dummyLogger) Debugf(msg string, args ...interface{}) {}
func (l dummyLogger) Error(msg string)                       {}
func (l dummyLogger) Errorf(msg string, args ...interface{}) {}
