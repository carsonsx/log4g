package log4g

import (
	"fmt"
	"math"
	"strconv"
)

type Level uint64

const (
	LEVEL_OFF   Level = 0
	LEVEL_PANIC Level = 100
	LEVEL_FATAL Level = 200
	LEVEL_ERROR Level = 300
	LEVEL_WARN  Level = 400
	LEVEL_INFO  Level = 500
	LEVEL_DEBUG Level = 600
	LEVEL_TRACE Level = 700
	LEVEL_ALL   Level = math.MaxUint64
)

var names = make(map[Level]string)
var alignNames = make(map[Level]string)

func initLevelName() {
	names[LEVEL_OFF] = "OFF"
	names[LEVEL_PANIC] = "PANIC"
	names[LEVEL_FATAL] = "FATAL"
	names[LEVEL_ERROR] = "ERROR"
	names[LEVEL_WARN] = "WARN"
	names[LEVEL_INFO] = "INFO"
	names[LEVEL_DEBUG] = "DEBUG"
	names[LEVEL_TRACE] = "TRACE"
	names[LEVEL_ALL] = "ALL"
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
