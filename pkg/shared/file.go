package shared

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type HashType string

const (
	HashTypeMD5    HashType = "md5"
	HashTypeSHA1   HashType = "sha1"
	HashTypeSHA256 HashType = "sha256"
	HashTypeSHA512 HashType = "sha512"
)

type File struct {
	Path   string              `json:"path"`
	Size   *int64              `json:"size,omitempty"`
	Hashes map[HashType]string `json:"hashes"`
}

func (f File) GetUUID() string {
	return CalculateUUIDv5FromMap(map[string]interface{}{
		"path": f.Path,
	})
}

func (f File) GetObservableType() ObservableType {
	return ObservableTypeFile
}

func ListFiles(ctx context.Context, root string) (chan *File, error) {
	paths, err := ListFilePaths(ctx, root)
	if err != nil {
		return nil, err
	}

	ch := make(chan *File)
	go func() {
		defer close(ch)

		wg := sync.WaitGroup{}
		for path := range paths {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()

				file, err := GetFile(path)
				if err != nil {
					return
				}
				ch <- file
			}(path)
		}
	}()
	return ch, nil
}

func ListFilePaths(ctx context.Context, root string) (chan string, error) {
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, errors.New("root is not a directory")
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			ch <- path
			return nil
		})
	}()
	return ch, nil
}

func GetFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file := &File{
		Path: path,
	}
	st, err := f.Stat()
	if err != nil {
		sz := st.Size()
		file.Size = &sz
	}
	hashes, err := GetReaderHashes(f)
	if err != nil {
		return nil, err
	}
	file.Hashes = hashes
	return file, nil
}

func GetFileHashes(path string) (map[HashType]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return GetReaderHashes(f)
}

func GetReaderHashes(rd io.Reader) (map[HashType]string, error) {
	md5h := md5.New()
	sha1h := sha1.New()
	sha256h := sha256.New()
	sha512h := sha256.New()

	pagesize := os.Getpagesize()
	r := bufio.NewReaderSize(rd, pagesize)
	mw := io.MultiWriter()

	_, err := io.Copy(mw, r)
	if err != nil {
		return nil, err
	}

	hashes := map[HashType]string{
		HashTypeMD5:    hex.EncodeToString(md5h.Sum(nil)),
		HashTypeSHA1:   hex.EncodeToString(sha1h.Sum(nil)),
		HashTypeSHA256: hex.EncodeToString(sha256h.Sum(nil)),
		HashTypeSHA512: hex.EncodeToString(sha512h.Sum(nil)),
	}
	return hashes, nil
}
