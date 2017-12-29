package log4g

import "encoding/json"

func Log(level Level, arg interface{}, args ...interface{}) {
	exportLoggers.Log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{}) {
	exportLoggers.Panic(arg, args...)
}

func Fatal(arg interface{}, args ...interface{}) {
	exportLoggers.Fatal(arg, args...)
}

func Error(arg interface{}, args ...interface{}) {
	exportLoggers.Error(arg, args...)
}

func ErrorIf(arg interface{}, args ...interface{}) {
	exportLoggers.ErrorIf(arg, args...)
}

func Warn(arg interface{}, args ...interface{}) {
	exportLoggers.Warn(arg, args...)
}

func Info(arg interface{}, args ...interface{}) {
	exportLoggers.Info(arg, args...)
}

func Debug(arg interface{}, args ...interface{}) {
	exportLoggers.Debug(arg, args...)
}

func Trace(arg interface{}, args ...interface{}) {
	exportLoggers.Trace(arg, args...)
}

func GetLevel() Level {
	return exportLoggers.GetLevel()
}

func SetLevel(level Level) {
	exportLoggers.SetLevel(level)
}

func IsLevelEnabled(level Level) bool {
	return exportLoggers.IsLevel(level)
}

func IsPanicEnabled() bool {
	return exportLoggers.IsPanicEnabled()
}

func IsFatalEnabled() bool {
	return exportLoggers.IsFatalEnabled()
}

func IsErrorEnabled() bool {
	return exportLoggers.IsErrorEnabled()
}

func IsWarnEnabled() bool {
	return exportLoggers.IsWarnEnabled()
}

func IsInfoEnabled() bool {
	return exportLoggers.IsInfoEnabled()
}

func IsDebugEnabled() bool {
	return exportLoggers.IsDebugEnabled()
}

func IsTraceEnabled() bool {
	return exportLoggers.IsTraceEnabled()
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

//var (
//	useEnvMode  bool
//	useFileMode bool
//)

//func SetEnv(env string) {
//	if useFileMode {
//		panic("can not set env if programmatically load config file")
//	}
//	setEnv(env)
//	useEnvMode = true
//}

//ensure
//func LoadConfig(filename string) {
//	if useEnvMode {
//		panic("can not programmatically load config file if set env")
//	}
//	err := loadConfig(filename)
//	if err != nil {
//		panic(err)
//	}
//	useEnvMode = true
//}
//
//func ReloadConfig() {
//	err := reloadConfig()
//	if err != nil {
//		panic(err)
//	}
//}

func Close() {
	exportLoggers.Close()
}
