package zim

import "github.com/tim-st/go-zim"

// File represents a zimfile
// and a local assigned path
type File struct {
	*zim.File
	Path string
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
