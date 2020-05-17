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
	DisplayName string
	Icon        string
	URL         string
	DisplayURL  string
	Enabled     bool
	Added       time.Time
}
