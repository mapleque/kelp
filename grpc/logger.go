package grpc

import (
	syslog "log"
	"os"
)

type logInterface interface {
	Debug(msg ...interface{})
	Info(msg ...interface{})
	Error(msg ...interface{})
	Warn(msg ...interface{})
	Fatal(msg ...interface{})
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

func (lg *logger) Info(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Warn(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Error(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Fatal(msg ...interface{}) {
	syslog.Println(msg...)
	os.Exit(1)
}
