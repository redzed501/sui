package config

import (
	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	DockerEnabled bool   `default:true`
	TraefikURL    string `default:""`
}

var (
	GlobalConfig Config
)

func LoadConfig() error {
	log.Debugf("Loading Config from ENV\n")
	err := envconfig.Process("SUI", &GlobalConfig)
	return err
}

func IsDockerEnabled() bool {
	return GlobalConfig.DockerEnabled
}

func IsTraefikEnabled() bool {
	if GlobalConfig.TraefikURL != "" {
		return true
	}
	return false
}

func GetTraefikURL() string {
	return GlobalConfig.TraefikURL
}
