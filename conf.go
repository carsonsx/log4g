package log4g

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
)

const config_file_path = "log4g.json"

var Config struct{
	Level    string `json:"level"`
	Filename string `json:"filename"`
	Flag     []string `json:"flag"`
}

func loadConfig() {

	Config.Level = GetLevelName(DEBUG)
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

}
