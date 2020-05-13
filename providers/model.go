package providers

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type ProviderType int

const (
	Docker ProviderType = iota
	Traefik
)

type Provider struct {
	Title     string
	Placement uint8
	Type      ProviderType
	Config    interface{}
	Apps      map[string]*App
}

type App struct {
	Name    string
	Icon    string
	URL     string
	Protect bool
	Added   time.Time
	Updated time.Time
}

type LostApp struct {
	BaseApp *App
	Lost    time.Time
}

func newApp() *App {
	return &App{
		Name:    "Default",
		Icon:    "None",
		URL:     "https://google.com",
		Protect: false,
	}
}

func (p *Provider) FetchApps() error {
	var err error
	switch p.Type {
	case Docker:
		err = fetchDockerApps(p)
		break
	case Traefik:
		err = fetchTraefikApps(p)
		break
	}

	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) AddApp(name string, icon string, URL string, protect bool) {
	app, ok := p.Apps[name]
	if !ok {
		app = newApp()
		app.Added = time.Now()
		p.Apps[name] = app
		log.Debugf("Adding new app (%s) to provider (%s)\n", name, p.Title)
	}
	app.Name = name
	app.URL = URL
	app.Icon = icon
	app.Protect = protect
	app.Updated = time.Now()
}

func (p *Provider) RemoveApp(name string) {
	app, ok := p.Apps[name]
	if ok {
		var lost LostApp
		lost.BaseApp = app
		lost.Lost = time.Now()
		delete(p.Apps, name)
	}
}

func fetchFileApps(p *Provider) error {
	_, valid := p.Config.(DockerProvider)
	if !valid {
		return errors.New("File provider has invalid config")
	}
	return nil
}

func fetchTraefikApps(p *Provider) error {
	_, valid := p.Config.(*DockerProvider)
	if !valid {
		return errors.New("Traefik provider has invalid config")
	}
	return nil
}
