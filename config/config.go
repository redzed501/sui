package config

import (
	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	TraefikURL   map[string]string
	TraefikUser  map[string]string `required:"false"`
	TraefikPass  map[string]string `required:"false"`
	LocalTraefik string            `default:"" required:"false"`

	Ignore []string `required:"false" envconfig:"IGNORE"`
	DEBUG  bool     `default:"false" envconfig:"DEBUG"`
}

type TraefikConfig struct {
	URL       string
	auth      bool
	user      string
	pass      string
	useDocker bool
}

var (
	GlobalConfig   Config
	traefikConfigs map[string]*TraefikConfig
)

func LoadConfig() error {
	log.Debugf("Loading Config from ENV\n")
	err := envconfig.Process("SUI", &GlobalConfig)
	if err != nil {
		return err
	}

	traefikConfigs = make(map[string]*TraefikConfig)
	for name, url := range GlobalConfig.TraefikURL {
		var tc *TraefikConfig
		tc.URL = url
		user, ok := GlobalConfig.TraefikUser[name]
		if ok {
			tc.auth = true
			tc.user = user
			pass, ok := GlobalConfig.TraefikUser[name]
			if ok {
				tc.pass = pass
			} else {
				log.Errorf("Username but no password provided for Traefik Provider: %s", name)
				continue
			}
			if GlobalConfig.LocalTraefik == name {
				tc.useDocker = true
			}
			traefikConfigs[name] = tc
		}
	}
	return nil
}

func getTraefikConfigs() map[string]*TraefikConfig {
	return traefikConfigs
}

func IsTraefikEnabled() bool {
	if GlobalConfig.TraefikURL != nil {
		return false
	}
	return true
}

func IsDebug() bool {
	return GlobalConfig.DEBUG
}
