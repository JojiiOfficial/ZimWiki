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

func (a ByPercentage) Len() int           { return len(a) }
func (a ByPercentage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPercentage) Less(i, j int) bool { return a[i].Similarity < a[j].Similarity }

// SearchForEntry in zim file
func (zf *File) SearchForEntry(query string, limit int) []SRes {
	zf.Mx.Lock()
	res := zf.File.EntriesWithNamespace(zim.NamespaceArticles, int(zf.ArticleCount()))
	zf.Mx.Unlock()

	mx := sync.Mutex{}
	dirEntries := make([]SRes, 0)

	gopool.New(len(res), runtime.NumCPU(), func(wg *sync.WaitGroup, pos, total, workerID int) interface{} {
		t := res[pos].Title()
		if len(t) == 0 || t[0] == 0 {
			return nil
		}

		title := string(t)

		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			// Follow redirect
			zf.Mx.Lock()
			if res[pos].IsRedirect() {
				fl, err := zf.FollowRedirect(&res[pos])
				if err == nil {
					res[pos] = fl
				} else {
					logrus.Warn(err)
				}
			}
			zf.Mx.Unlock()

			mx.Lock()
			dirEntries = append(dirEntries, SRes{
				File:           zf,
				DirectoryEntry: &res[pos],
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
	var add float32
	if strings.HasPrefix(strings.ToLower(a), strings.ToLower(b)) {
		add = 0.6
	}

	return int(float32(levenshtein.Similarity(a, b, nil))+add) * 100
}
