package zim

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/tim-st/go-zim"
)

// File represents a zimfile
// and a local assigned path
type File struct {
	*zim.File
	Path      string
	IndexFile string

	// Use Mutex cause we can only
	// have one routine reading
	// at the same time
	Mx sync.RWMutex
}

// FileInfo for zim archive
type FileInfo struct {
	Name            string
	Title           string
	ArticleCount    uint32
	Description     string
	LongDescription string
	Date            string
	FileSize        int
	Language        string
}

// GetDescription tries to return the long
// Description
func (zfi FileInfo) GetDescription() string {
	// Prefer long description
	if len(zfi.LongDescription) > 0 {
		return zfi.LongDescription
	}

	// Use normal description
	if len(zfi.Description) > 0 {
		return zfi.Description
	}

	return "No description provided"
}

// GetInfos for zim archive
func (zf *File) GetInfos() *FileInfo {
	return &FileInfo{
		Name:            zf.Name(),
		Title:           zf.Title(),
		ArticleCount:    zf.ArticleCount(),
		Description:     zf.Description(),
		LongDescription: zf.LongDescription(),
		Date:            zf.Date(),
		FileSize:        zf.Filesize(),
		Language:        zf.Language(),
	}
}

// GetID for zimfiles
func (zf *File) GetID() string {
	return zf.UUID().String()
}

// generateFileIndex for search
func (zf *File) generateFileIndex(w io.Writer) (uint32, error) {

	progress := uint32(0)
	writtenBytes := uint32(0)
	done := make(chan error, 1)

	go func() {
		// Loop through all entries
		err := zf.ForEachEntryAfterPosition(0, 0, func(entry *zim.DirectoryEntry, pos uint32) error {
			progress++

			if entry.IsArticle() {
				if len(entry.URL()) < 2 {
					return nil
				}

				sPos := strconv.FormatUint(uint64(pos), 36)

				// Build index entry
				var data []byte
				data = append(data, entry.URL()...)
				data = append(data, byte('\n'))
				data = append(data, []byte(sPos)...)
				data = append(data, byte('\n'))

				// Write data
				n, err := w.Write(data)
				if err != nil {
					return err
				}

				writtenBytes += uint32(n)
			}
			return nil
		})

		done <- err
	}()

	// Print progress
f:
	for {
		select {
		case <-done:
			fmt.Printf("\rIndexing %s ...done\n", zf.Name())
			break f
		case <-time.After(time.Second * 2):
			if progress > 0 {
				fmt.Printf("\rIndexing %s: %d%s", zf.Name(), progress*100/zf.ArticleCount(), "%")
			}
		}
	}

	return writtenBytes, nil
}

// HasIndexFile return true if wiki
// has an index file
func (zf *File) HasIndexFile() bool {
	if len(zf.IndexFile) == 0 {
		return false
	}
	s, err := os.Stat(zf.IndexFile)
	return err == nil && s.Size() > 0
}

// GetIndexReader gets indexReader for file
func (zf *File) GetIndexReader() *IndexReader {
	if !zf.HasIndexFile() {
		return nil
	}

	return &IndexReader{
		IndexFile: zf.IndexFile,
	}
}
