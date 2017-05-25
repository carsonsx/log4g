package log4g

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"flag"
)

const config_file_path = "log4g.json"

var Config struct{
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	Flag     string `json:"flag"`
	Filename string `json:"filename"`
}

func loadConfig() {

	// default
	Config.Level = GetLevelName(DEBUG)
	Config.Flag = "date|microseconds|shortfile"

	// load form config file
	if _, err := os.Stat(config_file_path); err == nil {
		data, err := ioutil.ReadFile(config_file_path)
		if err != nil {
			log.Print(err)
			return
		}
		err = json.Unmarshal(data, &Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	// override from os arguments
	flag.StringVar(&Config.Prefix, "prefix", Config.Prefix, "set log4g prefix")
	flag.StringVar(&Config.Level, "level", Config.Level, "set log4g level")
	flag.StringVar(&Config.Flag, "flag", Config.Flag, "set log4g flag, separated by '|'")
	flag.StringVar(&Config.Filename, "filename", Config.Filename, "set log4g filename")
	flag.Parse()
}
