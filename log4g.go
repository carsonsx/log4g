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
	"log"
	"strconv"
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
	for _, logger := range ls {
		logger.BeforeLog()
		n, err := logger.Log(level, arg, args...)
		if err == nil {
			logger.AfterLog(n)
		}
	}
}

func (ls Loggers) Close() {
	for _, logger := range ls {
		logger.Close()
	}
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
		for _, loggerConfig := range Config.Loggers {
			if loggerConfig.Disabled {
				continue
			}
			prefix := gPrefix
			if loggerConfig.Prefix != "" {
				prefix = loggerConfig.Prefix
			}
			flag := parseFlag(loggerConfig.Flag, gFlag)
			switch loggerConfig.Output {
			case "stdout":
				if logger := newLogger(prefix, flag, os.Stdout); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "stderr":
				if logger := newLogger(prefix, flag, os.Stderr); logger != nil {
					*loggers = append(*loggers, logger)
				}
			case "file":
				if logger := NewFileLogger(prefix, flag, loggerConfig.Filename, loggerConfig.Maxlines, loggerConfig.Maxsize, loggerConfig.MaxCount, loggerConfig.Daily); logger != nil {
					*loggers = append(*loggers, logger)
				}
			}
		}
	}
}

func lineCounter(filename string) int {
	file, err := os.OpenFile(filename, os.O_RDONLY | os.O_CREATE, 0660)
	if err != nil && !os.IsNotExist(err) {
		return 0
	}
	defer file.Close()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := file.Read(buf)
		if err == nil {
			count += bytes.Count(buf[:c], lineSep)
		} else if err == io.EOF {
			return count
		} else {
			return 0
		}
	}
}

func NewFileLogger(prefix string, flag int, filename string, maxlines int, maxsize int64, maxcount int, daily bool) Logger {

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)

	fileLogger := new (FileLogger)
	fileLogger.filename = filename
	fileLogger.filedir = filepath.Dir(filename)
	fileLogger.maxlines = maxlines
	fileLogger.maxsize = maxsize * 1024 * 1024
	fileLogger.maxcount = maxcount
	fileLogger.format = "%s.%0" + strconv.Itoa(len(strconv.Itoa(maxcount-1))) +  "d"
	fileLogger.daily = daily
	fileLogger.lines  = lineCounter(filename)

	output, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Print(err)
		return nil
	}
	fileLogger.file = output
	info, err := output.Stat()
	if err != nil {
		log.Print(err)
		return nil
	}

	fileLogger.size = info.Size()
	fileLogger.lastTime = info.ModTime()
	//for test
	//fileLogger.lastTime = info.ModTime().Add(- 24 * time.Hour)

	filepath.Walk(fileLogger.filedir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(filename)) {
			fileLogger.count++
		}
		return nil
	})

	fileLogger.GenericLogger = newLogger(prefix, flag, output)

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

	prefix    string
	flag      int
	stop      bool
	now  time.Time
}

func (l *GenericLogger) BeforeLog() {
	l.now = time.Now()
}

func (l *GenericLogger) Log(level Level, arg interface{}, args ...interface{}) (n int ,err error) {

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

type FileLogger struct {
	*GenericLogger
	filename string
	filedir string
	file     *os.File
	maxlines int
	maxsize  int64
	maxcount int
	daily    bool
	lines    int
	size     int64
	count    int
	format   string
	lastTime time.Time
}

func (l *FileLogger) BeforeLog() {
	if l.daily {
		l.dailyBackup()
	}
}

func (l *FileLogger) dailyBackup() {
	l.now = time.Now()
	if !l.lastTime.IsZero() {
		ltYear, ltMonth, ltDay := l.lastTime.Date()
		nowYear, nowMonth, nowDay := l.now.Date()
		if ltDay != nowDay || ltMonth != nowMonth || ltYear != nowYear {

			l.mu.Lock()
			defer l.mu.Unlock()

			strDate := fmt.Sprintf("%d%02d%02d", ltYear, ltMonth, ltDay)
			dateDir := filepath.Join(l.filedir, strDate)
			err := os.MkdirAll(dateDir, os.ModePerm)
			if err == nil {
				l.Close()
				//move all file to date director
				err = filepath.Walk(l.filedir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					if strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(l.filename)) {
						os.Remove(filepath.Join(l.filedir, strDate, info.Name()))
						return os.Rename(filepath.Join(l.filedir, info.Name()), filepath.Join(l.filedir, strDate, info.Name()))
					}
					return nil
				})
				if err != nil {
					Error(err)
					return
				}
				l.count = 0
				l.newOutput()
			} else {
				Error(err)
			}
		}
	}
}


func (l *FileLogger) newOutput() {
	//create new log file
	output, err := os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		l.stop = true
	}
	l.file = output
	l.out = output
	l.lines = 0
	l.size = 0
	l.count++
}

func (l *FileLogger) AfterLog(n int) {

	if n <= 0 {
		return
	}

	l.lastTime = l.now

	l.mu.Lock()
	defer l.mu.Unlock()

	l.lines++
	l.size += int64(n)
	if (l.maxlines > 0 && l.lines >= l.maxlines) || (l.maxsize > 0 && l.size > l.maxsize) {

		log.Printf("lines=%d,maxlines=%d,count=%d", l.lines, l.maxlines, l.count)

		//close log file
		l.Close()

		//remove the oldest log
		if l.count == l.maxcount {
			if os.Remove(fmt.Sprintf(l.format, l.filename, l.maxcount-1)) != nil {
				l.stop = true
				return
			}
			l.count--
		}

		//try to rename log files
		var err error
		for i := l.count; i > 0; i-- {
			var oldpath string
			if i == 1 {
				oldpath = l.filename
			} else {
				oldpath = fmt.Sprintf(l.format, l.filename, i-1)
			}
			newpath := fmt.Sprintf(l.format, l.filename, i)
			err = os.Rename(oldpath, newpath)
			if err != nil {
				log.Println(err)
				l.stop = true
				return
			}
		}

		l.newOutput()

	}
}

func (l *FileLogger) Close() {
	l.file.Close()
}
