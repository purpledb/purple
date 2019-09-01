package strato

import (
	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	dir = "tmp"
)

func TestDiskKV(t *testing.T) {
	is := assert.New(t)

	disk, err := NewDisk(dir)
	is.NoError(err)
	is.NotNil(disk)

	loc := &Location{
		Bucket: "test",
		Key: "test",
	}

	val := &Value{
		Content: []byte("here is some test content"),
	}

	is.NoError(disk.KVPut(loc, val))

	fetched, err := disk.KVGet(loc)
	is.NoError(err)
	is.Equal(fetched, val)

	is.NoError(disk.KVDelete(loc))

	fetched, err = disk.KVGet(loc)
	is.Error(err)
	is.Equal(err.Error(), badger.ErrKeyNotFound.Error())
	is.Nil(fetched)

	teardown(is, dir)
}

func teardown(is *assert.Assertions, dir string) {
	is.NoError(os.RemoveAll(dir))
}