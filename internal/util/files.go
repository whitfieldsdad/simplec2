package util

import (
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"
	"github.com/mitchellh/go-homedir"
)

type File struct {
	Path       string          `json:"path,omitempty"`
	Name       string          `json:"name,omitempty"`
	Directory  string          `json:"directory,omitempty"`
	Size       *int64          `json:"size,omitempty"`
	Hashes     []Hash          `json:"hashes,omitempty"`
	Timestamps *FileTimestamps `json:"timestamps,omitempty"`
}

func (f File) GetArtifactType() ArtifactType {
	return ArtifactTypeFile
}

type FileTimestamps struct {
	Modified *time.Time `json:"modified,omitempty"`
	Accessed *time.Time `json:"accessed,omitempty"`
	Changed  *time.Time `json:"changed,omitempty"`
	Born     *time.Time `json:"born,omitempty"`
}

func GetFile(path string) (*File, error) {
	path = getRealFilePath(path)
	return getFile(path)
}

func getFile(path string) (*File, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file := File{
		Path:      path,
		Name:      filepath.Base(path),
		Directory: filepath.Dir(path),
	}
	file.Timestamps, _ = getFileTimestamps(path)

	if !st.IsDir() {
		file.Hashes, _ = getFileHashes(path)

		size := st.Size()
		file.Size = &size
	}
	return &file, nil
}

func getFileHashes(path string) ([]Hash, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return getReaderHashes(file)
}

func getFileTimestamps(path string) (*FileTimestamps, error) {
	ts, err := times.Stat(path)
	if err != nil {
		return nil, err
	}
	m := ts.ModTime()
	a := ts.AccessTime()
	c := ts.ChangeTime()

	var b *time.Time
	if ts.HasBirthTime() {
		_b := ts.BirthTime()
		b = &_b
	}
	return &FileTimestamps{
		Modified: &m,
		Accessed: &a,
		Changed:  &c,
		Born:     b,
	}, nil
}

func getRealFilePath(path string) string {
	path, _ = homedir.Expand(path)
	path = os.ExpandEnv(path)
	path, _ = filepath.Abs(path)
	return path
}
