package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/willfantom/sui/providers"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Infof("SUI - home server dashboard")
	r := mux.NewRouter()

	// Add Test Docker Provider for Testing
	provider, err := providers.NewDockerProvider("test", 1, "/var/run/docker.sock")
	if err != nil {
		log.Panicf("Connection to provider failed\n")
	}
	provider.FetchApps()

	serveAssets(r)
	r.HandleFunc("/", serveIndex)

	http.ListenAndServe(":80", r)
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
