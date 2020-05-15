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

func (ap *AppProvider) renameApp(oldName, newName string) error {
	if !ap.appExist(oldName) {
		return fmt.Errorf("could not delete app")
	}
	ap.Apps[newName] = ap.Apps[oldName]
	return ap.deleteApp(oldName)
}

func (ap *AppProvider) deleteApp(appName string) error {
	if !ap.appExist(appName) {
		return fmt.Errorf("could not delete app")
	}
	delete(ap.Apps, appName)
	return nil
}

func (ap *AppProvider) appExist(appName string) bool {
	_, exist := ap.Apps[appName]
	return exist
}

func getDefaultIcon(name string) string {
	_, exist := iconDefault[name]
	if exist {
		return iconDefault[name]
	}
	return "application"
}