package zim

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// New create a new zimservice
func New(dir string) *Handler {
	return &Handler{
		Dir:       dir,
		fileCache: make(map[string]*File),
	}
}

// Start starts the zimservice
func (zs *Handler) Start(libPath string) error {

	// Load all zimfiles in given directorys
	if err := zs.loadFiles(); err != nil {
		return err
	}

	log.Infof("Successfully loaded %d zim file(s)", len(zs.files))

	return zs.GenerateIndex(libPath)
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
		if !info.IsDir() && !strings.HasSuffix(path, ".ix") && !strings.HasSuffix(path, ".ix.db") {
			// We want to use the real
			// path ond disk
			realPath := path

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
				Path: realPath,
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

// GenerateIndex for search queries
func (zs *Handler) GenerateIndex(libPath string) error {
	s := uint32(0)

	indexDB, err := NewIndexDB(libPath)
	if err != nil {
		return err
	}

	// Create index for all files
	for i := range zs.files {
		file := &zs.files[i]

		fmt.Printf("\rIndexing %s", file.Name())

		// Set index file
		fdir, fname := filepath.Split(file.Path)
		fname = "." + fname + ".ix"
		file.IndexFile = filepath.Join(fdir, fname)

		// Check file validation
		ok, err := indexDB.CheckFile(file.IndexFile)
		if err != nil {
			return err
		}
		// Skip file if index is still valid
		if ok {
			fmt.Printf("\rIndexing %s ...exists\n", file.Filename())
			continue
		}

		// Create new Index file
		f, err := os.OpenFile(file.IndexFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}

		// Generate index
		size, err := file.generateFileIndex(f)
		if err != nil {
			return err
		}

		f.Close()

		// Add index to DB
		err = indexDB.AddIndexFile(file.IndexFile)
		if err != nil {
			return err
		}

		s += size
	}

	if s > 0 {
		fmt.Printf("Full index size: %dMB\n", s/1000/1000)
	}

	return nil
}
