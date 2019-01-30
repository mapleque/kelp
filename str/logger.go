package str

import (
	syslog "log"
)

type logInterface interface {
	Error(msg ...interface{})
}

type logger struct{}

var log logInterface

func init() {
	log = &logger{}
}

func SetLogger(logger logInterface) {
	log = logger
}

func (lg *logger) Error(msg ...interface{}) {
	syslog.Println(msg...)
}
