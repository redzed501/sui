package config

import (
	"os"
	"strings"

	"gopkg.in/square/go-jose.v2/json"

	log "github.com/sirupsen/logrus"
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

/*

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
	Auth      bool
	User      string
	Pass      string
	UseDocker bool
}

const (
	URLPrefix string = "https://"
)

var (
	GlobalConfig   Config
	traefikConfigs map[string]*TraefikConfig = make(map[string]*TraefikConfig)
)


func


func LoadConfig() error {
	log.Debugf("Loading Config from ENV\n")
	err := envconfig.Process("SUI", &GlobalConfig)
	if err != nil {
		return err
	}

	fmt.Println("map:", GlobalConfig.TraefikURL)
	for name, url := range GlobalConfig.TraefikURL {
		log.Infof("Adding Provider: %s", name)
		var tc TraefikConfig
		tc.URL = URLPrefix + url
		user, ok := GlobalConfig.TraefikUser[name]
		if ok {
			tc.Auth = true
			tc.User = user
			pass, ok := GlobalConfig.TraefikPass[name]
			if ok {
				tc.Pass = pass
			} else {
				log.Errorf("Username but no password provided for Traefik Provider: %s", name)
				continue
			}
			if GlobalConfig.LocalTraefik == name {
				tc.UseDocker = true
			}
		}
		traefikConfigs[name] = &tc
	}
	return nil
}

func LoadConfigTest() error {
	log.Debugf("Loading Config from ENV\n")
	err := envconfig.Process("SUI", &GlobalConfig)
	if err != nil {
		return err
	}

	log.Infof("Adding Provider: TEST")
	var tc TraefikConfig
	tc.URL = URLPrefix + "traefik.fantom.host"
	tc.User = "fantom"
	tc.Pass = "0FK.yhip17mB>}oJ"
	tc.UseDocker = true
	traefikConfigs["TEST"] = &tc
	return nil
}

func GetTraefikConfigs() map[string]*TraefikConfig {
	return traefikConfigs
}

func IsTraefikEnabled() bool {
	if GlobalConfig.TraefikURL != nil {
		return true
	}
	return false
}

func GetProviderCount() int {
	return len(traefikConfigs)
}

func IsDebug() bool {
	return GlobalConfig.DEBUG
}*/
