package log4g

import (
	"fmt"
	"math"
	"strconv"
)

type Level uint64

const (
	level_OFF   Level = 0
	level_PANIC Level = 100
	level_FATAL Level = 200
	level_ERROR Level = 300
	level_WARN  Level = 400
	level_INFO  Level = 500
	level_DEBUG Level = 600
	level_TRACE Level = 700
	level_ALL   Level = math.MaxUint64
)

var names = make(map[Level]string)
var alignNames = make(map[Level]string)

func initLevelName() {
	names[level_OFF] = "OFF"
	names[level_PANIC] = "PANIC"
	names[level_FATAL] = "FATAL"
	names[level_ERROR] = "ERROR"
	names[level_WARN] = "WARN"
	names[level_INFO] = "INFO"
	names[level_DEBUG] = "DEBUG"
	names[level_TRACE] = "TRACE"
	names[level_ALL] = "ALL"
}

func alignName(level Level) {
	maxNameLen := 0
	for l, n := range names {
		if l <= level {
			if len(n) > maxNameLen {
				maxNameLen = len(n)
			}
		}
		alignNames[l] = fmt.Sprintf("%"+strconv.Itoa(maxNameLen)+"s", n)
	}
}

// custom log level
func ForName(name string, intLevel uint64) Level {
	l := Level(intLevel)
	if hasLevel(l) {
		panic(fmt.Sprintf("the level %d has existed", intLevel))
	}
	names[l] = name
	alignName(gLevel)
	return l
}

// check log level
func hasLevel(l Level) bool {
	_, ok := names[l]
	return ok
}

//
func getLevelDisplayName(l Level) string {
	if name, ok := alignNames[l]; ok {
		return "[" + name + "]"
	} else {
		panic(fmt.Sprintf("invalid log level %v", l))
	}
}

func getLevelName(l Level) string {
	if name, ok := names[l]; ok {
		return name
	} else {
		panic(fmt.Sprintf("invalid log level %v", l))
	}
}

func getLevelByName(name string) Level {
	for l, n := range names {
		if n == name {
			return l
		}
	}
	panic("invalid log name " + name)
}
