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
