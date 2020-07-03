package zim

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// Handler manage zim files
type Handler struct {
	Dir   string
	files []File

	// Cache for faster file search
	fileCache map[string]*File
}

// NewZim create a new zimservice
func NewZim(dir string) *Handler {
	return &Handler{
		Dir:       dir,
		fileCache: make(map[string]*File),
	}
}

// Start starts the zimservice
func (zs *Handler) Start() error {
	// Load all zimfiles in given directorys
	if err := zs.loadFiles(); err != nil {
		return err
	}

	log.Infof("Successfully loaded %d zim file(s)", len(zs.files))

	return nil
}

// GetFiles in dir
func (zs Handler) GetFiles() []File {
	return zs.files
}

// Load all files in given Dir
func (zs *Handler) loadFiles() error {
	err := filepath.Walk(zs.Dir, func(path string, info os.FileInfo, err error) error {
		// Ignore non regular files
		if !info.IsDir() {

			// Follow sysmlinks
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				// Follow link
				path, err = os.Readlink(path)
				if err != nil {
					return err
				}
			}

			// Try to open any file
			f, err := zim.Open(path)

			if err != nil {
				// Ignore errors for non .zim files
				if strings.HasSuffix(info.Name(), ".zim") {
					return err
				}

				return nil
			}

			zs.files = append(zs.files, File{
				File: f,
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// FindWikiFile by ID. File gets cached into a map
func (zs *Handler) FindWikiFile(zimFileID string) *File {
	if fil, has := zs.fileCache[zimFileID]; has {
		return fil
	}

	// Loop all files and find matching
	for i := range zs.files {
		file := &zs.files[i]

		if file.GetID() == zimFileID {
			zs.fileCache[file.GetID()] = file
			return file
		}
	}

	return nil
}
