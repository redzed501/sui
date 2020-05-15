package providers

import (
	"strings"

	"github.com/willfantom/sui/config"
)

const (
	labelRoot string = "sui"
)

func NewAppProvider(name string, ptype string) (*AppProvider, error) {
	var err error
	provider := &AppProvider{
		Protected: false,
		Priority:  0,
		Apps: make(map[string]*App),
		PType: ptype,
	}

	switch strings.ToLower(ptype) {
	case "docker":
		provider.TypeConfig, err = newDocker(name)
		break
	case "traefik":
		break
	}

	return provider, err
}

func (ap *AppProvider) RefreshApps() error {
	switch ap.PType {
	case "docker":
		dkrCnf, err := toDocker(ap.TypeConfig)
		if err != nil {
			return err
		}
		ap.Apps = dkrCnf.GetApps()
		break
	}
	return nil
}

func getFileConfigRoot() string {
	return config.GetFileConfigRoot()
}