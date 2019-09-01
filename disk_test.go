package strato

import (
	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	dir = "tmp/badger"
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

func TestDiskCounter(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key := "some-counter-key"

	count, err := disk.CounterGet(key)
	is.NoError(err)
	is.Zero(count)

	is.NoError(disk.CounterIncrement(key, int64(100)))

	count, err = disk.CounterGet(key)
	is.NoError(err)
	is.Equal(count, int64(100))

	is.NoError(disk.CounterIncrement(key, int64(-200)))
	count, err = disk.CounterGet(key)
	is.NoError(err)

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

func TestDiskSet(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key := "some-set"

	set, err := disk.GetSet(key)
	is.NoError(err)
	is.Empty(set)

	is.NoError(disk.AddToSet(key, "some-item"))

	set, err = disk.GetSet(key)
	is.NoError(err)
	is.Len(set, 1)
	is.Equal(set[0], "some-item")

	is.NoError(disk.RemoveFromSet(key, "never-existed-anyway"))

	is.NoError(disk.RemoveFromSet(key, "some-item"))
	set, err = disk.GetSet(key)
	is.NoError(err)
	is.Empty(set)

	set, err = disk.GetSet("no-set-here")
	is.NoError(err)
	is.Empty(set)

	clean(is)
}

func TestDiskHelperFunctions(t *testing.T) {
	is := assert.New(t)

	set := []string{"apple", "orange", "banana"}

	bs, err := setToBytes(set)
	is.NoError(err)
	is.NotNil(bs)

	s, err := bytesToSet(bs)
	is.NoError(err)
	is.Equal(s, set)
}

func setup(is *assert.Assertions) *Disk {
	is.NoError(os.MkdirAll(dir, os.ModePerm))

	disk, err := NewDisk(dir)
	is.NoError(err)
	is.NotNil(disk)
	return disk
}

func clean(is *assert.Assertions) {
	is.NoError(os.RemoveAll(dir))
}
