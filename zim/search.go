package zim

import (
	"hash/fnv"
	"strings"

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
	if zf.HasIndexFile() {
		// Try to do an index search
		res, success := zf.indexSearch(query)
		if success {
			return res
		}
	}

	// If index search was not successful
	// or index file doesn't exists
	// do a full search
	return zf.fullSearch(query)
}

func (zf *File) indexSearch(query string) ([]SRes, bool) {
	// Get the index reader for the file
	indexReader := zf.GetIndexReader()
	if indexReader == nil {
		return nil, false
	}

	type wasAdded = struct{}
	alreadyAdded := make(map[uint32]wasAdded, 0)
	dirEntries := make([]SRes, 0)

	// Find in index file
	indexReader.ForEachSimilar(query, func(title string, pos uint32) error {
		zf.Mx.Lock()
		entry, err := zf.EntryAtTitlePosition(pos)
		zf.Mx.Unlock()
		if err != nil {
			return err
		}

		var key = hash(entry.URL())
		if _, ok := alreadyAdded[key]; ok {
			return nil
		}

		// Follow redirect
		if entry.IsRedirect() {
			fl, err := zf.FollowRedirect(&entry)
			if err == nil {
				entry = fl
			} else {
				logrus.Warn(err)
			}

			// Also add redirected to map
			alreadyAdded[hash(entry.URL())] = wasAdded{}
		}

		alreadyAdded[key] = wasAdded{}

		dirEntries = append(dirEntries, SRes{
			File:           zf,
			DirectoryEntry: &entry,
			Similarity:     getStrDest(title, query),
		})
		return nil
	})

	return dirEntries, true
}

func (zf *File) fullSearch(query string) []SRes {
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

			dirEntries = append(dirEntries, SRes{
				File:           zf,
				DirectoryEntry: entry,
				Similarity:     getStrDest(title, query),
			})
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
