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

// 						    //
// --- Template functions   //
// 						    //

// Load and execute
func serveTemplate(tmpFile string, tmpData interface{}, w http.ResponseWriter) error {
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
		// TemplateCache[tmplName] = tmpl
	}

	// Execute template
	return tmpl.ExecuteTemplate(w, path.Base(BaseTemplate), tmpData)
}
