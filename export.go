package log4g

func Log(level Level, arg interface{}, args ...interface{})  {
	std.Log(level, arg, args...)
}

func Panic(arg interface{}, args ...interface{})  {
	std.Log(PANIC, arg, args...)
}

func Fatal(arg interface{}, args ...interface{})  {
	std.Log(FATAL, arg, args...)
}

func Error(arg interface{}, args ...interface{})  {
	std.Log(ERROR, arg, args...)
}

func Warn(arg interface{}, args ...interface{})  {
	std.Log(WARN, arg, args...)
}

func Info(arg interface{}, args ...interface{})  {
	std.Log(INFO, arg, args...)
}

func Debug(arg interface{}, args ...interface{})  {
	std.Log(DEBUG, arg, args...)
}

func Trace(arg interface{}, args ...interface{})  {
	std.Log(TRACE, arg, args...)
}

func IsLevel(level Level) bool {
	return std.IsLevel(level)
}

func IsFatal() bool {
	return IsLevel(FATAL)
}

func IsError() bool {
	return IsLevel(ERROR)
}

func IsWarn() bool {
	return IsLevel(WARN)
}

func IsInfo() bool {
	return IsLevel(INFO)
}

func IsDebug() bool {
	return IsLevel(DEBUG)
}

func IsTrace() bool {
	return IsLevel(TRACE)
}