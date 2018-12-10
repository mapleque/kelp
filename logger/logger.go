package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

type Loggerer interface {
	Log(tag string, message ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

type Logger struct {
	filepath         string
	rotateSize       int64
	rotateFiles      int
	outputSetting    *tagSetting
	callstackSetting *tagSetting

	logHandler       *handler
	nextRotateNumber int
	mux              *sync.Mutex
}

// DefaultLogger A logger for user call Log method with package directly.
// By default, the output target is stdout,
// which can be redirect to a file by using method RedirectTo.
var DefaultLogger *Logger

var pool map[string]*Logger

func init() {
	DefaultLogger = New("")
	pool = map[string]*Logger{}
}

// New Create a Logger.
// If filepath is empty, it will output to stdout.
// By default, it will rotate by size limit and file limit.
// Call SetRotateSize method to change size limit, default 10M.
// Call SetRotateFiles method to change file limit, default 2.
func New(filepath string) *Logger {
	logger := &Logger{
		filepath:         filepath,
		rotateSize:       10 * 1024 * 1024,
		rotateFiles:      2,
		outputSetting:    newTagSetting().SetAll(true),
		callstackSetting: newTagSetting().SetAll(false).Set(FATAL, true).Set(ERROR, true).Set(WARN, true).Set(DEBUG, true),

		logHandler:       &handler{},
		nextRotateNumber: 0,
		mux:              new(sync.Mutex),
	}
	logger.logHandler.Open(filepath)
	return logger
}

// Add Create a Logger and add to pool with key name
func Add(name, filepath string) *Logger {
	l := New(filepath)
	pool[name] = l
	return l
}

// Get Get the Logger you have added.
// If not exist, it will return DefaultLogger.
func Get(name string) *Logger {
	if l, ok := pool[name]; ok {
		return l
	}
	return DefaultLogger
}

// RedirectTo Redirect DefaultLogger output target.
// If filepath is empty, it will output to stdout.
func RedirectTo(filepath string) *Logger {
	DefaultLogger.RedirectTo(filepath)
	return DefaultLogger
}

// Log Log with DefaultLogger.
func Log(tag string, message ...interface{}) {
	DefaultLogger.Log(tag, message...)
}

// Debug Log with DefaultLogger.
func Debug(message ...interface{}) {
	DefaultLogger.Debug(message...)
}

// Info Log with DefaultLogger.
func Info(message ...interface{}) {
	DefaultLogger.Info(message...)
}

// Warn Log with DefaultLogger.
func Warn(message ...interface{}) {
	DefaultLogger.Warn(message...)
}

// Error Log with DefaultLogger.
func Error(message ...interface{}) {
	DefaultLogger.Error(message...)
}

// Fatal Log with DefaultLogger.
func Fatal(message ...interface{}) {
	DefaultLogger.Fatal(message...)
}

// SetOutput To set all log whether output or not.
func (this *Logger) SetOutput(is bool) *Logger {
	this.outputSetting.SetAll(is)
	return this
}

// SetTagOutput To set log with this tag whether output or not.
func (this *Logger) SetTagOutput(tag string, is bool) *Logger {
	this.outputSetting.Set(tag, is)
	return this
}

// WithCallstack To set callstack whether output or not when log.
func (this *Logger) WithCallstack(is bool) *Logger {
	this.callstackSetting.SetAll(is)
	return this
}

// WithTagCallstack To set callstack whether output or not when log with this tag.
func (this *Logger) WithTagCallstack(tag string, is bool) *Logger {
	this.callstackSetting.Set(tag, is)
	return this
}

// SetRotateSize Set the max file size limit.
// It will be rotate when file size exceeding the limit.
func (this *Logger) SetRotateSize(limit int64) *Logger {
	this.rotateSize = limit
	return this
}

// SetRotateFiles Set the number of rotate files saved.
// It will override the earliest file if exeeding the number.
func (this *Logger) SetRotateFiles(number int) *Logger {
	this.rotateFiles = number
	return this
}

// RedirectTo Redirect output to filepath.
// If filepath is empty, it will output to stdout.
func (this *Logger) RedirectTo(filepath string) *Logger {
	this.lock()
	defer this.unlock()
	this.logHandler.Close()
	this.filepath = filepath
	this.nextRotateNumber = 0
	this.logHandler.Open(filepath)
	return this
}

func (this *Logger) Debug(message ...interface{}) {
	this.Log(DEBUG, message...)
}

func (this *Logger) Info(message ...interface{}) {
	this.Log(INFO, message...)
}

func (this *Logger) Warn(message ...interface{}) {
	this.Log(WARN, message...)
}

func (this *Logger) Error(message ...interface{}) {
	this.Log(ERROR, message...)
}

func (this *Logger) Fatal(message ...interface{}) {
	this.Log(FATAL, message...)
	os.Exit(1)
}

func (this *Logger) Log(tag string, message ...interface{}) {
	this.lock()
	defer this.unlock()
	if this.outputSetting.Get(tag) {
		message = append([]interface{}{fmt.Sprintf("[%s]", tag)}, message...)
		this.logHandler.Println(message...)
		if this.callstackSetting.Get(tag) {
			this.logHandler.Println(callstack()...)
		}
		this.checkRotate()
	}
}

func (this *Logger) lock() {
	this.mux.Lock()
}

func (this *Logger) unlock() {
	this.mux.Unlock()
}

func (this *Logger) checkRotate() {
	if this.filepath == "" {
		// output to stdout, no need to rotate
		return
	}
	if fileSize(this.filepath) > this.rotateSize {
		this.rotate()
	}
}

func (this *Logger) rotate() {
	this.logHandler.Close()

	tar := fmt.Sprintf(this.filepath+".%d", this.nextRotateNumber)
	this.nextRotateNumber = (this.nextRotateNumber + 1) % this.rotateFiles
	if fileIsExist(tar) {
		os.Remove(tar)
	}
	os.Rename(this.filepath, tar)

	this.logHandler.Open(this.filepath)
}

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func fileSize(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		// return empty
		return 0
	}
	return fileInfo.Size()
}

func callstack() []interface{} {
	var cs []interface{}
	cs = append(cs, "Callstack:\n") // start with a new line
	for skip := 0; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// remove golang std package callstack
		if !strings.Contains(file, "/golang/src/") {
			cs = append(cs, fmt.Sprintf("%s:%d\n", file, line))
		}
	}
	return cs
}
