package log4g

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"flag"
)

const config_file_path = "log4g.json"

var config struct{
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	Flag     string `json:"flag"`
	Filename string `json:"filename"`
}

func loadConfig() {

	// default
	config.Level = LEVEL_DEBUG.Name()
	config.Flag = "date|microseconds|shortfile"

	// load form config file
	if _, err := os.Stat(config_file_path); err == nil {
		data, err := ioutil.ReadFile(config_file_path)
		if err != nil {
			log.Print(err)
			return
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	// override from os arguments
	if len(os.Args) > 1 {
		fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		fs.StringVar(&config.Prefix, "log4g-prefix", config.Prefix, "set log4g prefix")
		fs.StringVar(&config.Level, "log4g-level", config.Level, "set log4g level")
		fs.StringVar(&config.Flag, "log4g-flag", config.Flag, "set log4g flag, separated by '|'")
		fs.StringVar(&config.Filename, "log4g-filename", config.Filename, "set log4g filename")
		for i := 1; i < len(os.Args); i++ {
			 if fs.Parse(os.Args[i:]) == nil {
				 break
			 }
		}
	}

}
