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
// 当LogCallstack设置为false时，所有callstack都不输出

// AddSizeRotateLogger会按照文件大小切分轮转
// AddDateRotateLogger会按照日期切分轮转
import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
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

const (
	DATE_MODE int = iota
	SIZE_MODE
)

// LogPool 日志池对象
type LogPool struct {
	Pool map[string]*Logger
}

// Logger 日志对象
type Logger struct {
	path      string // 日志输出路径
	filename  string // 日志输出文件名
	mode      int    // 日志切分模式
	maxLevel  int    // 日志接受的最高级别
	minLevel  int    // 日志接受的最低级别
	maxNumber int    // 日志最大文件数，超过则循环替代
	maxSize   int64  // 日志单个文件最大size，单位byte

	suffix        int           // 当前rotate后缀
	mux           *sync.RWMutex // 并发锁
	logFile       *os.File      // 日志文件指针
	logger        *log.Logger   // go的logger指针
	lastWriteTime time.Time
}

// Log 全局变量
var Log *LogPool
var LogCallstack bool = true

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

	AddSizeRotateLogger(name, path, maxNumber, maxSize, maxLevel, minLevel)
}

// 根据日志大小切割轮转
func AddSizeRotateLogger(
	name, path string,
	maxNumber int, maxSize int64,
	maxLevel int, minLevel int) {
	Info("add logger", name, path)
	logger := &Logger{
		path:     path,
		filename: name,
		logFile:  openFile(path, name),

		mode:      SIZE_MODE,
		maxNumber: maxNumber,
		maxSize:   maxSize,

		maxLevel: maxLevel,
		minLevel: minLevel,

		suffix: 0,
		mux:    new(sync.RWMutex),
	}
	logger.logger = log.New(logger.logFile, "", log.Ldate|log.Ltime)

	Log.Pool[name] = logger
	go logger.ListenToReopenFile()
}

// 根据日期切割轮转
func AddDateRotateLogger(name, path string, maxLevel, minLevel int) {
	logger := &Logger{
		path:     path,
		filename: name,
		logFile:  openFile(path, name),

		mode:          DATE_MODE,
		lastWriteTime: time.Now(),

		maxLevel: maxLevel,
		minLevel: minLevel,

		suffix: 0,
		mux:    new(sync.RWMutex),
	}
	logger.logger = log.New(logger.logFile, "", log.Ldate|log.Ltime)

	Log.Pool[name] = logger
	go logger.ListenToReopenFile()
}

// logInfo是为了可变参数输出而定义的接口数据类型
type logInfo []interface{}

func baseLog(level int, prefix string, msg ...interface{}) {
	msg = append(logInfo{prefix}, msg...)
	if Log != nil {
		Log.log(level, msg...)
	}
}

func Debug(msg ...interface{}) {
	baseLog(DEBUG, "[DEBUG]", msg...)
	Callstack(DEBUG)
}

func Info(msg ...interface{}) {
	baseLog(INFO, "[INFO]", msg...)
}

func Warn(msg ...interface{}) {
	baseLog(WARN, "[WARN]", msg...)
	Callstack(WARN)
}

func Error(msg ...interface{}) {
	baseLog(ERROR, "[ERROR]", msg...)
	Callstack(ERROR)
}

func Fatal(msg ...interface{}) {
	baseLog(FATAL, "[FATAL]", msg...)
	Callstack(FATAL)
	os.Exit(1)
}

func Callstack(level int) {
	if LogCallstack {
		msg := getCallstack()
		baseLog(level, "", msg...)
	}
}

func (lp *LogPool) Get(name string) *Logger {
	return lp.Pool[name]
}

func (lp *LogPool) Debug(msg ...interface{}) {
	baseLog(DEBUG, "[DEBUG]", msg...)
	Callstack(DEBUG)
}

func (lp *LogPool) Info(msg ...interface{}) {
	baseLog(INFO, "[INFO]", msg...)
}

func (lp *LogPool) Warn(msg ...interface{}) {
	baseLog(WARN, "[WARN]", msg...)
	Callstack(WARN)
}

func (lp *LogPool) Error(msg ...interface{}) {
	baseLog(ERROR, "[ERROR]", msg...)
	Callstack(ERROR)
}

func (lp *LogPool) Fatal(msg ...interface{}) {
	baseLog(FATAL, "[FATAL]", msg...)
	Callstack(FATAL)
	os.Exit(1)
}

func (lp *LogPool) log(level int, msg ...interface{}) {
	if len(lp.Pool) < 1 {
		log.Println(msg...)
	} else {
		for _, logger := range lp.Pool {
			if logger.logger != nil &&
				((logger.maxLevel >= level && logger.minLevel <= level) ||
					level == ALL) {
				logger.logAndRotate(msg...)
			}
		}
	}
}

func (lg *Logger) logAndRotate(msg ...interface{}) {
	lg.mux.Lock()
	defer lg.mux.Unlock()
	lg.rotate()
	lg.lastWriteTime = time.Now()
	lg.logger.Println(msg...)
}

func (lg *Logger) Debug(msg ...interface{}) {
	lg.log(DEBUG, "[DEBUG]", msg...)
	lg.Callstack(DEBUG)
}

func (lg *Logger) Info(msg ...interface{}) {
	lg.log(INFO, "[INFO]", msg...)
}

func (lg *Logger) Warn(msg ...interface{}) {
	lg.log(WARN, "[WARN]", msg...)
	lg.Callstack(WARN)
}

func (lg *Logger) Error(msg ...interface{}) {
	lg.log(ERROR, "[ERROR]", msg...)
	lg.Callstack(ERROR)
}

func (lg *Logger) Callstack(level int) {
	msg := getCallstack()
	lg.log(level, "", msg...)
}

func (logger *Logger) log(level int, prefix string, msg ...interface{}) {
	msg = append(logInfo{prefix}, msg...)
	if logger.logger != nil &&
		((logger.maxLevel >= level && logger.minLevel <= level) ||
			level == ALL) {
		logger.logAndRotate(msg...)
	}
}

func (lg *Logger) rotate() {
	curFilename := lg.path + "/" + lg.filename
	switch lg.mode {
	case DATE_MODE:
		if lg.lastWriteTime.Day() != time.Now().Day() {
			lg.mvFile()
		}
	case SIZE_MODE:
		if fileSize(curFilename) > lg.maxSize {
			lg.mvFile()
		}
	}
}

func (lg *Logger) mvFile() {
	curFilename := lg.path + "/" + lg.filename
	if lg.logFile != nil {
		lg.logFile.Close()
	}
	switch lg.mode {
	case DATE_MODE:
		tarFilename := lg.lastWriteTime.Format("20060102")
		//is file exist, remove it
		if fileIsExist(tarFilename) {
			os.Remove(tarFilename)
		}
		os.Rename(curFilename, tarFilename)
	case SIZE_MODE:
		lg.suffix = int((lg.suffix + 1) % lg.maxNumber)
		tarFilename := curFilename + "." + strconv.Itoa(int(lg.suffix))
		//is file exist, remove it
		if fileIsExist(tarFilename) {
			os.Remove(tarFilename)
		}
		os.Rename(curFilename, tarFilename)
	}
	lg.logFile = openFile(lg.path, lg.filename)
	lg.logger = log.New(lg.logFile, "", log.Ldate|log.Ltime)
}

func (lg *Logger) ListenToReopenFile() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1)
	for {
		<-c
		lg.mux.Lock()
		lg.logFile.Close()
		lg.logFile = openFile(lg.path, lg.filename)
		lg.logger = log.New(lg.logFile, "", log.Ldate|log.Ltime)
		lg.mux.Unlock()
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
	if strings.Contains(file, "/golang/src/") {
		return true
	}
	return false
}
