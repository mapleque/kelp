package db

import (
	syslog "log"
)

// 实现一个简单的logger，用户记录相关信息
// 用户可以通过SetLogger方法重定向log输出
// logger只要实现分级输出方法即可

type logInterface interface {
	Debug(msg ...interface{})
	Info(msg ...interface{})
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

func (lg *logger) Info(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Warn(msg ...interface{}) {
	syslog.Println(msg...)
}

func (lg *logger) Error(msg ...interface{}) {
	syslog.Println(msg...)
}
