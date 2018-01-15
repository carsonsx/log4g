package log4g

import (
	"os"
)

func newStdoutLoggerItem(level Level, prefix string, flag int, calldepth int) *StdoutLoggerItem {
	item := new (StdoutLoggerItem)
	item.GenericLoggerItem = newLoggerItem(level, prefix, flag, os.Stdout, calldepth)
	return item
}

type StdoutLoggerItem struct {
	*GenericLoggerItem
}

func newStderrLoggerItem(level Level, prefix string, flag int, calldepth int) *StderrLoggerItem {
	item := new (StderrLoggerItem)
	item.GenericLoggerItem = newLoggerItem(level, prefix, flag, os.Stdout, calldepth)
	return item
}

type StderrLoggerItem struct {
	*GenericLoggerItem
}