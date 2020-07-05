package handlers

import (
	"fmt"
	"net/http"

	"git.jojii.de/jojii/zimserver/zim"
	"github.com/gorilla/mux"
	gzim "github.com/tim-st/go-zim"
)

func searchWiki(query string, wiki *zim.File) ([]SearchResult, error) {
	var results []SearchResult

	// Find
	dirs := wiki.EntriesWithSimilarity(gzim.NamespaceArticles, []byte(query), 100)
	results = make([]SearchResult, len(dirs))

	for i := range dirs {
		results[i] = SearchResult{
			Link:  zim.GetWikiURL(wiki, dirs[i]),
			Title: string(dirs[i].Title()),
			Wiki:  wiki.GetID(),
		}
	}

	return results, nil
}

// Search handles serach requests
func Search(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	vars := mux.Vars(r)
	wiki, ok := vars["wiki"]
	if !ok {
		return fmt.Errorf("Missing parameter")
	}

	// Get search query
	query, ok := r.URL.Query()["q"]
	if !ok || len(query) == 0 {
		http.Redirect(w, r, "/", http.StatusUnprocessableEntity)
		return nil
	}

	var results []SearchResult
	var err error

	if wiki == "-" {
		// Global search
		for i := range hd.ZimService.GetFiles() {
			r, err := searchWiki(query[0], &hd.ZimService.GetFiles()[i])
			if err != nil {
				return err
			}

			results = append(results, r...)
		}

	} else {
		// Wiki search
		z := hd.ZimService.FindWikiFile(wiki)
		if z == nil {
			http.NotFound(w, r)
			return ErrNotFound
		}

		// Search for query in WIKI
		results, err = searchWiki(query[0], z)
		if err != nil {
			return err
		}
	}

	return serveTemplate(SearchTemplate, w, TemplateData{
		Wiki: wiki,
		SearchTemplateData: SearchTemplateData{
			Results:   results,
			QueryText: query[0],
		},
	})
}
