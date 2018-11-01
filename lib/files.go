package lib

import (
	"io/ioutil"
	"path/filepath"
)

type FileListingEntry struct {
	Path string
	Name string
}

func ListFiles(dir string) ([]*FileListingEntry) {
	entries := make([]*FileListingEntry, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return entries
	}
	for _, file := range files {
		if !file.IsDir() {
			entries = append(entries, &FileListingEntry{
				Path: filepath.Join(dir, file.Name()),
				Name: file.Name(),
			})
		}
	}

	return entries
}
