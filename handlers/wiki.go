package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// WikiRaw handle direct wiki requests, without embedding into the webUI
func WikiRaw(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	vars := mux.Vars(r)

	if !checkVars(vars, "wikiID", "namespace", "file") {
		return fmt.Errorf("Missing parameter")
	}

	// Throw error for invalid namespaces
	reqNamespace := vars["namespace"]
	if !strings.ContainsAny(reqNamespace, "ABIJMUVWX-") || len(reqNamespace) > 1 {
		log.Error("Namespace not found!")
		http.NotFound(w, r)
		return nil
	}

	// Parse namespace
	namespace := zim.Namespace(reqNamespace[0])

	// reqFileURL is the url of the
	// requested file inside a wiki
	reqFileURL := vars["file"]

	// WikiID represents the zim UUID
	reqWikiID := vars["wikiID"]

	// Find requested wiki file by given ID
	z := hd.ZimService.FindWikiFile(reqWikiID)
	if z == nil {
		http.NotFound(w, r)
		return nil
	}

	entry, _, found := z.EntryWithURL(namespace, []byte(reqFileURL))
	if !found {
		fmt.Println("not found")
		http.NotFound(w, r)
		return nil
	}

	// Follow redirect
	if entry.IsRedirect() {
		entry, _ = z.FollowRedirect(&entry)
	}

	var blobReader, _, blobReaderErr = z.BlobReader(&entry)
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
	_, err := io.CopyBuffer(w, blobReader, buff)
	return err
}
