package handlers

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Assets for static assets
func Assets(w http.ResponseWriter, r *http.Request, hd HandlerData) error {
	vars := mux.Vars(r)

	if !checkVars(vars, "type", "file") {
		return fmt.Errorf("Missing parameter")
	}

	reqFile := vars["file"]
	assetType := vars["type"]

	// Get local path
	path := filepath.Clean(path.Join(AssetsPath, assetType, reqFile))

	// Serve Asset
	return serveRawFile(path, w)
}
