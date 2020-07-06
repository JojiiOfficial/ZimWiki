package zim

import (
	"hash/fnv"
	"strings"
	"sync"

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

	if len(a[i].Title()) != len(a[j].Title()) {
		return len(a[i].Title()) > len(a[j].Title())
	}

	// If we have two items with the same similarity,
	// we use the alphabet order to provide the
	// the same result for the same query
	return string(a[i].DirectoryEntry.Title()) < string(a[j].DirectoryEntry.Title())
}

// SearchForEntry in zim file
func (zf *File) SearchForEntry(query string) []SRes {
	mx := sync.Mutex{}
	dirEntries := make([]SRes, 0)

	type wasAdded = struct{}
	alreadyAdded := make(map[uint32]wasAdded, 0)

	zf.Mx.Lock()
	zf.ForEachEntryWithURLPrefix(zim.NamespaceArticles, nil, int(zf.ArticleCount()), func(entry *zim.DirectoryEntry) {
		if len(entry.Title()) == 0 {
			return
		}

		var key = hash(entry.URL())
		if _, ok := alreadyAdded[key]; ok {
			return
		}

		title := string(entry.Title())

		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			// Follow redirect
			if entry.IsRedirect() {
				fl, err := zf.FollowRedirect(entry)
				if err == nil {
					entry = &fl
				} else {
					logrus.Warn(err)
				}

				// Also add redirected to map
				alreadyAdded[hash(entry.URL())] = wasAdded{}
			}

			alreadyAdded[key] = wasAdded{}

			mx.Lock()
			dirEntries = append(dirEntries, SRes{
				File:           zf,
				DirectoryEntry: entry,
				Similarity:     getStrDest(title, query),
			})
			mx.Unlock()
		}
	})

	zf.Mx.Unlock()
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
		add = 70
	}

	return int(float32(levenshtein.Similarity(a, b, nil)))*100 + add
}

func hash(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}
