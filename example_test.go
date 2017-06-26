package log4g

import (
	"time"
)

func Example()  {

	defer Close()

	VERBOSE := ForLevelName("VERBOSE", 650)

	Trace(LEVEL_TRACE)
	Trace("message...")
	Trace("%d", LEVEL_TRACE)
	Trace("Config -> %v", Config)

	Log(VERBOSE, VERBOSE)
	Log(VERBOSE,"message...")
	Log(VERBOSE,"%d", VERBOSE)
	Log(VERBOSE,"Config -> %v", Config)

	Debug(LEVEL_DEBUG)
	Debug("message...")
	Debug("%d", LEVEL_DEBUG)
	Debug("Config -> %v", Config)

	Info(LEVEL_INFO)
	Info("message...")
	Info("%d", LEVEL_INFO)
	Info("Config -> %v", Config)

	Warn(LEVEL_WARN)
	Warn("message...")
	Warn("%d", LEVEL_WARN)
	Warn("Config -> %v", Config)

	Error(LEVEL_ERROR)
	Error("message...")
	Error("%d", LEVEL_ERROR)
	Error("Config -> %v", Config)

	Fatal(LEVEL_FATAL)
	Fatal("message...")
	Fatal("%d", LEVEL_FATAL)
	Fatal("Config -> %v", Config)

	Panic(LEVEL_PANIC)
	Panic("message...")
	Panic("%d", LEVEL_PANIC)
	Panic("Config -> %v", Config)

	// Output:
}

func ExampleDead() {

	for {
		Debug(time.Now())
		time.Sleep(100)
	}

	// Output:
}
