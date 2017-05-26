package log4g

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

var std loggers
var gLevel Level
var gFlag = 0
var gPrefix = ""

func init() {

	initLevelName()

	loadConfig()

	gLevel = getLevelByName(config.Level)

	//alignName(gLevel)

	flags := strings.Split(config.Flag, "|")
	for _, name := range flags {
		gFlag = gFlag | getFlagByName(name)
	}

	if config.Filename != "" {
		std = append(std, NewFileLogger(gLevel, config.Filename))
	} else {
		std = append(std, NewConsoleLogger(gLevel))
	}
}

func getFlagByName(name string) int {
	flags := make(map[string]int)
	flags["date"] = ldate
	flags["time"] = ltime
	flags["microseconds"] = lmicroseconds
	flags["longfile"] = llongfile
	flags["shortfile"] = lshortfile
	flags["UTC"] = lutc
	flags["stdFlags"] = lstdFlags
	return flags[name]
}

type Logger interface {
	Log(level Level, arg interface{}, args ...interface{})
}

type loggers []Logger

func (ls loggers) IsLevel(level Level) bool {
	return level <= gLevel
}

func (ls loggers) Log(level Level, arg interface{}, args ...interface{}) {
	for _, logger := range ls {
		logger.Log(level, arg, args...)
	}
}

func NewConsoleLogger(level Level) *loggerWrapper {
	return newLogger(level, os.Stdout)
}

func NewFileLogger(level Level, filename string) *loggerWrapper {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	return newLogger(level, output)
}

func newLogger(level Level, output io.Writer) *loggerWrapper {
	logger := new(loggerWrapper)
	logger.level = level
	logger.levelNames = NewAlignmentNames(level)
	logger.out = output
	return logger
}

type loggerWrapper struct {
	mu    sync.Mutex // ensures atomic writes; protects the following fields
	out   io.Writer  // destination for output
	buf   []byte     // for accumulating text to write
	level Level
	levelNames *alignmentNames
}

func (l *loggerWrapper) Log(level Level, arg interface{}, args ...interface{}) {
	if level <= l.level {
		var text string
		switch arg.(type) {
		case string:
			text = fmt.Sprintf(arg.(string), args...)
			l.Output(calldepth, level, text)
		default:
			text = fmt.Sprintf(fmt.Sprintf("%v", arg), args...)
			l.Output(calldepth, level, text)
		}
		if level == LEVEL_FATAL {
			os.Exit(1)
		} else if level == LEVEL_PANIC {
			panic(text)
		}
	}
}

func (l *loggerWrapper) Output(calldepth int, level Level, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if gFlag&(lshortfile|llongfile) != 0 {
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
	l.formatHeader(&l.buf, now, level, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// SetOutput sets the output destination for the logger.
func (l *loggerWrapper) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
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

func (l *loggerWrapper) formatHeader(buf *[]byte, t time.Time, level Level, file string, line int) {
	*buf = append(*buf, gPrefix...)
	if gFlag&lutc != 0 {
		t = t.UTC()
	}
	if gFlag&(ldate|ltime|lmicroseconds) != 0 {
		if gFlag&ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if gFlag&(ltime|lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if gFlag&lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	*buf = append(*buf, l.levelNames.Name(level)...)
	*buf = append(*buf, ' ')

	if gFlag&(lshortfile|llongfile) != 0 {
		if gFlag&lshortfile != 0 {
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
