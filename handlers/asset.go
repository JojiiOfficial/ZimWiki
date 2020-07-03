package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

// Assets for static assets
func Assets(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	assetType, has := vars["type"]
	if !has {
		return fmt.Errorf("Parameter 'type' not given")
	}

	reqFile, has := vars["file"]
	if !has {
		return fmt.Errorf("Parameter 'file' not given")
	}

	path := path.Join(assetPath, assetType, reqFile)
	m := mime.TypeByExtension(path[strings.LastIndex(path, "."):])
	w.Header().Set("Content-Type", m)

	return serveStaticFile(path, w)
}
