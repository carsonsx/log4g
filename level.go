package log4g

import (
	"fmt"
	"math"
	"strconv"
)

type Level uint64

const (
	OFF   Level = 0
	PANIC Level = 100
	FATAL Level = 200
	ERROR Level = 300
	WARN  Level = 400
	INFO  Level = 500
	DEBUG Level = 600
	TRACE Level = 700
	ALL   Level = math.MaxUint64
)

var names = make(map[Level]string)
var alignNames = make(map[Level]string)

func initLevelName() {
	names[OFF] = "OFF"
	names[PANIC] = "PANIC"
	names[FATAL] = "FATAL"
	names[ERROR] = "ERROR"
	names[WARN] = "WARN"
	names[INFO] = "INFO"
	names[DEBUG] = "DEBUG"
	names[TRACE] = "TRACE"
	names[ALL] = "ALL"
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
	if HasLevel(l) {
		panic(fmt.Sprintf("the level %d has existed", intLevel))
	}
	names[l] = name
	alignName(gLevel)
	return l
}

// check log level
func HasLevel(l Level) bool {
	_, ok := names[l]
	return ok
}

//
func GetLevelDisplayName(l Level) string {
	if name, ok := alignNames[l]; ok {
		return name
	} else {
		panic(fmt.Sprintf("invalid log level %v", l))
	}
}

func GetLevelName(l Level) string {
	if name, ok := names[l]; ok {
		return name
	} else {
		panic(fmt.Sprintf("invalid log level %v", l))
	}
}

func GetLevelByName(name string) Level {
	for l, n := range names {
		if n == name {
			return l
		}
	}
	panic("invalid log name " + name)
}
