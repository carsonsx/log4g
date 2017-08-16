package log4g

import "encoding/json"

func _log(level Level, arg interface{}, args ...interface{}) {
	if f, ok := arg.(func() (arg interface{}, args []interface{})); ok {
		if loggers.IsLevel(level) {
			arg, args := f()
			if arg == nil {
				return
			}
			if args == nil {
				loggers.Log(level, arg)
			} else {
				loggers.Log(level, arg, args...)
			}
		}
	} else {
		loggers.Log(level, arg, args...)
	}
}

func Log(level Level, arg interface{}, args ...interface{}) {
	_log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{}) {
	_log(LEVEL_PANIC, arg, args...)
}

func Fatal(arg interface{}, args ...interface{}) {
	_log(LEVEL_FATAL, arg, args...)
}

func Error(arg interface{}, args ...interface{}) {
	_log(LEVEL_ERROR, arg, args...)
}

func ErrorIf(arg interface{}, args ...interface{}) {
	if arg == nil {
		return
	}
	_log(LEVEL_ERROR, arg, args...)
}

func Warn(arg interface{}, args ...interface{}) {
	_log(LEVEL_WARN, arg, args...)
}

func Info(arg interface{}, args ...interface{}) {
	_log(LEVEL_INFO, arg, args...)
}

func Debug(arg interface{}, args ...interface{}) {
	_log(LEVEL_DEBUG, arg, args...)
}

func Trace(arg interface{}, args ...interface{}) {
	_log(LEVEL_TRACE, arg, args...)
}

func GetLevel() Level {
	return gLevel
}

func SetLevel(level Level) {
	gLevel = level
}

func IsLevelEnabled(level Level) bool {
	return loggers.IsLevel(level)
}

func IsPanicEnabled() bool {
	return IsLevelEnabled(LEVEL_PANIC)
}

func IsFatalEnabled() bool {
	return IsLevelEnabled(LEVEL_FATAL)
}

func IsErrorEnabled() bool {
	return IsLevelEnabled(LEVEL_ERROR)
}

func IsWarnEnabled() bool {
	return IsLevelEnabled(LEVEL_WARN)
}

func IsInfoEnabled() bool {
	return IsLevelEnabled(LEVEL_INFO)
}

func IsDebugEnabled() bool {
	return IsLevelEnabled(LEVEL_DEBUG)
}

func IsTraceEnabled() bool {
	return IsLevelEnabled(LEVEL_TRACE)
}

func JsonString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func JsonFunc(v interface{}) func() (arg interface{}, args []interface{}) {
	return func() (arg interface{}, args []interface{}) {
		return JsonString(v), nil
	}
}

var (
	useEnvMode  bool
	useFileMode bool
)

func SetEnv(env string) {
	if useFileMode {
		panic("can not set env if programmatically load config file")
	}
	setEnv(env)
	useEnvMode = true
}

//ensure
func LoadConfig(filename string) {
	if useEnvMode {
		panic("can not programmatically load config file if set env")
	}
	err := loadConfig(filename)
	if err != nil {
		panic(err)
	}
	useEnvMode = true
}

func ReloadConfig() {
	err := reloadConfig()
	if err != nil {
		panic(err)
	}
}

func Close() {
	loggers.Close()
}
