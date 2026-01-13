package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
)

type HashType string

const (
	HashTypeMD5    HashType = "md5"
	HashTypeSHA1   HashType = "sha1"
	HashTypeSHA256 HashType = "sha256"
)

var HashAlgorithms = []HashType{HashTypeMD5, HashTypeSHA1, HashTypeSHA256}

type Hash struct {
	Algorithm HashType `json:"algorithm"`
	Value     string   `json:"value"`
}

func getReaderHashes(rd io.Reader) ([]Hash, error) {
	writers := make(map[HashType]io.Writer)
	for _, algorithm := range HashAlgorithms {
		switch algorithm {
		case HashTypeMD5:
			writers[algorithm] = md5.New()
		case HashTypeSHA1:
			writers[algorithm] = sha1.New()
		case HashTypeSHA256:
			writers[algorithm] = sha256.New()
		}
	}
	mw := io.MultiWriter(writers[HashTypeMD5], writers[HashTypeSHA1], writers[HashTypeSHA256])
	if _, err := io.Copy(mw, rd); err != nil {
		return nil, err
	}
	hashes := make([]Hash, 0)
	for algorithm, writer := range writers {
		hashes = append(hashes, Hash{
			Algorithm: algorithm,
			Value:     hex.EncodeToString(writer.(hash.Hash).Sum(nil)),
		})
	}
	return hashes, nil
}
