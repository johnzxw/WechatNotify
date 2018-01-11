package Tool

import (
	"flag"
	"os"
	"path/filepath"

	"runtime"
	"strings"
	"github.com/larspensjo/config"
)

var currentPath = ""
var configFactory = make(map[string]map[string]string)

func GetCurrentDirectory() string {
	if currentPath != "" {
		return currentPath
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("currentpath error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func GetConfig(flagString string) map[string]string {
	//加缓存
	if data, err := configFactory[flagString]; err {
		return data
	}
	configFile := GetCurrentDirectory() + "/" + "config.ini"
	var TOPIC = make(map[string]string)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	//cfg, err := config.ReadDefault(*configFile)
	cfg, err := config.ReadDefault(configFile)
	if err != nil {
		//panic("Fail to find " + *configFile)
		panic("Fail to find " + configFile)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection(flagString) {
		section, err := cfg.SectionOptions(flagString)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(flagString, v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	configFactory[flagString] = make(map[string]string)
	configFactory[flagString] = TOPIC
	//Initialized topic from the configuration END
	return TOPIC
}
