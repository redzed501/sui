package main

import (
	"github.com/willfantom/sui/bookmarks"
	"github.com/willfantom/sui/providers"
	"github.com/willfantom/sui/search"
)

// IndexData is the data to be used in the template file
type IndexData struct {
	AppTitle      string
	AppProviders  map[string]*providers.AppProvider
	Bookmarks     map[string]*[]bookmarks.Bookmark
	SearchEngines map[string]*search.SearchEngine
}

// NewIndexData creates a new IndexData...
// Creates the required maps
// Returns a pointer to the the created indexdata
func NewIndexData() *IndexData {
	return &IndexData{
		AppProviders:  make(map[string]*providers.AppProvider),
		Bookmarks:     make(map[string]*[]bookmarks.Bookmark),
		SearchEngines: make(map[string]*search.SearchEngine),
	}
}
