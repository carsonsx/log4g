package log4g

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"bufio"
)

const bufferSize = 256 * 1024

func lineCounter(filename string) int {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0660)
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

func newFileLoggerItem(level Level, prefix string, flag int, filename string, buffer bool, maxlines int, maxsize int64, maxcount int, daily bool, calldepth int) LoggerItem {

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)

	fileLogger := new(FileLoggerItem)
	fileLogger.filename = filename
	fileLogger.filedir = filepath.Dir(filename)
	fileLogger.buffer = buffer
	fileLogger.maxlines = maxlines
	fileLogger.maxsize = maxsize * 1024 * 1024
	fileLogger.maxcount = maxcount
	fileLogger.format = "%s.%0" + strconv.Itoa(len(strconv.Itoa(maxcount-1))) + "d"
	fileLogger.daily = daily
	fileLogger.lines = lineCounter(filename)

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

	var out io.Writer
	if buffer {
		out = bufio.NewWriterSize(output, bufferSize)
	} else {
		out = output
	}

	fileLogger.GenericLoggerItem = newLoggerItem(level, prefix, flag, out, calldepth)

	return fileLogger
}

type FileLoggerItem struct {
	*GenericLoggerItem
	filename string
	filedir  string
	file     *os.File
	buffer   bool
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

func (l *FileLoggerItem) Before(t time.Time) {
	if l.daily {
		l.dailyBackup(t)
	}
}

func (l *FileLoggerItem) dailyBackup(t time.Time) {
	//l.now = time.Now()
	if !l.lastTime.IsZero() {
		ltYear, ltMonth, ltDay := l.lastTime.Date()
		nowYear, nowMonth, nowDay := t.Date()
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

func (l *FileLoggerItem) After(t time.Time, n int) {

	if n <= 0 {
		return
	}

	l.lastTime = t

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

func (l *FileLoggerItem) newOutput() {
	//create new log file
	output, err := os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		l.stop = true
	}
	l.file = output
	if l.buffer {
		l.out = bufio.NewWriterSize(output, bufferSize)
	} else {
		l.out = output
	}
	l.lines = 0
	l.size = 0
	l.count++
}

func (l *FileLoggerItem) Flush() {
	if writer, ok := l.out.(*bufio.Writer); ok {
		writer.Flush()
	}
}

func (l *FileLoggerItem) Close() {
	l.file.Close()
}
