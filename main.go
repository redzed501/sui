package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/willfantom/sui/config"
	"github.com/willfantom/sui/providers"
)

var (
	indexData *IndexData
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.Infof("SUI - home server dashboard")

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("problem loading config")
	}
	if config.IsDebug() {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("parsed Config\n")

	indexData = NewIndexData()

	addAppProviders()
	addSearchEngines()

	refreshApps()

	r := mux.NewRouter()

	serveAssets(r)
	r.HandleFunc("/", serveIndex)

	http.ListenAndServe(":80", r)
}

func addAppProviders() {
	var err error
	log.Debugf("adding docker providers\n")
	for name, cnf := range config.GetDockerCnfs() {
		indexData.AppProviders[name], err = providers.NewDockerProvider(cnf)
		if err != nil {
			panic(err)
		}
	}

	log.Debugf("adding docker x træfik providers\n")
	for name, cnf := range config.GetTraefikCnfs() {
		if cnf.PariedDocker != "" {
			dapp, exist := indexData.AppProviders[name]
			if exist {
				indexData.AppProviders[name], err = providers.NewDockerTraefikProvider(cnf, dapp)
			}
		}
	}

	log.Debugf("adding træfik providers\n")
	for name, cnf := range config.GetTraefikCnfs() {
		indexData.AppProviders[name], err = providers.NewTraefikProvider(cnf)
	}
	// Load other providers here
	
}

func addSearchEngines() {
	log.Debugf("adding search engines\n")
	indexData.SearchEngines = config.GetSearchEngines()
}


func refreshApps() {
	for name, prov := range indexData.AppProviders {
		err := prov.RefreshApps()
		if err != nil {
			panic(err)
		}
		log.Debugf("found apps | provider: %s | app count: %d", name, len(prov.Apps))
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving Index")
	var t = template.Must(template.ParseFiles("./templates/index.html"))

	err := t.Execute(w, indexData)
	if err != nil {
		panic(err)
	}
}

func serveAssets(r *mux.Router) {
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}

