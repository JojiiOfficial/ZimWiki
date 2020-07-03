package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.jojii.de/jojii/zimserver/zim"
	gzim "github.com/tim-st/go-zim"
)

// WikiRaw handle direct wiki requests, without embedding into the webUI
func WikiRaw(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	// Split requested path by /
	sPath := strings.Split(parseURLPath(r.URL), "/")

	var reqWikiID string
	var z *zim.File

	// We can use zim.File for getting
	// The desired wiki and mainpage
	if len(sPath) > 2 {
		// WikiID represents the zim UUID
		reqWikiID = sPath[2]

		hd.ZimService.Mx.Lock()
		// Find requested wiki file by given ID
		z = hd.ZimService.FindWikiFile(reqWikiID)
		hd.ZimService.Mx.Unlock()
		if z == nil {
			return nil
		}
	}

	// Something in the request is missing
	if len(sPath) < 5 {
		newLoc := "/"

		// Try to use main page if
		// the page is the only
		// thing missing
		if len(sPath) == 4 {
			if mainpage := zim.GetMainpageURLRaw(z); len(mainpage) > 0 {
				newLoc = mainpage
			}
		}

		// Something is missing in the given URL
		w.Header().Set("Location", newLoc)
		w.WriteHeader(http.StatusMovedPermanently)
		return nil
	}

	// Throw error for invalid namespaces
	reqNamespace := sPath[3]
	if !strings.ContainsAny(reqNamespace, "ABIJMUVWX-") || len(reqNamespace) > 1 {
		http.NotFound(w, r)
		return nil
	}

	// Parse namespace
	namespace := gzim.Namespace(reqNamespace[0])

	switch namespace {
	case gzim.NamespaceLayout, gzim.NamespaceArticles, gzim.NamespaceImagesFiles, gzim.NamespaceImagesText:
	default:
		http.NotFound(w, r)
		return nil
	}

	// reqFileURL is the url of the
	// requested file inside a wiki
	reqFileURL := strings.Join(sPath[4:], "/")

	z.Mx.Lock()
	entry, _, found := z.EntryWithURL(namespace, []byte(reqFileURL))
	z.Mx.Unlock()
	if !found {
		http.NotFound(w, r)
		return nil
	}

	// Follow redirect
	if entry.IsRedirect() {
		z.Mx.Lock()
		entry, _ = z.FollowRedirect(&entry)
		z.Mx.Unlock()
		http.Redirect(w, r, zim.GetRawWikiURL(z, entry), http.StatusNotFound)
		return nil
	}

	// Create reader from requested file
	z.Mx.Lock()
	defer z.Mx.Unlock()
	blobReader, _, err := z.BlobReader(&entry)
	if err != nil {
		return err
	}

	// Set Mimetype accordingly
	if mimetypeList := z.MimetypeList(); int(entry.Mimetype()) < len(mimetypeList) {
		w.Header().Set("Content-Type", mimetypeList[entry.Mimetype()])
	}

	// Copy response
	_, err = io.Copy(w, blobReader)
	return err
}

// WikiPreview sends a human friendly preview page for a WIKI site
func WikiPreview(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	return nil
}

func parseURLPath(u *url.URL) string {
	path := u.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return path
}
