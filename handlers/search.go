package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
	"strconv"

	"github.com/JojiiOfficial/ZimWiki/zim"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func searchSingle(query string, wiki *zim.File) ([]zim.SRes, string) {
	// Search entries
	entries := wiki.SearchForEntry(query)
	// Sort them by similarity
	sort.Sort(sort.Reverse(zim.ByPercentage(entries)))

	// Limit all results to 100
	e := 100
	if len(entries) < e {
		e = len(entries)
	}

	// Know the number of results
	resultText := "No result was found"

	if len(entries) == 1 {
		resultText = "One result was found"
	} else if len(entries) != 0 {
		resultText = strconv.Itoa(len(entries)) + " results was found"
	}

	return entries[:e], resultText
}

// Search in all available wikis
func searchGlobal(query string, handler *zim.Handler) ([]zim.SRes, string) {
	var results []zim.SRes
	mx := sync.Mutex{}
	wg := sync.WaitGroup{}
	files := handler.GetFiles()

	wg.Add(len(files))

	// Concurrent global search
	for i := range files {
		go func(index int) {
			// Search
			r := files[index].SearchForEntry(query)

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

	// Limit all results to 100
	e := 100
	if len(results) < e {
		e = len(results)
	}

	// Know the number of results
	resultText := "No result was found"

	if len(results) == 1 {
		resultText = "One result was found"
	} else if len(results) != 0 {
		resultText = strconv.Itoa(len(results)) + " results was found"
	}

	return results[:e], resultText
}

// Search handles serach requests
func Search(w http.ResponseWriter, r *http.Request, hd HandlerData) error {
	vars := mux.Vars(r)
	wiki, ok := vars["wiki"]
	if !ok {
		return fmt.Errorf("Missing parameter")
	}
	var query string

	if r.Method == "GET" {
		// Get GET search query
		getQuery, ok := r.URL.Query()["q"]
		if ok && len(getQuery) > 0 {
			query = getQuery[0]
		}
	} else if r.Method == "POST" {
		// Get Post search query
		query = r.PostFormValue("sQuery")
	}

	if len(query) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return nil
	}

	var res []zim.SRes
	var source string
	var resultText string

	start := time.Now()
	if wiki == "-" {
		source = "global search"

		res, resultText = searchGlobal(query, hd.ZimService)
	} else {
		// Wiki search
		z := hd.ZimService.FindWikiFile(wiki)
		if z == nil {
			http.NotFound(w, r)
			return ErrNotFound
		}
		source = z.File.Title()

		// Search for query in WIKI
		res, resultText = searchSingle(query, z)
	}

	favCache := make(map[string]string)
	results := make([]SearchResult, len(res))

	for i := range res {
		// Get Favicon of each file
		fav, has := favCache[res[i].GetID()]
		if !has {
			favIcon, err := res[i].Favicon()
			if err == nil {
				fav = zim.GetRawWikiURL(res[i].File, favIcon)
			} else {
				fav = ""
			}

			favCache[res[i].GetID()] = fav
		}

		results[i] = SearchResult{
			Img:   fav,
			Link:  zim.GetWikiURL(res[i].File, *res[i].DirectoryEntry),
			Title: string(res[i].Title()),
			Wiki:  wiki,
		}
	}

	timeTook := time.Since(start)

	log.Info(resultText, " in ", timeTook)

	// Redirect to wiki page if only
	// one search result was found
	if len(results) == 1 {
		http.Redirect(w, r, results[0].Link, http.StatusMovedPermanently)
		return nil
	}

	return serveTemplate(SearchTemplate, w, TemplateData{
		Wiki: wiki,
		SearchTemplateData: SearchTemplateData{
			Results:      results,
			QueryText:    query,
			SearchSource: source,
			ResultText:   resultText,
			TimeTook:     timeTook,
		},
	})
}
