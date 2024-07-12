package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type PathTransformFunc func(string) string

func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
}

func DefaultPathTransformFunc(key string) string {
	return key
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"

	pathAndFilename := pathName + "/" + filename

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written to (%d) bytes to disk: %s", n, pathAndFilename)

	return nil
}
