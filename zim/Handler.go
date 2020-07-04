package zim

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// Handler manage zim files
type Handler struct {
	Dir   string
	files []File

	// Cache for faster file search
	fileCache map[string]*File

	Mx sync.Mutex
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
func (zs *Handler) GetFiles() []File {
	return zs.files
}

// Load all files in given Dir
func (zs *Handler) loadFiles() error {
	var success, errs int

	filepath.Walk(zs.Dir, func(path string, info os.FileInfo, err error) error {
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
				errs++
				log.Error(errors.Wrap(err, path))

				// Ignore errors for now
				return nil
			}

			zs.files = append(zs.files, File{
				File: f,
				Path: path,
			})
			success++
		}

		return nil
	})

	if success == 0 && errs > 0 {
		log.Fatal("Too many errors")
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
