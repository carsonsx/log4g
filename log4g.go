package log4g

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
	"flag"
	"runtime/debug"
)

const (
	//log format
	Ldate                                                                                        = 1 << iota
	Ltime
	Lmicroseconds
	Llongfile
	Lshortfile
	LUTC
	LstdFlags              = Ldate | Ltime
	exportCallDepth  = 5
	customCallDepth  = 4
)

var exportLoggers = newLoggers(exportCallDepth, defaultConfigFilepath...)

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

func NewLoggers(filepath ...string) *Loggers {
	return newLoggers(customCallDepth, filepath...)
}

func newLoggers(calldepth int, filepath ...string) *Loggers {
	initLevelName()
	ls := new(Loggers)
	ls.calldepth = calldepth
	ls.LoadConfig(filepath...)
	return ls
}

type Loggers struct {
	loggers []Logger
	config *Config
	argLevel Level
	calldepth int
	closed bool
}

func (ls *Loggers) LoadConfig(filepath ...string) {
	ls.closed = true
	ls.Close()
	ls.config = NewConfig()
	for _, fp := range filepath {
		if AddFileChangedListener(fp, ls.config, loadConfig) == nil {
			break
		}
	}

	//clear loggers
	ls.loggers = []Logger{}

	if len(ls.config.Loggers) == 0 {
		ls.loggers = append(ls.loggers, newLogger(ls.GetLevel(), ls.config.Prefix, parseFlag(ls.config.Flag), os.Stdout, ls.calldepth))
	} else {
		for _, lc := range ls.config.Loggers {
			if lc.Disabled {
				continue
			}
			prefix := ls.config.Prefix
			if lc.Prefix != "" {
				prefix = lc.Prefix
			}
			flag := parseFlag(ls.config.Flag)
			if lc.Flag != "" {
				flag = parseFlag(lc.Flag)
			}
			level := GetLevelByName(ls.config.Level)
			if lc.Level != "" {
				level = GetLevelByName(lc.Level)
			}
			var logger Logger
			switch lc.Output {
			case "stdout":
				logger = newLogger(level, prefix, flag, os.Stdout, ls.calldepth)
			case "stderr":
				logger = newLogger(level, prefix, flag, os.Stderr, ls.calldepth)
			case "file":
				logger = newFileLogger(level, prefix, flag, lc.Filename, lc.MaxLines, lc.Maxsize, lc.MaxCount, lc.Daily, ls.calldepth)
			case "redis":
				logger = newRedisLogger(level, prefix, flag, lc, ls.calldepth)
			case "socket":
				logger = newSocketLogger(level, prefix, flag, lc, ls.calldepth)
			}
			if logger != nil {
				ls.loggers = append(ls.loggers, logger)
			}
		}
	}

	ls.closed = false
}

func (ls *Loggers) GetLevel() Level {
	if ls.argLevel > 0 {
		return ls.argLevel
	}
	return GetLevelByName(ls.config.Level)
}

func (ls *Loggers) SetLevel(level Level) {
	ls.argLevel = level
}

func (ls *Loggers) Panic(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_PANIC, arg, args...)
}

func (ls *Loggers) Fatal(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_FATAL, arg, args...)
}

func (ls *Loggers) Error(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_ERROR, arg, args...)
}

func (ls *Loggers) ErrorIf(arg interface{}, args ...interface{}) {
	if arg == nil {
		return
	}
	ls.Log(LEVEL_ERROR, arg, args...)
}

func (ls *Loggers) Warn(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_WARN, arg, args...)
}

func (ls *Loggers) Info(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_INFO, arg, args...)
}

func (ls *Loggers) Debug(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_DEBUG, arg, args...)
}

func (ls *Loggers) Trace(arg interface{}, args ...interface{}) {
	ls.Log(LEVEL_TRACE, arg, args...)
}

func (ls *Loggers) IsLevelEnabled(level Level) bool {
	return ls.IsLevel(level)
}

func (ls *Loggers) IsPanicEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_PANIC)
}

func (ls *Loggers) IsFatalEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_FATAL)
}

func (ls *Loggers) IsErrorEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_ERROR)
}

func (ls *Loggers) IsWarnEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_WARN)
}

func (ls *Loggers) IsInfoEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_INFO)
}

func (ls *Loggers) IsDebugEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_DEBUG)
}

func (ls *Loggers) IsTraceEnabled() bool {
	return ls.IsLevelEnabled(LEVEL_TRACE)
}

func (ls *Loggers) IsLevel(level Level) bool {
	for _, logger := range ls.loggers {
		if logger.GetLevel() >= level {
			return true
		}
	}
	return false
}

func (ls *Loggers) Log(level Level, arg interface{}, args ...interface{}) {

	//TODO add cache for logger reloading
	if ls.closed {
		return
	}

	if ls.argLevel > 0 && ls.argLevel < level  {
		return
	}

	if f, ok := arg.(func() (arg interface{}, args []interface{})); ok {
		if ls.IsLevel(level) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			arg, args = f()
		}
	}

	now := time.Now()

	for _, logger := range ls.loggers {
		logger.SetTime(now)
		n, err := logger.Log(level, arg, args...)
		if err == nil {
			logger.AfterLog(n)
		} else {
			log.Println(err)
		}
	}
}

func (ls *Loggers) Open() {
	ls.closed = false
}

func (ls *Loggers) Close() {
	ls.closed = true
	for _, logger := range ls.loggers {
		logger.Close()
	}
}

func newLogger(level Level, prefix string, flag int, output io.Writer, calldepth int) *GenericLogger {
	logger := new(GenericLogger)
	logger.level = level
	logger.prefix = prefix
	logger.flag = flag
	logger.out = output
	logger.calldepth = calldepth
	return logger
}


type Logger interface {
	GetLevel() Level
	BeforeLog()
	SetTime(t time.Time)
	Log(level Level, arg interface{}, args ...interface{}) (n int, err error)
	AfterLog(n int)
	Close()
}

type GenericLogger struct {
	mu        sync.Mutex // ensures atomic writes; protects the following fields
	out       io.Writer  // destination for output
	buf       []byte     // for accumulating text to write
	level     Level
	prefix    string
	flag      int
	stop      bool
	now       time.Time
	calldepth int
}


func (l *GenericLogger) GetLevel() Level {
	return l.level
}

func (l *GenericLogger) BeforeLog() {
}

func (l *GenericLogger) SetTime(t time.Time) {
	l.now = t
}


func (l *GenericLogger) Log(level Level, arg interface{}, args ...interface{}) (n int, err error) {

	if l.stop || level > l.level {
		return
	}

	var text string
	switch arg.(type) {
	case string:
		text = fmt.Sprintf(arg.(string), args...)
		n, err = l.Output(l.calldepth, level, text)
	default:
		text = fmt.Sprintf(fmt.Sprintf("%v", arg), args...)
		n, err = l.Output(l.calldepth, level, text)
	}
	if level == LEVEL_FATAL {
		os.Exit(1)
	} else if level == LEVEL_PANIC {
		panic(text)
	} else if level == LEVEL_ERROR {
		l.Output(l.calldepth, level, string(debug.Stack()))
	}

	return
}

func (l *GenericLogger) Output(calldepth int, level Level, s string) (n int, err error) {

	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, l.now, level, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	return l.out.Write(l.buf)
}

func (l *GenericLogger) AfterLog(n int) {

}

func (l *GenericLogger) Close() {

}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *GenericLogger) formatHeader(buf *[]byte, t time.Time, level Level, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	*buf = append(*buf, getAlignedName(level)...)
	*buf = append(*buf, ' ')

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}
