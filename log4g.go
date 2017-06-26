package log4g

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	//log format
	ldate = 1 << iota
	ltime
	lmicroseconds
	llongfile
	lshortfile
	lutc
	lstdFlags = ldate | ltime

	calldepth = 4
)

var loggers = new(Loggers)
var gPrefix string
var gLevel Level
var gFlag = 0

func init() {
	initLevelName()
	loadDefaultConfig()
}

type Logger interface {
	BeforeLog()
	Log(level Level, arg interface{}, args ...interface{}) (n int, err error)
	AfterLog(n int)
	Close()
}

type Loggers []Logger

func (ls Loggers) IsLevel(level Level) bool {
	return level <= gLevel
}

func (ls Loggers) Log(level Level, arg interface{}, args ...interface{}) {

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	for _, logger := range ls {
		logger.BeforeLog()
		n, err := logger.Log(level, arg, args...)
		if err == nil {
			logger.AfterLog(n)
		} else {
			log.Println(err)
		}
	}
}

func (ls Loggers) Close() {
	for _, logger := range ls {
		logger.Close()
	}
}

func initLoggers() {

	if loggers != nil {
		loggers.Close()
	}

	gPrefix = Config.Prefix
	gLevel = GetLevelByName(Config.Level)
	gFlag = parseFlag(Config.Flag, ldate|ltime|lshortfile)

	alignLevelName(gLevel)

	loggers = new(Loggers)
	if len(Config.Loggers) == 0 {
		*loggers = append(*loggers, newLogger(gPrefix, gFlag, os.Stdout))
	} else {
		for _, lc := range Config.Loggers {
			if lc.Disabled {
				continue
			}
			prefix := gPrefix
			if lc.Prefix != "" {
				prefix = lc.Prefix
			}
			flag := parseFlag(lc.Flag, gFlag)
			switch lc.Output {
			case "stdout":
				if logger := newLogger(prefix, flag, os.Stdout); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "stderr":
				if logger := newLogger(prefix, flag, os.Stderr); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "file":
				if logger := newFileLogger(prefix, flag, lc.Filename, lc.MaxLines, lc.Maxsize, lc.MaxCount, lc.Daily); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "redis":
				if logger := newRedisLogger(prefix, flag, lc); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "socket":
				if logger := newSocketLogger(prefix, flag, lc); logger != nil {
					*loggers = append(*loggers, logger)
				}
			}
		}
	}
}

func newLogger(prefix string, flag int, output io.Writer) *GenericLogger {
	logger := new(GenericLogger)
	logger.prefix = prefix
	logger.flag = flag
	logger.out = output
	return logger
}

type GenericLogger struct {
	mu  sync.Mutex // ensures atomic writes; protects the following fields
	out io.Writer  // destination for output
	buf []byte     // for accumulating text to write

	prefix string
	flag   int
	stop   bool
	now    time.Time
}

func (l *GenericLogger) BeforeLog() {
	l.now = time.Now()
}

func (l *GenericLogger) Log(level Level, arg interface{}, args ...interface{}) (n int, err error) {

	if l.stop || level > gLevel {
		return
	}

	var text string
	switch arg.(type) {
	case string:
		text = fmt.Sprintf(arg.(string), args...)
		n, err = l.Output(calldepth, level, text)
	default:
		text = fmt.Sprintf(fmt.Sprintf("%v", arg), args...)
		n, err = l.Output(calldepth, level, text)
	}
	if level == LEVEL_FATAL {
		os.Exit(1)
	} else if level == LEVEL_PANIC {
		panic(text)
	}

	return
}

func (l *GenericLogger) Output(calldepth int, level Level, s string) (n int, err error) {

	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(lshortfile|llongfile) != 0 {
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
	if l.flag&lutc != 0 {
		t = t.UTC()
	}
	if l.flag&(ldate|ltime|lmicroseconds) != 0 {
		if l.flag&ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(ltime|lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	*buf = append(*buf, getAlignedName(level)...)
	*buf = append(*buf, ' ')

	if l.flag&(lshortfile|llongfile) != 0 {
		if l.flag&lshortfile != 0 {
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
		*buf = append(*buf, ' ')
	}
}
