package upload

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	tables := []struct {
		s       string
		hash    string
		maxSize int
		e       bool
	}{
		{"testing", "ae2b1fca515949e5d54fb22b8ed95575", 7, false},
		{"testing", "ae2b1fca515949e5d54fb22b8ed95575", 6, true},
	}

	for _, v := range tables {
		fs := &FileStore{
			maxFileSize: v.maxSize,
		}

		hash, err := fs.calculateHash(strings.NewReader(v.s))
		if v.e {
			assert.Equal(t, "", hash, "should not calculate hash on error")
			assert.EqualError(t, err, "file too big", "should recognize too big files")
		} else {
			assert.Equal(t, v.hash, hash, "should calculate hash correctly")
			assert.Nil(t, err)
		}
	}
}

func TestNew(t *testing.T) {
	fs, err := New("./testing/", 10)

	expected := FileStore{
		dir:         "./testing/",
		index:       make(map[string]string),
		maxFileSize: 10,
	}
	assert.Nil(t, err, "should not cause error creating")
	assert.Equal(t, fs, expected, "should create as expected")

	err = os.Remove("./testing/")
	assert.Nil(t, err, "should not cause error cleaning up")
}

func TestFileStore_Write(t *testing.T) {
	dir := "testing"
	name := t.Name()
	fs, err := New(dir, 10)
	assert.Nil(t, err, "should not cause error creating struct")
	testMessage := "testing"

	h, err := fs.Write(name, strings.NewReader(testMessage))
	assert.Nil(t, err, "writing should not cause error")
	assert.Equal(t, Handle("ae2b1fca515949e5d54fb22b8ed95575"), h, "should return the correct handle")

	assert.Equal(t, name, fs.index["ae2b1fca515949e5d54fb22b8ed95575"], "should store the index correctly")

	b, err := ioutil.ReadFile(filepath.Join(dir, name))
	assert.Equal(t, testMessage, string(b), "file content should match")

	err = os.Remove(name)
	assert.Nil(t, err, "should not cause error cleaning up testfile")
}
