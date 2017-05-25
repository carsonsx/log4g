package log4g


func Example()  {

	VERBOSE := ForName("VERBOSE", 650)

	Trace(level_TRACE)
	Trace("message...")
	Trace("%d", level_TRACE)
	Trace("config -> %v", config)

	Log(VERBOSE, VERBOSE)
	Log(VERBOSE,"message...")
	Log(VERBOSE,"%d", VERBOSE)
	Log(VERBOSE,"config -> %v", config)

	Debug(level_DEBUG)
	Debug("message...")
	Debug("%d", level_DEBUG)
	Debug("config -> %v", config)

	Info(level_INFO)
	Info("message...")
	Info("%d", level_INFO)
	Info("config -> %v", config)

	Warn(level_WARN)
	Warn("message...")
	Warn("%d", level_WARN)
	Warn("config -> %v", config)

	Error(level_ERROR)
	Error("message...")
	Error("%d", level_ERROR)
	Error("config -> %v", config)

	Fatal(level_FATAL)
	Fatal("message...")
	Fatal("%d", level_FATAL)
	Fatal("config -> %v", config)

	Panic(level_PANIC)
	Panic("message...")
	Panic("%d", level_PANIC)
	Panic("config -> %v", config)

	// Output:
}
