package handlers

import (
	"html/template"
	"net/http"
	"path"
)

var (
	// TemplateCache is a cache of templates
	TemplateCache = make(map[string]*template.Template)
)

// 						      //
// --- Template data structs  //
// 						      //

// TemplateData data for the base template
type TemplateData struct {
	Favtype      string
	FavIcon      string
	Wiki         string
	Namespace    string

	HomeTemplateData
	WikiViewTemplateData
	SearchTemplateData
}

// HomeCards contain data for the
// cards displayed on the home site
type HomeCards struct {
	Image string
	Title string
	Text  string
	Link  string
}

// HomeTemplateData contain data
// for the home template
type HomeTemplateData struct {
	Cards []HomeCards
}

// WikiViewTemplateData data for wiki view page
type WikiViewTemplateData struct {
	Source string
}

// SearchTemplateData data for search
type SearchTemplateData struct {
	SearchSource string
	QueryText    string
	Results      []SearchResult
}

// SearchResult represents a single result of a search query
type SearchResult struct {
	Img   string
	Wiki  string
	Title string
	Link  string
}

// 						    //
// --- Template functions   //
// 						    //

// Load and execute
func serveTemplate(tmpFile string, w http.ResponseWriter, btd TemplateData) error {
	var err error
	tmplName := path.Base(tmpFile)

	w.Header().Set("Content-Type", "text/html utf-8")

	// Find in cache
	tmpl, has := TemplateCache[tmplName]

	if !has {
		// Parse if not in cache
		tmpl, err = template.New(tmplName).ParseFiles(BaseTemplate, tmpFile)
		if err != nil {
			return err
		}

		// Cache template
		TemplateCache[tmplName] = tmpl
	}

	if len(btd.Wiki) == 0 {
		btd.Wiki = "-"
	}

	// Execute template
	return tmpl.ExecuteTemplate(w, path.Base(BaseTemplate), btd)
}
