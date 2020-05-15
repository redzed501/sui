package config

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/willfantom/sui/bookmarks"
	"github.com/willfantom/sui/search"
)

var (
	cnf *Config
)

const (
	fileConfigRoot string = "/sui"
	mainConfigFile string = "config.json"
)

func LoadConfig() error {
	cnf = NewConfig()
	err := loadFromFile(fmt.Sprintf("%s/%s", fileConfigRoot, mainConfigFile))
	//TODO: ENV, CMDARGS
	return err
}

func loadFromFile(path string) error {
	cnfFile, err := os.Open(path)
	if err != nil {
		log.Errorf("no config file found at %s", path)
		return err
	}
	defer cnfFile.Close()
	err = json.NewDecoder(cnfFile).Decode(cnf)
	if err != nil {
		log.Errorf("config file could not be parsed | %s", path)
		return err
	}
	return nil
}

func IsDebug() bool {
	return cnf.Debug
}

func GetAppProviderConfigs() []AppProviderConfig {
	return cnf.ProviderConfigs
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

func GetFileConfigRoot() string {
	return fileConfigRoot
}
