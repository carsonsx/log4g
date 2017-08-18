package log4g

import "encoding/json"

func Log(level Level, arg interface{}, args ...interface{}) {
	exportLogger.Log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{}) {
	exportLogger.Panic(arg, args...)
}

func Fatal(arg interface{}, args ...interface{}) {
	exportLogger.Fatal(arg, args...)
}

func Error(arg interface{}, args ...interface{}) {
	exportLogger.Error(arg, args...)
}

func ErrorIf(arg interface{}, args ...interface{}) {
	exportLogger.ErrorIf(arg, args...)
}

func Warn(arg interface{}, args ...interface{}) {
	exportLogger.Warn(arg, args...)
}

func Info(arg interface{}, args ...interface{}) {
	exportLogger.Info(arg, args...)
}

func Debug(arg interface{}, args ...interface{}) {
	exportLogger.Debug(arg, args...)
}

func Trace(arg interface{}, args ...interface{}) {
	exportLogger.Trace(arg, args...)
}

func GetLevel() Level {
	return exportLogger.GetLevel()
}

func SetLevel(level Level) {
	exportLogger.SetLevel(level)
}

func IsLevelEnabled(level Level) bool {
	return exportLogger.IsLevel(level)
}

func IsPanicEnabled() bool {
	return exportLogger.IsPanicEnabled()
}

func IsFatalEnabled() bool {
	return exportLogger.IsFatalEnabled()
}

func IsErrorEnabled() bool {
	return exportLogger.IsErrorEnabled()
}

func IsWarnEnabled() bool {
	return exportLogger.IsWarnEnabled()
}

func IsInfoEnabled() bool {
	return exportLogger.IsInfoEnabled()
}

func IsDebugEnabled() bool {
	return exportLogger.IsDebugEnabled()
}

func IsTraceEnabled() bool {
	return exportLogger.IsTraceEnabled()
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
	exportLogger.Close()
}
