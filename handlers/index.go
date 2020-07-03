package handlers

import (
	"net/http"
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request) error {
	return serveTemplate(HomeTemplate, HomeTemplateData{
		Cards: []HomeCards{},
	}, w)
}
