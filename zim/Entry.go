package zim

import (
	"fmt"

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
