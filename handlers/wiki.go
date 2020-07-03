package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.jojii.de/jojii/zimserver/zim"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	gzim "github.com/tim-st/go-zim"
)

// Find a given wiki page
func findWikiPage(vars map[string]string, w http.ResponseWriter, r *http.Request, hd *HandlerData) (*gzim.Namespace, *zim.File, *gzim.DirectoryEntry, error) {
	if !checkVars(vars, "wikiID", "namespace", "file") {
		return nil, nil, nil, fmt.Errorf("Missing parameter")
	}

	// Throw error for invalid namespaces
	reqNamespace := vars["namespace"]
	if !strings.ContainsAny(reqNamespace, "ABIJMUVWX-") || len(reqNamespace) > 1 {
		http.NotFound(w, r)
		return nil, nil, nil, ErrNamespaceNotFound
	}

	// Parse namespace
	namespace := gzim.Namespace(reqNamespace[0])

	// reqFileURL is the url of the
	// requested file inside a wiki
	reqFileURL := vars["file"]

	// WikiID represents the zim UUID
	reqWikiID := vars["wikiID"]

	// Find requested wiki file by given ID
	z := hd.ZimService.FindWikiFile(reqWikiID)
	if z == nil {
		return nil, nil, nil, ErrNotFound
	}

	entry, _, found := z.EntryWithURL(namespace, []byte(reqFileURL))
	if !found {
		http.NotFound(w, r)
		return nil, nil, nil, ErrNotFound
	}

	// Follow redirect
	if entry.IsRedirect() {
		entry, _ = z.FollowRedirect(&entry)
	}

	return &namespace, z, &entry, nil
}

// WikiRaw handle direct wiki requests, without embedding into the webUI
func WikiRaw(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	vars := mux.Vars(r)

	_, z, entry, err := findWikiPage(vars, w, r, hd)
	if err != nil {
		return err
	}

	var blobReader, _, blobReaderErr = z.BlobReader(entry)
	if blobReaderErr != nil {
		log.Printf("Entry found but loading blob data failed for URL: %s with error %s\n", r.URL.Path, blobReaderErr)
		http.Error(w, blobReaderErr.Error(), http.StatusFailedDependency)
		return nil
	}

	// If file was found, set Mimetype accordingly
	if mimetypeList := z.MimetypeList(); int(entry.Mimetype()) < len(mimetypeList) {
		w.Header().Set("Content-Type", mimetypeList[entry.Mimetype()])
	}

	// Copy response
	buff := make([]byte, 1024*1024)
	_, err = io.CopyBuffer(w, blobReader, buff)
	return err
}

// WikiPreview sends a human friendly preview page for a WIKI site
func WikiPreview(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {

	return nil
}
