package config

import (
	"encoding/json"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/willfantom/sui/bookmarks"
	"github.com/willfantom/sui/search"
)

var (
	cnf *Config
)

const (
	cnfFilePath string = "/sui/config.json"
)

func LoadConfig() error {
	cnf = NewConfig()

	//File
	loadFromFile(cnfFilePath)
	//ENV
	// ARGS

	parseConfig()
	return nil
}

func loadFromFile(path string) {
	cnfFile, err := os.Open(path)
	if err != nil {
		log.Errorf("no config file found at %s", path)
		return
	}
	defer cnfFile.Close()
	err = json.NewDecoder(cnfFile).Decode(cnf)
	if err != nil {
		log.Errorf("config file could not be parsed | %s", path)
		return
	}
}

func parseConfig() {
	for name, c := range cnf.DockerConfigs {
		if c.CnfDType == "socket" {
			cnf.DockerConfigs[name].DType = Socket
		} else if c.CnfDType == "tcp" {
			cnf.DockerConfigs[name].DType = TCP
		} else {
			log.Errorf("invaid docker type for %s\n", name)
			delete(cnf.DockerConfigs, name)
		}
	}
	for name, c := range cnf.TraefikConfigs {
		if c.CnfIgnoredList != "" {
			cnf.TraefikConfigs[name].IgnoredList = strings.Split(strings.ToLower(c.CnfIgnoredList), " ")
		}
	}
}

func IsDebug() bool {
	return cnf.Debug
}

func GetDockerCnfs() map[string]*DockerConfig {
	return cnf.DockerConfigs
}

func GetTraefikCnfs() map[string]*TraefikConfig {
	return cnf.TraefikConfigs
}

func GetSearchEngines() map[string]*search.SearchEngine {
	return cnf.SearchEngines
}

func GetBookmarks() map[string]*[]bookmarks.Bookmark {
	return cnf.BookmarkCats
}

func GetAppRefresh() int {
	return cnf.AppRefresh
}
