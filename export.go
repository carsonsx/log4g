package log4g

func Log(level Level, arg interface{}, args ...interface{})  {
	std.Log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{})  {
	std.Log(level_PANIC, arg, args...)
}

func Fatal(arg interface{}, args ...interface{})  {
	std.Log(level_FATAL, arg, args...)
}

func Error(arg interface{}, args ...interface{})  {
	std.Log(level_ERROR, arg, args...)
}

func Warn(arg interface{}, args ...interface{})  {
	std.Log(level_WARN, arg, args...)
}

func Info(arg interface{}, args ...interface{})  {
	std.Log(level_INFO, arg, args...)
}

func Debug(arg interface{}, args ...interface{})  {
	std.Log(level_DEBUG, arg, args...)
}

func Trace(arg interface{}, args ...interface{})  {
	std.Log(level_TRACE, arg, args...)
}

func IsLevelEnabled(level Level) bool {
	return std.IsLevel(level)
}

func IsPanicEnabled() bool {
	return IsLevelEnabled(level_PANIC)
}

func IsFatalEnabled() bool {
	return IsLevelEnabled(level_FATAL)
}

func IsErrorEnabled() bool {
	return IsLevelEnabled(level_ERROR)
}

func IsWarnEnabled() bool {
	return IsLevelEnabled(level_WARN)
}

func IsInfoEnabled() bool {
	return IsLevelEnabled(level_INFO)
}

func IsDebugEnabled() bool {
	return IsLevelEnabled(level_DEBUG)
}

func IsTraceEnabled() bool {
	return IsLevelEnabled(level_TRACE)
}