package log4g

import (
	"time"
)

var DBLog = NewLogger("log4g-db.json")

func Example()  {

	defer Close()

	VERBOSE := ForLevelName("VERBOSE", 650)

	Trace(LEVEL_TRACE)
	Trace("message...")
	Trace("%d", LEVEL_TRACE)
	Trace("Config -> %v", exportLoggers)

	Log(VERBOSE, VERBOSE)
	Log(VERBOSE,"message...")
	Log(VERBOSE,"%d", VERBOSE)
	Log(VERBOSE,"Config -> %v", exportLoggers)

	Debug(LEVEL_DEBUG)
	Debug("message...")
	Debug("%d", LEVEL_DEBUG)
	Debug("Config -> %v", exportLoggers)

	Info(LEVEL_INFO)
	Info("message...")
	Info("%d", LEVEL_INFO)
	Info("Config -> %v", exportLoggers)

	Warn(LEVEL_WARN)
	Warn("message...")
	Warn("%d", LEVEL_WARN)
	Warn("Config -> %v", exportLoggers)

	Error(LEVEL_ERROR)
	Error("message...")
	Error("%d", LEVEL_ERROR)
	Error("Config -> %v", exportLoggers)


	DBLog.Info("only me....")

	Fatal(LEVEL_FATAL)
	Fatal("message...")
	Fatal("%d", LEVEL_FATAL)
	Fatal("Config -> %v", exportLoggers)

	Panic(LEVEL_PANIC)
	Panic("message...")
	Panic("%d", LEVEL_PANIC)
	Panic("Config -> %v", exportLoggers)




	// Output:
}

func ExampleDead() {

	for {
		Debug(time.Now())
		time.Sleep(100)
	}

	// Output:
}
