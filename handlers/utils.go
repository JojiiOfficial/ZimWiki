package handlers

import (
	"io"
	"mime"
	"net/http"
	"strings"
)

// Serve file and set content-type accordingly
func serveRawFile(path string, w http.ResponseWriter) error {
	// Cache files
	w.Header().Set("Cache-Control", "max-age=31536000, public")

	// Detect and set mimetype
	m := mime.TypeByExtension(path[strings.LastIndex(path, "."):])
	w.Header().Set("Content-Type", m)

	return serveStaticFile(path, w)
}

func serveStaticFile(path string, w io.Writer) error {
	// Try to open file
	f, err := WebFS.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}

	// Send file
	buff := make([]byte, 1024*1024)
	_, err = io.CopyBuffer(w, f, buff)

	return err
}

// Check given vars
func checkVars(vars map[string]string, keys ...string) bool {
	for _, key := range keys {
		if _, ok := vars[key]; !ok {
			return false
		}
	}

	return true
}
