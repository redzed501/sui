package providers

import (
	"fmt"
	"time"

	"github.com/prometheus/common/log"
	"github.com/willfantom/sui/config"
)

type ProviderType int

const (
	Docker ProviderType = iota
	Traefik
)

type AppProvider struct {
	PType      ProviderType
	TypeConfig interface{}
	Protected  bool
	Priority   uint8
	Apps       map[string]*App
}

type App struct {
	Icon      string
	URL       string
	Protected bool
	Added     time.Time
}

func NewAppProvider(pType ProviderType, cnf interface{}) (*AppProvider, error) {
	log.Debugf("creating new provider")
	err := fmt.Errorf("could not create provider")
	var provider interface{}
	switch pType {
	case Docker:
		provider, err = NewDockerProvider(cnf.(*config.DockerConfig))
		break
	case Traefik:
		provider, err = NewTraefikProvider(cnf.(*config.TraefikConfig))
		break
	}
	if err != nil {
		return nil, err
	}
	return &AppProvider{
		Protected: false,
		Priority:  0,
		Apps: make(map[string]*App),
		PType: pType,
		TypeConfig: provider,
	}, nil
}

func newApp() *App {
	return &App{
		Icon: "application",
		URL: "youtube.com/watch?v=dQw4w9WgXcQ",
		Protected: false,
		Added: time.Now(),
	}
}
