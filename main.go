package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/willfantom/sui/config"
	"github.com/willfantom/sui/providers"
)

const (
	dockerSock string = "/var/run/docker.sock"
)

var (
	provs map[string]*providers.Provider
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Infof("SUI - home server dashboard")

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Problem loading config")
	}

	provs = make(map[string]*providers.Provider)
	if config.IsDockerEnabled() {
		loadDockerProvider()
	}

	refreshApps()

	r := mux.NewRouter()

	serveAssets(r)
	r.HandleFunc("/", serveIndex)

	http.ListenAndServe(":80", r)
}

func loadDockerProvider() {
	provider, err := providers.NewDockerProvider(0, dockerSock)
	if err != nil {
		log.Fatalf("Could not connect to docker")
	}
	provs["docker"] = provider
}

func refreshApps() {
	for _, prov := range provs {
		err := prov.FetchApps()
		if err != nil {
			panic(err)
		}
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving Index")
	var t = template.Must(template.ParseFiles("./templates/index.html"))

	err := t.Execute(w, IndexData{})
	if err != nil {
		panic(err)
	}
}

func serveAssets(r *mux.Router) {
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}
