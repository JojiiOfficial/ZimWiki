package zim

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// ...
const (
	WikiRawPrefix = "/wiki/raw/"
	WikiPrefix    = "/wiki/preview/"
)

// GetRawWikiURL returns the wiki url for the given DE
func GetRawWikiURL(zimFile File, entry zim.DirectoryEntry) string {
	return fmt.Sprintf(WikiRawPrefix+"%s/%s/%s", zimFile.GetID(), string(entry.Namespace()), string(entry.URL()))
}

// GetWikiURL returns the wiki url for the given DE
func GetWikiURL(zimFile File, entry zim.DirectoryEntry) string {
	return fmt.Sprintf(WikiPrefix+"%s/%s/%s", zimFile.GetID(), string(entry.Namespace()), string(entry.URL()))
}

// GetMainpageName of zimFile
func GetMainpageName(zimFile File) *zim.DirectoryEntry {
	mp, err := zimFile.MainPage()
	if err != nil {
		log.Error(err)
		return nil
	}

	return &mp
}

// GetMainpageURL of zimFile
func GetMainpageURL(zimFile File) string {
	mp := GetMainpageName(zimFile)
	if mp == nil {
		return ""
	}

	return GetWikiURL(zimFile, *mp)
}
