package zim

import (
	"runtime"
	"strings"
	"sync"

	"github.com/JojiiOfficial/gopool"
	"github.com/agext/levenshtein"
	"github.com/sirupsen/logrus"
	"github.com/tim-st/go-zim"
)

// SRes search result with similarity
type SRes struct {
	*File
	*zim.DirectoryEntry
	Similarity int
}

// ByPercentage sort by perc
type ByPercentage []SRes

func (a ByPercentage) Len() int      { return len(a) }
func (a ByPercentage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPercentage) Less(i, j int) bool {
	if a[i].Similarity != a[j].Similarity {
		return a[i].Similarity < a[j].Similarity
	}

	// If we have two items with the same similarity,
	// we use the alphabet order to provide the
	// the same result for the same query
	return string(a[i].DirectoryEntry.Title()) < string(a[j].DirectoryEntry.Title())
}

// SearchForEntry in zim file
func (zf *File) SearchForEntry(query string, limit int) []SRes {
	zf.Mx.Lock()
	res := zf.File.EntriesWithNamespace(zim.NamespaceArticles, int(zf.ArticleCount()))
	zf.Mx.Unlock()

	mx := sync.Mutex{}
	dirEntries := make([]SRes, 0)

	gopool.New(len(res), runtime.NumCPU(), func(wg *sync.WaitGroup, pos, total, workerID int) interface{} {
		entry := &res[pos]
		if len(entry.Title()) == 0 {
			return nil
		}
		title := string(entry.Title())

		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			// Follow redirect
			zf.Mx.Lock()
			if entry.IsRedirect() {
				fl, err := zf.FollowRedirect(entry)
				if err == nil {
					entry = &fl
				} else {
					logrus.Warn(err)
				}
			}
			zf.Mx.Unlock()

			mx.Lock()
			dirEntries = append(dirEntries, SRes{
				File:           zf,
				DirectoryEntry: entry,
				Similarity:     getStrDest(title, query),
			})
			mx.Unlock()
		}

		return nil
	}).Run().Wait()

	return dirEntries
}

// Get string destination in percent and
// prefer equal same prefixes
func getStrDest(a, b string) int {
	al := strings.ToLower(a)
	bl := strings.ToLower(b)

	if al == bl {
		return 200
	}

	var add int
	if strings.HasPrefix(al, bl) {
		add = 60
	}

	return int(float32(levenshtein.Similarity(a, b, nil)))*100 + add
}
