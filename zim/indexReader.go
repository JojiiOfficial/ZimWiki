package zim

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// IndexReader reads elements
// from a wiki index file
type IndexReader struct {
	IndexFile string
}

type foundFunc func(string, uint32) error

// ForEachSimilar finds similar entries
// inside an index file
func (indexReader IndexReader) ForEachSimilar(subs string, foundFunc foundFunc) error {
	f, err := os.OpenFile(indexReader.IndexFile, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	ql := strings.ToLower(subs)

	var i int
	var matchStr string
	for scanner.Scan() {
		line := scanner.Text()

		if i%2 != 0 {
			// Here is an integer
			pos, err := strconv.ParseUint(line, 36, 32)
			if err != nil {
				return err
			}

			if len(matchStr) > 0 {
				foundFunc(matchStr, uint32(pos))
				matchStr = ""
			}
		} else {
			l := strings.ToLower(line)
			if strings.Contains(l, ql) || ql == l {
				matchStr = line
			}
		}

		i++
	}

	return nil
}
