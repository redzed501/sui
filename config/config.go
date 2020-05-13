package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Docker      bool `default:"true"`
	TraefikURL  map[string]string
	TraefikUser string `required:"false"`
	TraefikPass string `required:"false"`
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
	return GlobalConfig.Docker
}

func IsTraefikEnabled() bool {
	if GlobalConfig.TraefikURL != nil {
		return false
	}
	return true
}

func GetTraefikURL(name string) (string, error) {
	value, ok := GlobalConfig.TraefikURL[name]
	if !ok {
		return "", fmt.Errorf("No traefik URL for instance %s", name)
	}
	return value, nil
}
