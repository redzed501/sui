package config

import (
	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	DockerEnabled bool `default:"true" envconfig:"DOCKER"`
	TraefikURL    map[string]string
	TraefikUser   string   `required:"false"`
	TraefikPass   string   `required:"false"`
	Ignore        []string `required:"false" envconfig:"IGNORE"`
	DEBUG         bool     `default:"false" envconfig:"DEBUG"`
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
	if GlobalConfig.TraefikURL != nil {
		return false
	}
	return true
}

func GetTraefikURLS() map[string]string {
	return GlobalConfig.TraefikURL
}

func GetProviderCount() int {
	count := 0
	if IsDockerEnabled() {
		count++
	}
	count += len(GetTraefikURLS())
	return count
}

func IsDebug() bool {
	return GlobalConfig.DEBUG
}
