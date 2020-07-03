package services

import "github.com/tim-st/go-zim"

// ZimFile represents a zimfile
// and a local assigned path
type ZimFile struct {
	*zim.File
	Path string
}

// ZimFileInfo for zim archive
type ZimFileInfo struct {
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
func (zfi ZimFileInfo) GetDescription() string {
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
func (zf *ZimFile) GetInfos() *ZimFileInfo {
	return &ZimFileInfo{
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
