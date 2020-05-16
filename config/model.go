package config

import (
	"github.com/willfantom/sui/bookmarks"
	"github.com/willfantom/sui/search"
)

type DockerType uint8

const (
	Socket DockerType = iota
	TCP
)

type Config struct {
	AppTitle        string                           `json:"title"`
	Debug           bool                             `json:"debug"`
	AppRefresh      int                              `json:"app_refresh"`
	ProviderConfigs []AppProviderConfig              `json:"appproviders"`
	SearchEngines   map[string]*search.SearchEngine  `json:"engines"`
	BookmarkCats    map[string]*[]bookmarks.Bookmark `json:"bookmarks"`
}

type AppProviderConfig struct {
	Name  string `json:"name"`
	PType string `json:"type"`
}

func NewConfig() *Config {
	return &Config{
		AppTitle:       "SUI Dashboard",
		AppRefresh:     300,
		SearchEngines:  make(map[string]*search.SearchEngine),
		BookmarkCats:   make(map[string]*[]bookmarks.Bookmark),
	}
}
