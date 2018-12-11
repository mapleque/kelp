package http

import (
	syslog "log"
)

type loggerer interface {
	Log(tag string, msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type logger struct{}

var log loggerer

func init() {
	log = &logger{}
}

func SetLogger(logger loggerer) {
	log = logger
}

func (lg *logger) Log(tag string, msg ...interface{}) {
	syslog.Println(append([]interface{}{tag}, msg...))
}

func (lg *logger) Debug(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Error(msg ...interface{}) {
	syslog.Println(msg...)
}
