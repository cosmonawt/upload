package upload

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"errors"
	"path/filepath"
)

type FileStore struct {
	dir         string            // base directory
	index       map[string]string // stores references from hash to filename
	maxFileSize int               // max file size in bytes
}

type Handle string

func New(directory string, maxFileSize int) (store FileStore, err error) {
	err = os.Mkdir(directory, os.FileMode(0777))
	if err != nil && !os.IsExist(err) {
		return FileStore{}, err
	}

	return FileStore{directory, make(map[string]string), maxFileSize}, nil
}

func (fs *FileStore) calculateHash(src io.Reader) (hash string, err error) {
	h := md5.New()
	s, err := io.Copy(h, src)
	if err != nil {
		return "", err
	}
	if s > int64(fs.maxFileSize) {
		return "", errors.New("file too big")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (fs *FileStore) Write(filename string, src io.Reader) (handle Handle, err error) {
	hash, err := fs.calculateHash(src)
	if err != nil {
		return Handle(""), err
	}

	file, err := os.Create(filepath.Join("./testing/", filename))
	if err != nil {
		return Handle(""), err
	}
	defer file.Close()

	if _, err := io.Copy(file, src); err != nil {
		return Handle(""), err
	}

	fs.index[hash] = filename
	return Handle(hash), nil
}

/*func (fs *FileStore) Read(handle Handle, dst *io.Writer) (err error) {

}

func (fs *FileStore) Delete(handle Handle) (err error) {

}*/
