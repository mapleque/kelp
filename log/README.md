Log package
====

Init
----

```
log.AddLogger(
    "name",                 // log file name, also use as logger name
    "/var/log/project",     // log path
    5,                      // max log file number
    1000000,                // max file size, byte
    5,                      // max log level
    1,                      // min log level
)
```

Level
----

Log level:
- 1:DEBUG
- 2:INFO
- 3:WARN
- 4:ERROR
- 5:FATAL

Log api:
- Debug
- Info
- Warn
- Error
- Fatal

Target
----

Only write one logger.
```
log.Log.Get("name").Info()
```

Rotate
----

Size rotate
```
log.AddSizeRotateLogger(
    "name",                 // log file name, also use as logger name
    "/var/log/project",     // log path
    5,                      // max log file number
    1000000,                // max file size, byte
    5,                      // max log level
    1,                      // min log level
)
```

Date rotate
```
log.AddDateRotateLogger(
    "name",                 // log file name, also use as logger name
    "/var/log/project",     // log path
    5,                      // max log level
    1,                      // min log level
)
```

External rotate
```
kill -SIGUSR1 binary_file
```
