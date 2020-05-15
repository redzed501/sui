package providers

import (
	"time"
)

var enabledProviders = [...]string{"docker", "traefik"}

type AppProvider struct {
	PType      string
	TypeConfig interface{}
	Protected  bool
	Priority   uint8
	Apps       map[string]*App
}

type App struct {
	Icon    string
	URL     string
	Enabled bool
	Added   time.Time
}

// func NewAppProvider(pType ProviderType, cnf interface{}) (*AppProvider, error) {
// 	log.Debugf("creating new provider")
// 	err := fmt.Errorf("could not create provider")
// 	var provider interface{}
// 	switch pType {
// 	case Docker:
// 		provider, err = NewDockerProvider(cnf.(*config.DockerConfig))
// 		break
// 	case Traefik:
// 		provider, err = NewTraefikProvider(cnf.(*config.TraefikConfig))
// 		break
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &AppProvider{
// 		Protected: false,
// 		Priority:  0,
// 		Apps: make(map[string]*App),
// 		PType: pType,
// 		TypeConfig: provider,
// 	}, nil
// }

func newApp() *App {
	return &App{
		Icon:    "application",
		URL:     "youtube.com/watch?v=dQw4w9WgXcQ",
		Enabled: true,
		Added:   time.Now(),
	}
}
