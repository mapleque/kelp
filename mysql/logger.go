package mysql

import (
	syslog "log"
)

type logInterface interface {
	Debug(msg ...interface{})
	Error(msg ...interface{})
	Warn(msg ...interface{})
}

type logger struct{}

var log logInterface

func init() {
	log = &logger{}
}

func SetLogger(logger logInterface) {
	log = logger
}

func (lg *logger) Debug(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Warn(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Error(msg ...interface{}) {
	syslog.Println(msg...)
}
