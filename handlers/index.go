package handlers

import (
	"net/http"
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request) error {
	return serveStaticFile(BaseTemplate, w)
}
