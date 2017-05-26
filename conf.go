package log4g

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"strings"
)

const config_file_path = "log4g.json"


type loggerConfig struct{
	Disabled bool `json:"disabled"`
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	Flag     string `json:"flag"`
	Output   string `json:"output"`
	Filename string `json:"filename"`
	Maxsize  int64 `json:"maxsize"`
	Maxlines int `json:"maxlines"`
	MaxCount int `json:"maxcount"`
	Daily    bool `json:"daily"`
}

var Config struct{
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	Flag     string `json:"flag"`
	Loggers []*loggerConfig `json:"Loggers"`
}

func loadConfig() {

	// default
	Config.Level = LEVEL_DEBUG.Name()
	Config.Flag = "date|time|shortfile"

	// load form Config file
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
	//if len(os.Args) > 1 {
	//	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	//	fs.StringVar(&Config.Prefix, "log4g-prefix", Config.Prefix, "set log4g prefix")
	//	fs.StringVar(&Config.Level, "log4g-level", Config.Level, "set log4g level")
	//	fs.StringVar(&Config.Flag, "log4g-flag", Config.Flag, "set log4g flag, separated by '|'")
	//	fs.StringVar(&Config.Filename, "log4g-filename", Config.Filename, "set log4g filename")
	//	for i := 1; i < len(os.Args); i++ {
	//		 if fs.Parse(os.Args[i:]) == nil {
	//			 break
	//		 }
	//	}
	//}

}

func getFlagByName(name string) int {
	flags := make(map[string]int)
	flags["date"] = ldate
	flags["time"] = ltime
	flags["microseconds"] = lmicroseconds
	flags["longfile"] = llongfile
	flags["shortfile"] = lshortfile
	flags["UTC"] = lutc
	flags["stdFlags"] = lstdFlags
	return flags[name]
}

func parseFlag(strFlag string, defaultValue int) int {
	if strFlag == "" {
		return defaultValue
	}
	flags := strings.Split(strFlag, "|")

	flag := 0
	for _, name := range flags {
		flag = flag | getFlagByName(name)
	}
	return flag
}
