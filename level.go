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

var (
	levelNames = []string{"OFF", "PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	names      = make(map[Level]string)

)

func initLevelName() {
	for i, name := range levelNames {
		names[Level(i*100)] = name
	}
	names[LEVEL_ALL] = "ALL"
}

func (l Level) Name() string {
	if name, ok := names[l]; ok {
		return name
	}
	return "UNKNOWN"
}

func NewAlignmentNames(level Level) *alignmentNames {
	an := &alignmentNames{make(map[Level]string)}
	an.align(level)
	return an
}

type alignmentNames struct {
	names map[Level]string
}

func (a *alignmentNames) align(level Level) {
	maxNameLen := 0
	for l, n := range names {
		if l <= level && len(n) > maxNameLen {
			maxNameLen = len(n)
		}
	}
	for l, n := range names {
		a.names[l] = fmt.Sprintf("%"+strconv.Itoa(maxNameLen)+"s", n)
	}
}

func (a *alignmentNames) Name(level Level) string {
	return names[level]
}

// custom log level
func ForName(name string, intLevel uint64) Level {
	l := Level(intLevel)
	if hasLevel(l) {
		panic(fmt.Sprintf("the level %d has existed", intLevel))
	}
	names[l] = name
	//alignName(gLevel)
	return l
}

// check log level
func hasLevel(l Level) bool {
	_, ok := names[l]
	return ok
}

func getLevelByName(name string) Level {
	for l, n := range names {
		if n == name {
			return l
		}
	}
	panic("invalid log name " + name)
}
