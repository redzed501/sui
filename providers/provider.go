package providers

import (
	"strings"
	"time"

	"github.com/willfantom/sui/config"
)

const (
	labelRoot string = "sui"

	maxDisplayNameLen int = 16
	maxDisplayURLLen int = 19
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
		provider.TypeConfig, err = newTraefik(name)
		break
	}

	return provider, err
}

func newApp() *App {
	return &App{
		DisplayName: "",
		Icon:    "application",
		DisplayURL: "no url",
		URL:     "/",
		Enabled: true,
		Added:   time.Now(),
	}
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
	case "traefik":
		trCnf, err := toTraefik(ap.TypeConfig)
		if err != nil {
			return err
		}
		ap.Apps = trCnf.GetApps()
		break
	}
	return nil
}

func (app *App) FormatApp(name string) {
	if len(name) > maxDisplayNameLen {
		app.DisplayName = name[:(maxDisplayNameLen-3)] + "..."
	} else {
		app.DisplayName = name
	}
	if len(app.URL) > maxDisplayURLLen {
		app.DisplayURL = app.URL[:(maxDisplayURLLen-3)] + "..."
	} else {
		app.DisplayURL = app.URL
	}
}

func getDefaultIcon(name string) string {
	_, exist := iconDefault[name]
	if exist {
		return iconDefault[name]
	}
	return "application"
}

func getFileConfigRoot() string {
	return config.GetFileConfigRoot()
}