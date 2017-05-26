package log4g

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"strings"
	"bytes"
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

var loggers *Loggers
var gPrefix string
var gLevel Level
var gFlag = 0

func init() {
	initLevelName()
	reload()
}

type Logger interface {
	Log(level Level, arg interface{}, args ...interface{})
	Handle()
	Close()
}

type Loggers []Logger

func (ls Loggers) IsLevel(level Level) bool {
	return level <= gLevel
}

func (ls Loggers) Log(level Level, arg interface{}, args ...interface{}) {
	for _, logger := range ls {
		logger.Log(level, arg, args...)
		logger.Handle()
	}
}
func (ls Loggers) Close() {
	for _, logger := range ls {
		logger.Close()
	}
}

func reload()  {
	loadConfig()
	initLoggers()
}

func initLoggers()  {

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
		for _, logger := range Config.Loggers {
			if logger.Disabled {
				continue
			}
			prefix := gPrefix
			if logger.Prefix != "" {
				prefix = logger.Prefix
			}
			flag := parseFlag(logger.Flag, gFlag)
			switch logger.Output {
			case "stdout":
				*loggers = append(*loggers, newLogger(prefix, flag, os.Stdout))
			case "stderr":
				*loggers = append(*loggers, newLogger(prefix, flag, os.Stderr))
			case "file":
				*loggers = append(*loggers, NewFileLogger(prefix, flag, logger.Filename, logger.Maxlines, logger.Maxsize, logger.Daily))
			}
		}
	}

}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func NewFileLogger(prefix string, flag int, filename string, maxlines int, maxsize int64, daily bool) Logger {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	fileLogger := new (FileLogger)
	fileLogger.file = output
	fileLogger.filename = filename
	fileLogger.maxlines = maxlines
	fileLogger.maxsize = maxsize
	fileLogger.daily = daily
	fileLogger.GenericLogger = newLogger(prefix, flag, output)
	lines, err := lineCounter(output)
	if err == nil {
		fileLogger.lines = lines
	}
	info, err := output.Stat()
	if err == nil {
		fileLogger.size = info.Size()
	}

	filepath.Walk(filepath.Dir(filename), func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Join(path, info.Name()), filename) {
			fileLogger.count++
		}
		return nil
	})

	return fileLogger
}

func newLogger(prefix string, flag int, output io.Writer) *GenericLogger {
	logger := new(GenericLogger)
	logger.prefix = prefix
	logger.flag = flag
	logger.out = output
	return logger
}

type GenericLogger struct {
	mu    sync.Mutex // ensures atomic writes; protects the following fields
	out   io.Writer  // destination for output
	buf   []byte     // for accumulating text to write

	prefix string
	flag int
	lastWrittenCount int
}

func (l *GenericLogger) Log(level Level, arg interface{}, args ...interface{}) {
	if level <= gLevel {
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

func (l *GenericLogger) Output(calldepth int, level Level, s string) error {
	now := time.Now() // get this early.
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
	l.formatHeader(&l.buf, now, level, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	n, err := l.out.Write(l.buf)
	if err == nil {
		l.lastWrittenCount = n
	} else {
		l.lastWrittenCount = 0
	}

	return err
}

func (l *GenericLogger) Close() {

}

func (l *GenericLogger) Handle() {

}

// SetOutput sets the output destination for the logger.
func (l *GenericLogger) SetOutput(w io.Writer) {
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
		*buf = append(*buf, ": "...)
	}
}

type FileLogger struct {
	*GenericLogger
	filename string
	file     *os.File
	maxlines int
	maxsize  int64
	daily    bool
	lines    int
	size     int64
	count    int
}

func (l *FileLogger) Handle() {

	if l.lastWrittenCount <= 0 {
		return
	}

	l.lines++
	l.size += int64(l.lastWrittenCount)
	if (l.maxlines > 0 && l.lines >= l.maxlines) || (l.maxsize > 0 && l.size > l.maxsize) {
		l.Close()
		//try to rename log files
		for i := l.count; i > 0; i-- {
			var err error
			if i == 1 {
				err = os.Rename(l.filename, l.filename + ".001")
			} else {
				err = os.Rename(fmt.Sprintf("%s.%3d", l.filename, i-1), fmt.Sprintf("%s.%3d", l.filename, i))
			}
			if err != nil {
				//stop this logger
				break
			}
		}
		output, err := os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			panic(err)
		}
		l.SetOutput(output)
		l.lines = 0
		l.size = 0
	}
}

func (l *FileLogger) Close() {
	l.Close()
}
