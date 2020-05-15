package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

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
	addBookmarks()

	go refreshApps()

	r := mux.NewRouter()

	serveAssets(r)
	r.HandleFunc("/", serveIndex)
	r.HandleFunc("/js/search.js", serveSearchJS)

	http.ListenAndServe(":6999", r)
}

func addAppProviders() {
	var err error
	for _, provider := range config.GetAppProviderConfigs() {
		log.Debugf("adding provider | %s", provider.Name)
		indexData.AppProviders[provider.Name], err = providers.NewAppProvider(provider.Name, provider.PType)
		if err != nil {
			panic(err)
		}
	}

}

func addSearchEngines() {
	log.Debugf("adding search engines\n")
	indexData.SearchEngines = config.GetSearchEngines()
}

func addBookmarks() {
	log.Debugf("adding bookmarks\n")
	indexData.Bookmarks = config.GetBookmarks()
	for name, bmks := range indexData.Bookmarks {
		log.Debugf("added bookmarks | category: %s | count: %d", name, len(*bmks))
	}
}


func refreshApps() {
	for true {
		for name, prov := range indexData.AppProviders {
			err := prov.RefreshApps()
			if err != nil {
				panic(err)
			}
			log.Debugf("found apps | provider: %s | app count: %d", name, len(prov.Apps))
		}
		time.Sleep(time.Duration(config.GetAppRefresh())*time.Second)
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

func serveSearchJS(w http.ResponseWriter, r *http.Request) {
	var t = template.Must(template.ParseFiles("./templates/search.js"))

	err := t.Execute(w, indexData)
	if err != nil {
		panic(err)
	}
}

func serveAssets(r *mux.Router) {
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}

