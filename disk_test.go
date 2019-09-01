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

func TestGenericDiskFunctions(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key := []byte("some-key")
	value := []byte("some value")

	val, err := disk.read(key)
	is.Equal(err, badger.ErrKeyNotFound)
	is.Nil(val)

	is.NoError(disk.write(key, value))

	val, err = disk.read(key)
	is.NoError(err)
	is.Equal(val, value)

	is.NoError(disk.delete(key))

	val, err = disk.read(key)
	is.Equal(err, badger.ErrKeyNotFound)
	is.Nil(val)

	clean(is)
}

func TestDiskCache(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key, value := "some-cache-key", "some-value"
	ttl := int32(3600)

	val, err := disk.CacheGet(key)
	is.Equal(err, badger.ErrKeyNotFound)
	is.Empty(val)

	is.NoError(disk.CacheSet(key, value, ttl))

	val, err = disk.CacheGet(key)
	is.NoError(err)
	is.Equal(val, value)

	is.NoError(disk.CacheSet(key, value, 0))
	val, err = disk.CacheGet(key)
	is.Empty(val)

	clean(is)
}

func TestDiskKV(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	loc := &Location{
		Bucket: "test",
		Key:    "test",
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

	clean(is)
}

func setup(is *assert.Assertions) *Disk {
	disk, err := NewDisk(dir)
	is.NoError(err)
	is.NotNil(disk)
	return disk
}

func clean(is *assert.Assertions) {
	is.NoError(os.RemoveAll(dir))
}
