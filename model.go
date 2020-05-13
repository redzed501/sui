package main

import (
	"github.com/willfantom/sui/providers"
)

type Config struct {
	Providers map[string]*providers.Provider
}

type ProviderConfig struct {
	Name     string
	Priority uint8
}

type IndexData struct {
	Providers []providers.Provider
	Query     []QueryEngines
}

type QueryEngines struct {
	URL    string
	Name   string
	Prefix string
}
