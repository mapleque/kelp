package http

import (
	syslog "log"
)

type loggerer interface {
	Log(tag string, msg ...interface{})
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

func Info(msg ...interface{}) {
	log.Log("INFO", msg...)
}

func Warn(msg ...interface{}) {
	log.Log("WARN", msg...)
}

func Error(msg ...interface{}) {
	log.Log("ERROR", msg...)
}

func Debug(msg ...interface{}) {
	log.Log("DEBUG", msg...)
}
