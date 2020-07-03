package handlers

import (
	"io"
	"os"
)

func serveStaticFile(path string, w io.Writer) error {
	// Try to open file
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	// Send file
	buff := make([]byte, 1024*1024)
	_, err = io.CopyBuffer(w, f, buff)

	return err
}
