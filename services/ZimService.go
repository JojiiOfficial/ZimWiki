package services

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// ZimService manage zim files
type ZimService struct {
	Dir   string
	files []ZimFile
}

// NewZimService create a new zimservice
func NewZimService(dir string) *ZimService {
	return &ZimService{
		Dir: dir,
	}
}

// Start starts the zimservice
func (zs *ZimService) Start() error {
	// Load all zimfiles in given directorys
	if err := zs.loadFiles(); err != nil {
		return err
	}

	log.Infof("Successfully loaded %d zim file(s)", len(zs.files))

	return nil
}

// GetFiles in dir
func (zs ZimService) GetFiles() []ZimFile {
	return zs.files
}

// Load all files in given Dir
func (zs *ZimService) loadFiles() error {
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

			zs.files = append(zs.files, ZimFile{
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
