package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"git.jojii.de/jojii/zimserver/zim"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func searchSingle(query string, wiki *zim.File) []zim.SRes {
	// Search entries
	entries := wiki.SearchForEntry(query, 100)
	// Sort them by similarity
	sort.Sort(sort.Reverse(zim.ByPercentage(entries)))

	return entries
}

// Search in all available wikis
func searchGlobal(query string, handler *zim.Handler) []zim.SRes {
	var results []zim.SRes
	mx := sync.Mutex{}
	wg := sync.WaitGroup{}
	files := handler.GetFiles()

	wg.Add(len(files))

	// Concurrent global search
	for i := range files {
		go func(index int) {
			// Search
			r := files[index].SearchForEntry(query, 100)

			// Append results
			mx.Lock()
			results = append(results, r...)
			mx.Unlock()

			wg.Done()
		}(i)
	}

	wg.Wait()

	// Sort by similarity
	sort.Sort(sort.Reverse(zim.ByPercentage(results)))

	e := 100
	if len(results) < e {
		e = len(results)
	}

	return results[:e]
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

	var res []zim.SRes
	var source string

	start := time.Now()
	if wiki == "-" {
		source = "global search"

		res = searchGlobal(query[0], hd.ZimService)
	} else {
		// Wiki search
		z := hd.ZimService.FindWikiFile(wiki)
		if z == nil {
			http.NotFound(w, r)
			return ErrNotFound
		}
		source = z.File.Title()

		// Search for query in WIKI
		res = searchSingle(query[0], z)
	}

	results := make([]SearchResult, len(res))
	for i := range res {
		results[i] = SearchResult{
			Link:  zim.GetWikiURL(res[i].File, *res[i].DirectoryEntry),
			Title: string(res[i].Title()),
			Wiki:  wiki,
		}
	}

	log.Info("Searching took ", time.Since(start))

	return serveTemplate(SearchTemplate, w, TemplateData{
		Wiki: wiki,
		SearchTemplateData: SearchTemplateData{
			Results:      results,
			QueryText:    query[0],
			SearchSource: source,
		},
	})
}
