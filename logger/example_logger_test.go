package logger_test

import (
	"github.com/mapleque/kelp/logger"
)

func Example_simple() {
	logger.Info("Here is info log!")
	logger.Warn("Here is warn log!")
	logger.Error("Here is error log!")
	logger.Debug("Here is debug log!")

	// diy log tag
	logger.Log("TAG", "Here is tag log define by user!")

	// call os.Exis(1) after print log
	logger.Fatal("Here is fatal log!")
}

func Example_redirect() {
	logger.RedirectTo("/tmp/kelp-logger.log")
}

func Example_pool() {
	// register logger
	accessLogger := logger.Add("access_log", "/tmp/access_log.log")
	errorLogger := logger.Add("error_log", "/tmp/error_log.log")

	// get logger in other scope
	// accessLogger := logger.Get("access_log")
	// errorLogger := logger.Get("error_log")

	// this log will output in /tmp/access_log.log
	accessLogger.Info("Here is info log!")

	// this log will output in /tmp/error_log.log
	errorLogger.Error("Here is error log!")
}

func Example_tagOutput() {
	accessLogger := logger.Add("access_log", "/tmp/access_log.log")
	errorLogger := logger.Add("error_log", "/tmp/error_log.log")
	// do not output DEBUG log
	accessLogger.SetTagOutput(logger.DEBUG, false)

	// only output ERROR log
	errorLogger.SetOutput(false).SetTagOutput(logger.ERROR, true)
}

func Exampe_callstack() {
	accessLogger := logger.Add("access_log", "/tmp/access_log.log")
	// only output callstack in ERROR log
	accessLogger.WithCallstack(false).WithTagCallstack(logger.ERROR, true)
}
