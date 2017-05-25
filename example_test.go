package log4g_test

import "github.com/carsonsx/log4g"

func Example()  {

	VERBOSE := log4g.ForName("VERBOSE", 650)

	log4g.Trace(log4g.level_TRACE)
	log4g.Trace("message...")
	log4g.Trace("%d", log4g.level_TRACE)
	log4g.Trace("config -> %v", log4g.config)

	log4g.Log(VERBOSE, VERBOSE)
	log4g.Log(VERBOSE,"message...")
	log4g.Log(VERBOSE,"%d", VERBOSE)
	log4g.Log(VERBOSE,"config -> %v", log4g.config)

	log4g.Debug(log4g.level_DEBUG)
	log4g.Debug("message...")
	log4g.Debug("%d", log4g.level_DEBUG)
	log4g.Debug("config -> %v", log4g.config)

	log4g.Info(log4g.levl_INFO)
	log4g.Info("message...")
	log4g.Info("%d", log4g.levl_INFO)
	log4g.Info("config -> %v", log4g.config)

	log4g.Warn(log4g.level_WARN)
	log4g.Warn("message...")
	log4g.Warn("%d", log4g.level_WARN)
	log4g.Warn("config -> %v", log4g.config)

	log4g.Error(log4g.level_ERROR)
	log4g.Error("message...")
	log4g.Error("%d", log4g.level_ERROR)
	log4g.Error("config -> %v", log4g.config)

	log4g.Fatal(log4g.level_FATAL)
	log4g.Fatal("message...")
	log4g.Fatal("%d", log4g.level_FATAL)
	log4g.Fatal("config -> %v", log4g.config)

	log4g.Panic(log4g.level_PANIC)
	log4g.Panic("message...")
	log4g.Panic("%d", log4g.level_PANIC)
	log4g.Panic("config -> %v", log4g.config)

	// Output:
}
