package providers

import "time"

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

func newAppProvider() *AppProvider {
	//Specific provider constructor must fill in other details
	return &AppProvider{
		Protected: false,
		Priority:  0,
		Apps: make(map[string]*App),
	}
}

func newApp() *App {
	return &App{
		Icon: "application",
		URL: "youtube.com/watch?v=dQw4w9WgXcQ",
		Protected: false,
		Added: time.Now(),
	}
}
