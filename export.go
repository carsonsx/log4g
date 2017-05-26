package log4g

func Log(level Level, arg interface{}, args ...interface{})  {
	loggers.Log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_PANIC, arg, args...)
}

func Fatal(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_FATAL, arg, args...)
}

func Error(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_ERROR, arg, args...)
}

func Warn(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_WARN, arg, args...)
}

func Info(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_INFO, arg, args...)
}

func Debug(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_DEBUG, arg, args...)
}

func Trace(arg interface{}, args ...interface{})  {
	loggers.Log(LEVEL_TRACE, arg, args...)
}

func GetLevel() Level {
	return gLevel
}

func SetLevel(level Level)  {
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