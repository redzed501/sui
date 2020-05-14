package providers

import (
	"fmt"
)

func (ap *AppProvider) RefreshApps() error {
	switch ap.PType {
	case Docker:
		dp, valid := ap.TypeConfig.(*DockerProvider)
		if !valid {
			return fmt.Errorf("invalid config found!!")
		}
		return dp.GetApps(ap.Apps)
	case Traefik:
		tp, valid := ap.TypeConfig.(*TraefikProvider)
		if !valid {
			return fmt.Errorf("invalid config found!!")
		}
		return tp.GetApps(ap.Apps)
	}
	
	return nil
}

func getDefaultIcon(name string) string {
	_, exist := iconDefault[name]
	if exist {
		return iconDefault[name]
	}
	return "application"
}