package zim

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const dbFileName = ".ix.db"

// ...
var (
	ErrAlreadyInDB = errors.New("Index already exists")
)

// IndexDB a full db of index entries
type IndexDB struct {
	File    string       `json:"-"`
	LibPath string       `json:"-"`
	Entries []IndexEntry `json:"ies"`
}

// IndexEntry represents an
// entry (a index file) in an
//
type IndexEntry struct {
	IndexFile string `json:"if"`
	Checksum  string `json:"cs"`
}

// NewIndexDB read indexDB from file or create a new one
func NewIndexDB(libPath string) (*IndexDB, error) {
	var indexDB IndexDB

	path := filepath.Join(libPath, dbFileName)
	s, err := os.Stat(path)

	if err == nil && s.Size() > 0 {
		// File exists and is not empty
		r, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		indexDB.File = path
		indexDB.LibPath = libPath
		// Read IndexDB
		return &indexDB, json.Unmarshal(r, &indexDB)
	}

	// Create index DB and write to file
	indexDB = IndexDB{Entries: []IndexEntry{}, File: path, LibPath: libPath}
	return &indexDB, indexDB.Save(path)
}

// Save the db
func (indexDB *IndexDB) Save(path string) error {
	// Open the file. Create if not exists
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	// Create empty db
	d, err := json.Marshal(indexDB)
	if err != nil {
		return err
	}

	// Write indexDB
	_, err = f.Write(d)
	return err
}

// addIndex to db
func (indexDB *IndexDB) addIndex(entry IndexEntry) error {
	indexDB.Entries = append(indexDB.Entries, entry)
	return indexDB.Save(indexDB.File)
}

// AddIndexFile to DB. Calculates checksum automatically
func (indexDB *IndexDB) AddIndexFile(file string) error {
	if indexDB.GetEntry(file) != nil {
		return ErrAlreadyInDB
	}

	// If file still contains a dir, remove it
	file = indexDB.addPathPrefix(file)

	// Open file to index
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	// Generate hash of file
	sHash, err := fileChecksum(f)
	if err != nil {
		return err
	}

	// Add index
	return indexDB.addIndex(IndexEntry{
		Checksum:  sHash,
		IndexFile: removePathPrefix(file),
	})
}

// CheckFile if it is in DB and the Checksum matches
func (indexDB *IndexDB) CheckFile(file string) (bool, error) {
	// If file still contains a dir, remove it
	file = indexDB.addPathPrefix(file)

	// Check if file is in DB
	entry := indexDB.GetEntry(file)
	if entry == nil {
		return false, nil
	}

	// Open file
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return false, err
	}

	// Create checksum
	sum, err := fileChecksum(f)
	if err != nil {
		return false, err
	}

	// return sum == , nil
	return entry.Checksum == sum, nil
}

// GetEntry in IndexDB
func (indexDB *IndexDB) GetEntry(file string) *IndexEntry {
	file = removePathPrefix(file)
	for i := range indexDB.Entries {
		if indexDB.Entries[i].IndexFile == file {
			return &indexDB.Entries[i]
		}
	}

	return nil
}

func fileChecksum(f *os.File) (string, error) {
	hash := crc32.NewIEEE()
	buff := make([]byte, 1024*1024)
	_, err := io.CopyBuffer(hash, f, buff)
	if err != nil {
		return "", err
	}
	sHash := hex.EncodeToString(hash.Sum(nil))
	return sHash, nil
}

func removePathPrefix(file string) string {
	_, file = filepath.Split(file)
	return file
}

func (indexDB IndexDB) addPathPrefix(file string) string {
	file = removePathPrefix(file)
	return filepath.Join(indexDB.LibPath, file)
}
