package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
	"strconv"
	"math"

	"github.com/JojiiOfficial/ZimWiki/zim"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func searchSingle(query string, nbResultsPerPage int, resultsUntil int, wiki *zim.File) ([]zim.SRes, string, int, int) {
	// Search entries
	entries := wiki.SearchForEntry(query)
	// Sort them by similarity
	sort.Sort(sort.Reverse(zim.ByPercentage(entries)))

	// Calculate the first result displayed in the page
	resultsStart := resultsUntil - nbResultsPerPage

	// If the result of the calculation is negative, set the variable to 0
	if resultsStart < 0 {
		resultsStart = 0
	}

	// If the number of results to be displayed per page is lower than the total number of results
	if len(entries) < nbResultsPerPage {
		resultsUntil = len(entries)
	}

	// If the result to be the maximum result of the page is greater than the number of total results
	if resultsUntil > len(entries) {
		resultsUntil = len(entries)
	}

	// Calculate the total number of pages
	nbPages := int(math.Ceil(float64(len(entries)) / float64(nbResultsPerPage)))

	// Know the number of results
	resultText := "No result"

	if len(entries) == 1 {
		resultText = "1 result"
	} else if len(entries) != 0 {
		resultText = strconv.Itoa(len(entries)) + " results"
	}

	return entries[resultsStart:resultsUntil], resultText, len(entries), nbPages
}

// Search in all available wikis
func searchGlobal(query string, nbResultsPerPage int, resultsUntil int, handler *zim.Handler) ([]zim.SRes, string, int, int) {
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

	resultsStart := resultsUntil - nbResultsPerPage

	// If the result of the calculation is negative, set the variable to 0
	if resultsStart < 0 {
		resultsStart = 0
	}

	// If the number of results to be displayed per page is lower than the total number of results
	if len(results) < nbResultsPerPage {
		resultsUntil = len(results)
	}

	// If the result to be the maximum result of the page is greater than the number of total results
	if resultsUntil > len(results) {
		resultsUntil = len(results)
	}

	// Calculate the total number of pages
	nbPages := int(math.Ceil(float64(len(results)) / float64(nbResultsPerPage)))

	// Know the number of results
	resultText := "No result"

	if len(results) == 1 {
		resultText = "1 result"
	} else if len(results) != 0 {
		resultText = strconv.Itoa(len(results)) + " results"
	}

	return results[resultsStart:resultsUntil], resultText, len(results), nbPages
}

// Search handles serach requests
func Search(w http.ResponseWriter, r *http.Request, hd HandlerData) error {
	vars := mux.Vars(r)
	wiki, ok := vars["wiki"]
	if !ok {
		return fmt.Errorf("Missing parameter")
	}
	var query string
	var actualPageNb int

	if r.Method == "POST" {
		// Get Post search query
		query = r.PostFormValue("sQuery")
		// Get the number of the current page
		actualPageNb, _ = strconv.Atoi(r.PostFormValue("pageNumber"))
	}

	if len(query) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return nil
	}

	// If the variable is not defined, we set it to 1
	if actualPageNb == 0 {
		actualPageNb = 1
	}

	// The number of results displayed per page
	nbResultsPerPage := 12

	// A small calculus to know the last result to display
	resultsUntil := nbResultsPerPage * actualPageNb

	var res []zim.SRes
	var source string
	var resultText string
	var nbResults int
	var nbPages int

	start := time.Now()
	if wiki == "-" {
		source = "global search"

		res, resultText, nbResults, nbPages = searchGlobal(query, nbResultsPerPage, resultsUntil, hd.ZimService)
	} else {
		// Wiki search
		z := hd.ZimService.FindWikiFile(wiki)
		if z == nil {
			http.NotFound(w, r)
			return ErrNotFound
		}
		source = z.File.Title()

		// Search for query in WIKI
		res, resultText, nbResults, nbPages = searchSingle(query, nbResultsPerPage, resultsUntil, z)
	}

	var previousPage int
	var nextPage int

	if actualPageNb != 1 {
		previousPage = actualPageNb - 1
	}

	if  nbResults - (nbResultsPerPage * (actualPageNb)) > 0 {
		nextPage = actualPageNb + 1
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
			ActualPageNb: actualPageNb,
			NbPages:      nbPages,
			PreviousPage: previousPage,
			NextPage:     nextPage,
		},
	})
}
