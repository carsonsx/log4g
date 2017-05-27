package log4g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

var (
	gEnv         string
	gFile        string
	filePrefixes = []string{"./", "conf/", "config/"}
)

type loggerConfig struct {
	Disabled bool   `json:"disabled"`
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	Flag     string `json:"flag"`
	Output   string `json:"output"`
	Filename string `json:"filename"`
	Maxsize  int64  `json:"maxsize"`
	Maxlines int    `json:"maxlines"`
	MaxCount int    `json:"maxcount"`
	Daily    bool   `json:"daily"`
}

var Config struct {
	Prefix  string          `json:"prefix"`
	Level   string          `json:"level"`
	Flag    string          `json:"flag"`
	Loggers []*loggerConfig `json:"Loggers"`
}

func setEnv(env string) {
	gEnv = env
	loadDefaultConfig()
}

func loadDefaultConfig() {
	found := false
	basename := "log4g"
	if gEnv != "" {
		basename = "log4g-" + gEnv
	}
	basename += ".json"
	for _, prefix := range filePrefixes {
		filepath := prefix  + basename
		if err := loadConfig(filepath); err == nil {
			found = true
			break
		}
	}
	if !found {
		//log.Printf("not found any %s config file in [%s] , use default stdout", basename, strings.Join(filePrefixes, ","))
		initLoggers()
	}
}

func loadConfig(filename string) error {
	if gFile == filename {
		return nil
	}
	gFile = filename
	return reloadConfig()
}

func reloadConfig() error {

	// default
	Config.Level = LEVEL_DEBUG.Name()
	Config.Flag = "date|time|shortfile"

	// load form Config file
	_, err := os.Stat(gFile)
	if err == nil { //file exist
		data, err := ioutil.ReadFile(gFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &Config)
		if err != nil {
			panic(err)
		}
		initLoggers()
		//initialized = true
	}
	return err
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
