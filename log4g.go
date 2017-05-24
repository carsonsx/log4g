package log4g

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var std loggers
var gLevel Level
var flags map[string]int

func init() {

	initLevelName()

	loadConfig()

	gLevel = GetLevelByName(Config.Level)

	alignName(gLevel)

	flags = make(map[string]int)
	flags["Ldate"] = log.Ldate
	flags["Ltime"] = log.Ltime
	flags["Lmicroseconds"] = log.Lmicroseconds
	flags["Llongfile"] = log.Llongfile
	flags["Lshortfile"] = log.Lshortfile
	flags["LUTC"] = log.LUTC
	flags["LstdFlags"] = log.LstdFlags

	intFlag := 0
	for _, f := range Config.Flag {
		intFlag = intFlag | flags[f]
	}

	std = append(std, newConsoleLogger(gLevel, intFlag))
	if Config.Filename != "" {
		std = append(std, newFileLogger(gLevel, intFlag, Config.Filename))
	}
}

type Logger interface {
	Log(level Level, arg interface{}, args ...interface{})
	Close()
}

type loggers []Logger

func (logs loggers) IsLevel(level Level) bool {
	return level <= gLevel
}

func (logs loggers) Log(level Level, arg interface{}, args ...interface{}) {
	for _, logger := range logs {
		logger.Log(level, arg, args...)
	}
}

func newConsoleLogger(level Level, flag int) *LogWrapper {
	return newFileLogger(level, flag, "")
}

func newFileLogger(level Level, flag int, filename string) *LogWrapper {
	logger := new(LogWrapper)
	logger.level = level
	var output io.Writer
	if filename != "" {
		os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		// Open the log file
		fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			panic(err)
		}
		output = fd
	} else {
		output = os.Stdout
	}
	logger.Logger = log.New(output, "", flag)
	return logger
}

type LogWrapper struct {
	*log.Logger
	file  *os.File
	level Level
}

func (lw *LogWrapper) Log(level Level, arg interface{}, args ...interface{}) {
	if level <= lw.level {
		var text string
		switch arg.(type) {
		case string:
			text = fmt.Sprintf(GetLevelDisplayName(level)+" "+arg.(string), args...)
			lw.Output(4, text)
		default:
			text = fmt.Sprintf(fmt.Sprintf("%s %v", GetLevelDisplayName(level), arg), args...)
			lw.Output(4, text)
		}
		if level == FATAL {
			os.Exit(1)
		} else if level == PANIC {
			panic(text)
		}
	}
}

func (lw *LogWrapper) Close() {
	if lw.file != nil {
		lw.file.Close()
	}
}
