package config

import "github.com/willfantom/sui/search"

type DockerType uint8

const (
	Socket DockerType = iota
	TCP
)

type Config struct {
	AppTitle       string                    `json:"title"`
	Debug          bool                      `json:"debug"`
	AppRefresh     int                       `json:"app_refresh"`
	DockerConfigs  map[string]*DockerConfig  `json:"dockers"`
	TraefikConfigs map[string]*TraefikConfig `json:"traefiks"`
	SearchEngines  map[string]*search.SearchEngine `json:"engines"`
}

type DockerConfig struct {
	Path     string `json:"path"`
	CnfDType string `json:"type"`
	DType    DockerType
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

type TraefikConfig struct {
	URL            string `json:"url"`
	PariedDocker   string `json:"docker"`
	CnfIgnoredList string `json:"ignored"`
	IgnoredList    []string
	User           string `json:"user"`
	Pass           string `json:"pass"`
}

func NewConfig() *Config {
	return &Config{
		AppTitle:       "SUI",
		AppRefresh:     300,
		DockerConfigs:  make(map[string]*DockerConfig),
		TraefikConfigs: make(map[string]*TraefikConfig),
	}
}
