package db

import (
	syslog "log"
)

// 实现一个简单的logger，记录相关信息
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
	msg = append([]interface{}{"[Debug][kelp.db]"}, msg...)
	syslog.Println(msg...)
}

func (lg *logger) Info(msg ...interface{}) {
	msg = append([]interface{}{"[Info][kelp.db]"}, msg...)
	syslog.Println(msg...)
}

func (lg *logger) Warn(msg ...interface{}) {
	msg = append([]interface{}{"[Warn][kelp.db]"}, msg...)
	syslog.Println(msg...)
}

func (lg *logger) Error(msg ...interface{}) {
	msg = append([]interface{}{"[Error][kelp.db]"}, msg...)
	syslog.Println(msg...)
}
