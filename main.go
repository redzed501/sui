package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("SUI - home server dashboard")
	r := mux.NewRouter()
	var i IndexData

	serveAssets(r)
	r.HandleFunc("/", serveIndex)

	http.ListenAndServe(":80", r)
}

func parseConfig() {

}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving Index")
	var t = template.Must(template.ParseFiles("./templates/index.html"))

	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func serveAssets(r *mux.Router) {
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}
