package log4g

import (
	"io"
	"time"
	"sync"
	"fmt"
	"os"
	"runtime"
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
)

func newLoggerItem(level Level, prefix string, flag int, output io.Writer, calldepth int) *GenericLoggerItem {
	logger := new(GenericLoggerItem)
	logger.level = level
	logger.prefix = prefix
	logger.flag = flag
	logger.out = output
	logger.calldepth = calldepth
	return logger
}

type LoggerItem interface {
	GetLevel() Level
	Before(t time.Time)
	Log(t time.Time, level Level, arg interface{}, args ...interface{}) (n int, err error)
	After(t time.Time, n int)
	Flush()
	Close()
}

type GenericLoggerItem struct {
	mu        sync.Mutex // ensures atomic writes; protects the following fields
	out       io.Writer  // destination for output
	buf       []byte     // for accumulating text to write
	level     Level
	prefix    string
	flag      int
	stop      bool
	calldepth int
}


func (l *GenericLoggerItem) GetLevel() Level {
	return l.level
}

func (l *GenericLoggerItem) Before(t time.Time) {

}

func (l *GenericLoggerItem) Log(t time.Time, level Level, arg interface{}, args ...interface{}) (n int, err error) {

	if l.stop || level > l.level {
		return
	}

	var text string
	switch arg.(type) {
	case string:
		text = fmt.Sprintf(arg.(string), args...)
		n, err = l.output(t, l.calldepth, level, text)
	default:
		text = fmt.Sprintf(fmt.Sprintf("%v", arg), args...)
		n, err = l.output(t, l.calldepth, level, text)
	}
	if level == LEVEL_FATAL {
		os.Exit(1)
	} else if level == LEVEL_PANIC {
		panic(text)
	}

	return
}

func (l *GenericLoggerItem) output(t time.Time, calldepth int, level Level, s string) (n int, err error) {

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
	l.formatHeader(&l.buf, t, level, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	return l.out.Write(l.buf)
}

func (l *GenericLoggerItem) After(t time.Time, n int) {

}

func (l *GenericLoggerItem) Flush() {

}

func (l *GenericLoggerItem) Close() {

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

func (l *GenericLoggerItem) formatHeader(buf *[]byte, t time.Time, level Level, file string, line int) {
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
