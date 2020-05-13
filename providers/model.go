package providers

import (
	"errors"
	"time"
)

type ProviderType int

const (
	Docker ProviderType = iota
	File
	Traefik
)

type Provider struct {
	Title     string
	Placement uint8
	Type      ProviderType
	Config    interface{}
	Apps      map[string]App
}

type App struct {
	Visible bool
	Icon    string
	URL     string
	Protect bool
	Added   time.Time
}

type LostApp struct {
	BaseApp *App
	Lost    time.Time
}

func (p *Provider) FetchApps() error {
	var err error
	switch p.Type {
	case Docker:
		err = fetchDockerApps(p)
		break
	case File:
		err = fetchFileApps(p)
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
