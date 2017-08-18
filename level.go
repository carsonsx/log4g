package log4g

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

var (
	levelNames   = []string{"OFF", "PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	names        = make(map[Level]string)
	alignedNames = make(map[Level]string)
)

func initLevelName() {
	for i, name := range levelNames {
		names[Level(i*100)] = name
	}
	names[LEVEL_ALL] = "ALL"
	alignLevelName()
}

func (l Level) Name() string {
	if name, ok := names[l]; ok {
		return name
	}
	return "UNKNOWN"
}

func alignLevelName() {
	maxNameLen := 0
	for _, n := range names {
		if len(n) > maxNameLen {
			maxNameLen = len(n)
		}
	}
	for l, n := range names {
		alignedNames[l] = fmt.Sprintf("%"+strconv.Itoa(maxNameLen)+"s", n)
	}
}

func getAlignedName(level Level) string {
	return alignedNames[level]
}

// custom log level
func ForLevelName(name string, intLevel uint64) Level {
	l := Level(intLevel)
	if hasLevel(l) {
		panic(fmt.Sprintf("the level %d has existed", intLevel))
	}
	names[l] = name
	alignLevelName()
	return l
}

// check log level
func hasLevel(l Level) bool {
	_, ok := names[l]
	return ok
}

func GetLevelByName(name string) Level {
	upname := strings.ToUpper(name)
	for l, n := range names {
		if n == upname {
			return l
		}
	}
	panic("invalid log name " + name)
}
