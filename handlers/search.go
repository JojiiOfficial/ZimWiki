package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.jojii.de/jojii/zimserver/zim"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	gzim "github.com/tim-st/go-zim"
)

func searchWiki(query string, wiki *zim.File) ([]SearchResult, error) {
	var results []SearchResult

	var dirs []gzim.DirectoryEntry

	// Search
	wiki.Mx.Lock()
	dirs = append(dirs, wiki.EntriesWithSimilarity(gzim.NamespaceArticles, []byte(query), 100)...)
	dirs = append(dirs, wiki.EntriesWithTitlePrefix(gzim.NamespaceArticles, []byte(query), 100)...)
	wiki.Mx.Unlock()

	results = make([]SearchResult, len(dirs))

	var err error
	for i := range dirs {
		// Follow redirects
		wiki.Mx.Lock()
		if dirs[i].IsRedirect() {
			dirs[i], err = wiki.FollowRedirect(&dirs[i])
			if err != nil {
				wiki.Mx.Unlock()
				return nil, err
			}
		}
		wiki.Mx.Unlock()

		// Add result
		results[i] = SearchResult{
			Link:  zim.GetWikiURL(wiki, dirs[i]),
			Title: string(dirs[i].Title()),
			Wiki:  wiki.GetID(),
		}
	}

	return removeDoubles(results), nil
}

// Remove duplicate search results
func removeDoubles(inp []SearchResult) []SearchResult {
	keys := make(map[string]bool)
	newArr := make([]SearchResult, 0)

	for i := range inp {
		k := inp[i].Link + inp[i].Wiki

		if _, val := keys[k]; !val {
			keys[k] = true
			newArr = append(newArr, inp[i])
		}
	}

	return newArr
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
	var source string

	start := time.Now()
	if wiki == "-" {
		source = "global search"

		mx := sync.Mutex{}
		wg := sync.WaitGroup{}
		files := hd.ZimService.GetFiles()
		resChan := make(chan error, len(files))

		wg.Add(len(files))

		// Concurrent global search
		for i := range files {
			go func(index int) {
				// Search
				r, err := searchWiki(query[0], &files[index])
				if err != nil {
					resChan <- err
				} else {
					mx.Lock()
					results = append(results, r...)
					mx.Unlock()
					resChan <- nil
				}

				wg.Done()
			}(i)
		}

		wg.Wait()
		close(resChan)

		// Wait for queries
		for i := 0; i < len(files); i++ {
			err = <-resChan
			if err != nil {
				return err
			}
		}
	} else {
		// Wiki search
		z := hd.ZimService.FindWikiFile(wiki)
		if z == nil {
			http.NotFound(w, r)
			return ErrNotFound
		}
		source = z.File.Title()

		// Search for query in WIKI
		results, err = searchWiki(query[0], z)
		if err != nil {
			return err
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
