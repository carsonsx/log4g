package log4g


func Example()  {

	VERBOSE := ForName("VERBOSE", 650)

	Trace(LEVEL_TRACE)
	Trace("message...")
	Trace("%d", LEVEL_TRACE)
	Trace("config -> %v", config)

	Log(VERBOSE, VERBOSE)
	Log(VERBOSE,"message...")
	Log(VERBOSE,"%d", VERBOSE)
	Log(VERBOSE,"config -> %v", config)

	Debug(LEVEL_DEBUG)
	Debug("message...")
	Debug("%d", LEVEL_DEBUG)
	Debug("config -> %v", config)

	Info(LEVEL_INFO)
	Info("message...")
	Info("%d", LEVEL_INFO)
	Info("config -> %v", config)

	Warn(LEVEL_WARN)
	Warn("message...")
	Warn("%d", LEVEL_WARN)
	Warn("config -> %v", config)

	Error(LEVEL_ERROR)
	Error("message...")
	Error("%d", LEVEL_ERROR)
	Error("config -> %v", config)

	Fatal(LEVEL_FATAL)
	Fatal("message...")
	Fatal("%d", LEVEL_FATAL)
	Fatal("config -> %v", config)

	Panic(level_PANIC)
	Panic("message...")
	Panic("%d", level_PANIC)
	Panic("config -> %v", config)

	// Output:
}
