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
	config.Level = getLevelName(level_DEBUG)
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
	flag.StringVar(&config.Prefix, "prefix", config.Prefix, "set log4g prefix")
	flag.StringVar(&config.Level, "level", config.Level, "set log4g level")
	flag.StringVar(&config.Flag, "flag", config.Flag, "set log4g flag, separated by '|'")
	flag.StringVar(&config.Filename, "filename", config.Filename, "set log4g filename")
	flag.Parse()
}
