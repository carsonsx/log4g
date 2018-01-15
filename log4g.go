package log4g

import (
	"log"
	"os"
	"time"
	"flag"
	"github.com/carsonsx/gutil"
	"runtime/debug"
)

const (
	exportCallDepth  = 5
	customCallDepth  = 4
)

var exportLoggers = newLogger(exportCallDepth, defaultConfigFilepath...)

func init() {
	argLevel := parseArgLevel()
	if argLevel != "" {
		exportLoggers.argLevel = GetLevelByName(argLevel)
	}
}

func parseArgLevel() string {
	var level string
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cmd.StringVar(&level, "log4g.level", "", "set log4g log level")
	cmd.Parse(os.Args[1:])
	return level
}

func NewLogger(filepath ...string) *Logger {
	return newLogger(customCallDepth, filepath...)
}

func newLogger(calldepth int, filepath ...string) *Logger {
	initLevelName()
	ls := new(Logger)
	ls.calldepth = calldepth
	ls.LoadConfig(filepath...)
	return ls
}

type Logger struct {
	items     []LoggerItem
	config    *Config
	argLevel  Level
	calldepth int
	closed    bool
}

func (l *Logger) LoadConfig(filepath ...string) {
	l.closed = true
	l.Close()
	l.config = NewConfig()
	gutil.ListenFirstValidJsonFile(l.config, loadConfig, filepath...)

	//clear loggers
	l.items = []LoggerItem{}

	if len(l.config.Items) == 0 {
		l.items = append(l.items, newLoggerItem(l.GetLevel(), l.config.Prefix, parseFlag(l.config.Flag), os.Stdout, l.calldepth))
	} else {
		for _, lc := range l.config.Items {
			if lc.Disabled {
				continue
			}
			prefix := l.config.Prefix
			if lc.Prefix != "" {
				prefix = lc.Prefix
			}
			flag := parseFlag(l.config.Flag)
			if lc.Flag != "" {
				flag = parseFlag(lc.Flag)
			}
			level := GetLevelByName(l.config.Level)
			if lc.Level != "" {
				level = GetLevelByName(lc.Level)
			}
			var logger LoggerItem
			switch lc.Output {
			case "stdout":
				logger = newStdoutLoggerItem(level, prefix, flag, l.calldepth)
			case "stderr":
				logger = newStderrLoggerItem(level, prefix, flag, l.calldepth)
			case "file":
				logger = newFileLoggerItem(level, prefix, flag, lc.Filename, lc.Buffer, lc.MaxLines, lc.Maxsize, lc.MaxCount, lc.Daily, l.calldepth)
			case "redis":
				logger = newRedisLoggerItem(level, prefix, flag, lc, l.calldepth)
			case "socket":
				logger = newSocketLoggerItem(level, prefix, flag, lc, l.calldepth)
			}
			if logger != nil {
				l.items = append(l.items, logger)
			}
		}
	}

	l.closed = false
}

func (l *Logger) GetLevel() Level {
	if l.argLevel > 0 {
		return l.argLevel
	}
	return GetLevelByName(l.config.Level)
}

func (l *Logger) SetLevel(level Level) {
	l.argLevel = level
}

func (l *Logger) Panic(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_PANIC, arg, args...)
}

func (l *Logger) Fatal(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_FATAL, arg, args...)
}

func (l *Logger) Error(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_ERROR, arg, args...)
}

func (l *Logger) ErrorIf(arg interface{}, args ...interface{}) {
	if arg == nil {
		return
	}
	l.Log(LEVEL_ERROR, arg, args...)
}

func (l *Logger) ErrorStack(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_ERROR, arg, args...)
	l.Log(LEVEL_ERROR, debug.Stack())
}

func (l *Logger) Warn(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_WARN, arg, args...)
}

func (l *Logger) Info(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_INFO, arg, args...)
}

func (l *Logger) Debug(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_DEBUG, arg, args...)
}

func (l *Logger) Trace(arg interface{}, args ...interface{}) {
	l.Log(LEVEL_TRACE, arg, args...)
}

func (l *Logger) IsLevelEnabled(level Level) bool {
	return l.IsLevel(level)
}

func (l *Logger) IsPanicEnabled() bool {
	return l.IsLevelEnabled(LEVEL_PANIC)
}

func (l *Logger) IsFatalEnabled() bool {
	return l.IsLevelEnabled(LEVEL_FATAL)
}

func (l *Logger) IsErrorEnabled() bool {
	return l.IsLevelEnabled(LEVEL_ERROR)
}

func (l *Logger) IsWarnEnabled() bool {
	return l.IsLevelEnabled(LEVEL_WARN)
}

func (l *Logger) IsInfoEnabled() bool {
	return l.IsLevelEnabled(LEVEL_INFO)
}

func (l *Logger) IsDebugEnabled() bool {
	return l.IsLevelEnabled(LEVEL_DEBUG)
}

func (l *Logger) IsTraceEnabled() bool {
	return l.IsLevelEnabled(LEVEL_TRACE)
}

func (l *Logger) IsLevel(level Level) bool {
	for _, logger := range l.items {
		if logger.GetLevel() >= level {
			return true
		}
	}
	return false
}

func (l *Logger) Log(level Level, arg interface{}, args ...interface{}) {

	if l.closed {
		return
	}

	if l.argLevel > 0 && l.argLevel < level  {
		return
	}

	if f, ok := arg.(func() (arg interface{}, args []interface{})); ok {
		if l.IsLevel(level) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			arg, args = f()
		}
	}

	now := time.Now()
	for _, item := range l.items {
		item.Before(now)
		n, err := item.Log(now, level, arg, args...)
		if err == nil {
			item.After(now, n)
		} else {
			log.Println(err)
		}
	}
}

func (l *Logger) Open() {
	l.closed = false
}

func (l *Logger) Flush() {
	for _, logger := range l.items {
		logger.Flush()
	}
}

func (l *Logger) Close() {
	l.closed = true
	for _, logger := range l.items {
		logger.Close()
	}
}

