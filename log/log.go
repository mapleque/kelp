package log

// log模块封装，用于整个框架中任何地方
// 实现了日志指定大小和自动切分，日志分级输出等，参考配置文件选项
//
// import本包后直接使用各级别日志输出方法输出日志
// 日志将会被输出到所有已经添加到日志池的满足条件的logger
// 如果想指定输出的logger可以用Log.Pool[name].Debug等方法
//
// Debug 有callstack，可指定logger
// Info 可指定logger
// Warn 有callstack，可指定logger
// Error 有callstack，可指定logger
// Fatal 有callstack，且服务会停止
// Callstack 直接输出当前callstack，可指定logger

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	// 0 - 5
	ALL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// LogPool 日志池对象
type LogPool struct {
	Pool map[string]*Logger
}

// Logger 日志对象
type Logger struct {
	path      string // 日志输出路径
	filename  string // 日志输出文件名
	maxLevel  int    // 日志接受的最高级别
	minLevel  int    // 日志接受的最低级别
	maxNumber int    // 日志最大文件数，超过则循环替代
	maxSize   int64  // 日志单个文件最大size，单位byte

	suffix  int           // 当前rotate后缀
	mux     *sync.RWMutex // 并发锁
	logFile *os.File      // 日志文件指针
	logger  *log.Logger   // go的logger指针
}

// Log 全局变量
var Log *LogPool

func init() {
	if Log != nil {
		return
	}
	Info("init log module...")
	Log = &LogPool{}
	Log.Pool = make(map[string]*Logger)
}

// 增加一个日志输出模块
func AddLogger(
	name, path string,
	maxNumber int, maxSize int64,
	maxLevel int, minLevel int) {
	Info("add logger", name, path)
	logger := &Logger{}

	logger.path = path
	logger.filename = name
	logger.logFile = openFile(path, name)
	logger.maxNumber = maxNumber
	logger.maxSize = maxSize
	logger.maxLevel = maxLevel
	logger.minLevel = minLevel

	logger.suffix = 0
	logger.mux = new(sync.RWMutex)
	logger.logger = log.New(logger.logFile, "", log.Ldate|log.Ltime)

	Log.Pool[name] = logger
}

// logInfo是为了可变参数输出而定义的接口数据类型
type logInfo []interface{}

func baseLog(level int, prefix string, msg ...interface{}) {
	msg = append(logInfo{prefix}, msg...)
	if Log != nil {
		Log.log(level, msg...)
	}
	log.Println(msg...)
}

func Debug(msg ...interface{}) {
	baseLog(DEBUG, "[DEBUG]", msg...)
	Callstack()
}

func Info(msg ...interface{}) {
	baseLog(INFO, "[INFO]", msg...)
}

func Warn(msg ...interface{}) {
	baseLog(WARN, "[WARN]", msg...)
	Callstack()
}

func Error(msg ...interface{}) {
	baseLog(ERROR, "[ERROR]", msg...)
	Callstack()
}

func Fatal(msg ...interface{}) {
	baseLog(FATAL, "[FATAL]", msg...)
	Callstack()
	os.Exit(1)
}

func Callstack() {
	msg := getCallstack()
	baseLog(ALL, "", msg...)
}

func (lp *LogPool) log(level int, msg ...interface{}) {
	for _, logger := range lp.Pool {
		if logger.logger != nil &&
			((logger.maxLevel >= level && logger.minLevel <= level) ||
				level == ALL) {
			logger.rotate()
			logger.mux.RLock()
			defer logger.mux.RUnlock()
			logger.logger.Println(msg...)
		}
	}
}

func (lg *Logger) Debug(msg ...interface{}) {
	lg.log(DEBUG, "[DEBUG]", msg...)
	lg.Callstack()
}

func (lg *Logger) Info(msg ...interface{}) {
	lg.log(INFO, "[INFO]", msg...)
}

func (lg *Logger) Warn(msg ...interface{}) {
	lg.log(WARN, "[WARN]", msg...)
	lg.Callstack()
}

func (lg *Logger) Error(msg ...interface{}) {
	lg.log(ERROR, "[ERROR]", msg...)
	lg.Callstack()
}

func (lg *Logger) Callstack() {
	msg := getCallstack()
	lg.log(ALL, "", msg...)
}

func (logger *Logger) log(level int, prefix string, msg ...interface{}) {
	msg = append(logInfo{prefix}, msg...)
	if logger.logger != nil &&
		((logger.maxLevel >= level && logger.minLevel <= level) ||
			level == ALL) {
		logger.rotate()
		logger.mux.RLock()
		defer logger.mux.RUnlock()
		logger.logger.Println(msg...)
	}
	log.Println(msg...)
}

func (lg *Logger) rotate() {
	curFilename := lg.path + "/" + lg.filename
	if fileSize(curFilename) > lg.maxSize {
		lg.mux.Lock()
		defer lg.mux.Unlock()
		lg.suffix = int((lg.suffix + 1) % lg.maxNumber)
		if lg.logFile != nil {
			lg.logFile.Close()
		}
		tarFilename := curFilename + "." + strconv.Itoa(int(lg.suffix))
		//is file exist, remove it
		if fileIsExist(tarFilename) {
			os.Remove(tarFilename)
		}
		os.Rename(curFilename, tarFilename)
		lg.logFile = openFile(lg.path, lg.filename)
		lg.logger = log.New(lg.logFile, "", log.Ldate|log.Ltime)
	}
}

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func fileSize(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		// usually file is rotating now, so just return empty file
		return int64(0)
	}
	return fileInfo.Size()
}

func openFile(path, filename string) *os.File {
	pathInfo, err := os.Stat(path)
	if err != nil {
		Fatal("log path config error", err.Error())
	}
	if !pathInfo.IsDir() {
		Fatal("log path [" + path + "] is not a dir")
	}
	logFile, err := os.OpenFile(
		path+"/"+filename,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0666)
	if err != nil {
		Fatal("open log file error", err.Error())
	}
	return logFile
}

// callstack
func getCallstack() []interface{} {
	var callstack []interface{}
	callstack = append(callstack, "Callstack:\n") // start with a new line
	for skip := 0; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// remove program and framework callstack
		if !isFilterCallstack(file) {
			callstack = append(callstack, file+":"+strconv.Itoa(line)+"\n")
		}
	}
	return callstack
}

func isFilterCallstack(file string) bool {
	/*
		if strings.Contains(file, "/coral/") ||
			strings.Contains(file, "/golang/src/") {
			return true
		}
	*/
	if strings.Contains(file, "/golang/src/") {
		return true
	}
	return false
}
