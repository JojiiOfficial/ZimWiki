package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.jojii.de/jojii/zimserver/zim"
	gzim "github.com/tim-st/go-zim"
)

func parseWikiRequest(w http.ResponseWriter, r *http.Request, hd *HandlerData) (*zim.File, *gzim.Namespace, *gzim.DirectoryEntry, bool) {
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
			return nil, nil, nil, false
		}
	}

	// Something in the request is missing
	if len(sPath) < 5 {
		newLoc := "/"

		// Try to use main page if
		// the page is the only
		// thing missing
		if len(sPath) >= 3 {
			if mainpage := zim.GetMainpageURL(z); len(mainpage) > 0 {
				newLoc = mainpage
			}
		}

		// Something is missing in the given URL
		http.Redirect(w, r, newLoc, http.StatusMovedPermanently)
		return nil, nil, nil, false
	}

	// Throw error for invalid namespaces
	reqNamespace := sPath[3]
	if !strings.ContainsAny(reqNamespace, "ABIJMUVWX-") || len(reqNamespace) > 1 {
		http.NotFound(w, r)
		return nil, nil, nil, false
	}

	// Parse namespace
	namespace := gzim.Namespace(reqNamespace[0])

	switch namespace {
	case gzim.NamespaceLayout, gzim.NamespaceArticles, gzim.NamespaceImagesFiles, gzim.NamespaceImagesText:
	default:
		http.NotFound(w, r)
		return nil, nil, nil, false
	}

	// reqFileURL is the url of the
	// requested file inside a wiki
	reqFileURL := strings.Join(sPath[4:], "/")

	z.Mx.Lock()
	entry, _, found := z.EntryWithURL(namespace, []byte(reqFileURL))
	z.Mx.Unlock()
	if !found {
		http.NotFound(w, r)
		return nil, nil, nil, false
	}

	// Follow redirect
	if entry.IsRedirect() {
		z.Mx.Lock()
		entry, _ = z.FollowRedirect(&entry)
		z.Mx.Unlock()
		http.Redirect(w, r, zim.GetRawWikiURL(z, entry), http.StatusNotFound)
		return nil, nil, nil, false
	}

	return z, &namespace, &entry, true
}

// WikiRaw handle direct wiki requests, without embedding into the webUI
func WikiRaw(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	// Find file and dirEntry
	z, _, entry, success := parseWikiRequest(w, r, hd)
	if !success {
		// We already handled
		// http errors & redirects
		return nil
	}

	// Create reader from requested file
	z.Mx.Lock()
	defer z.Mx.Unlock()
	blobReader, _, err := z.BlobReader(entry)
	if err != nil {
		return err
	}

	// Set Mimetype accordingly
	if mimetypeList := z.MimetypeList(); int(entry.Mimetype()) < len(mimetypeList) {
		w.Header().Set("Content-Type", mimetypeList[entry.Mimetype()])
	}

	// Cache files
	w.Header().Set("Cache-Control", "max-age=31536000, public")

	// Send raw file
	// TODO replace absolute links
	buff := make([]byte, 1024*1024)
	_, err = io.CopyBuffer(w, blobReader, buff)
	return err
}

// WikiView sends a human friendly preview page for a WIKI site
func WikiView(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	// Find file and dirEntry
	z, _, entry, success := parseWikiRequest(w, r, hd)
	if !success {
		// We already handled
		// http errors & redirects
		return nil
	}

	var favurl, favType string
	z.Mx.Lock()
	favIcon, err := z.Favicon()
	z.Mx.Unlock()
	if err == nil {
		if mimetypeList := z.MimetypeList(); int(favIcon.Mimetype()) < len(mimetypeList) {
			favType = mimetypeList[favIcon.Mimetype()]
		}
		favurl = zim.GetRawWikiURL(z, favIcon)
	}

	// Cache files
	w.Header().Set("Cache-Control", "max-age=31536000, public")

	return serveTemplate(WikiPageTemplate, w, TemplateData{
		FavIcon: favurl,
		Favtype: favType,
		Wiki:    z.GetID(),
		WikiViewTemplateData: WikiViewTemplateData{
			Source: zim.GetRawWikiURL(z, *entry),
		},
	})
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
