package handlers

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/JojiiOfficial/gaw"
	"github.com/gorilla/mux"
)

// Assets for static assets
func Assets(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	if !checkVars(vars, "type", "file") {
		return fmt.Errorf("Missing parameter")
	}

	reqFile := vars["file"]
	assetType := vars["type"]

	// Get local path
	path := filepath.Clean(path.Join(AssetsPath, assetType, reqFile))

	// Check if file exists
	if !gaw.FileExists(path) {
		http.NotFound(w, r)
		return nil
	}

	// Serve Asset
	return serveRawFile(path, w)
}
