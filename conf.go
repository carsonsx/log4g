package log4g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

var (
	gEnv                  string
	gFile                 string
	defaultConfigFilepath = []string{"./log4g.json", "conf/log4g.json", "config/log4g.json"}
)

type loggerConfig struct {
	Disabled  bool   `json:"disabled"`
	Prefix    string `json:"prefix"`
	Level     string `json:"level"`
	Flag      string `json:"flag"`
	Output    string `json:"output"`
	Filename  string `json:"filename"`
	Maxsize   int64  `json:"maxsize"`
	MaxLines  int    `json:"max_lines"`
	MaxCount  int    `json:"max_count"`
	Daily     bool   `json:"daily"`
	Address   string `json:"address"`
	DB        int    `json:"db"`
	Password  string `json:"password"`
	RedisType string `json:"redis_type"`
	RedisKey  string `json:"redis_key"`
	Network   string `json:"network"`
	Codec     string `json:"codec"`
	JsonKey   string `json:"json_key"`
	JsonExt   string `json:"json_ext"`
}

func NewConfig() *Config {
	c := new(Config)
	c.initDefault()
	return c
}

type Config struct {
	Prefix  string          `json:"prefix"`
	Level   string          `json:"level"`
	Flag    string          `json:"flag"`
	Loggers []*loggerConfig `json:"Loggers"`
}

func (c *Config) initDefault()  {
	// default
	c.Level = LEVEL_DEBUG.Name()
	c.Flag = "date|time|shortfile"
}

func setEnv(env string) {
	gEnv = env
	//loadDefaultConfig()
}

func loadConfig(filepath string, mapping interface{}) error {

	config := mapping.(*Config)

	// load form Config file
	_, err := os.Stat(filepath)
	if err == nil { //file exist
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, config)
		if err != nil {
			return err
		}
	}
	return err
}

func parseFlag(strFlag string) int {
	flags := strings.Split(strFlag, "|")

	flag := 0
	for _, name := range flags {
		flag = flag | getFlagByName(name)
	}
	return flag
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


