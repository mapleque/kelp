Logger Package
====

本组建主要用于简化在程序中输出日志的复杂度，提供一些简单的输出日志接口和配置接口。

开箱即用
----

下面的例子展示了logger包开箱即用的日志输出方法。

```
import "github.com/mapleque/kelp/logger"

func main() {
  logger.Info("Here is info log!")
  logger.Warn("Here is warn log!")
  logger.Error("Here is error log!")
  logger.Debug("Here is debug log!")

  // diy log tag
  logger.Log("TAG", "Here is tag log define by user!")

  // call os.Exis(1) after print log
  logger.Fatal("Here is fatal log!")
}
```

输出到文件
----
上一节的例子中所有的日志将会输出到控制台。自然的，logger包提供了将日志输出重定向到指定文件的方法：

```
logger.RedirectTo("/tmp/kelp-logger.log")
```

一旦调用上面的方法，之后所有的日志都将输出到`/tmp/kelp-logger.log`文件中。前提是运行这段代码的用户具有指定文件的写权限。

特别的，当参数为空字符串时，日志输出将会被重新重定向到控制台。

使用多目标日志源
----
logger包提供了日志源池，用户可以通过注册多个日志源从而将日志输出到不同文件。

```
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
```

特别的，当用户尝试获取未注册的logger时，系统会返回一个包级别的默认logger。

日志文件轮转切割
----
对于每一个通过logger包创建的日志源，都可以单独设置其轮转切割规则。

```
// set the max file size limit, unit byte(b)
accessLogger.SetRotateSize(10*1024*1024)

// set the number of rotate files
accessLogger.SetRotateFiles(2)

// setting above, you will get 2 extra file for rotate
//   access_log.log
//   access_log.log.0
//   access_log.log.1
```

其中`access_log.log`文件将始终作为当前正在写的日志文件，带有`.n`后缀的文件将根据时间顺序先后创建，如果达到了用户设置的最大文件数，将会优先覆盖最早生成的文件。

选择性输出日志
----
对于每一个通过logger包创建的日志源，还可以单独设置不输出指定tag的日志。

```
// do not output DEBUG log
accessLogger.SetTagOutput(logger.DEBUG, false)

// only output ERROR log
errorLogger.SetOutput(false).SetTagOut(logger.ERROR, true)
```

其中是否输出某个tag的日志，由当前输出语句前最后一次调用的设置决定。

自定义是否输出调用栈
----

logger包提供了输出调用栈方法。

```
// only output callstack in ERROR log
accessLogger.WithCallstack(false).WithTagCallstack(logger.ERROR, true)
```

默认的，在预定义的tag中，除INFO日志外都会输出调用栈信息。
