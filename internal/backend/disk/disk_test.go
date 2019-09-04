package disk

import (
	"github.com/lucperkins/strato/internal/oops"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/services/kv"

	"github.com/lucperkins/strato/internal/data"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

func tmpDataDir() string {
	here, _ := os.Getwd()
	return filepath.Join(here, rootDataDir)
}

func TestGenericDiskFunctions(t *testing.T) {
	is := assert.New(t)

	is.NoError(os.MkdirAll(tmpDataDir(), os.ModePerm))

	db, err := badger.Open(badger.DefaultOptions(rootDataDir))
	is.NoError(err)

	key := []byte("some-key")
	value := []byte("some value")

	val, err := dbRead(db, key)
	is.True(oops.IsNotFound(err))
	is.Nil(val)

	is.NoError(dbWrite(db, key, value))

	val, err = dbRead(db, key)
	is.NoError(err)
	is.Equal(val, value)

	is.NoError(dbDelete(db, key))

	val, err = dbRead(db, key)
	is.True(oops.IsNotFound(err))
	is.Nil(val)

	clean(is)
}

func TestDiskCache(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key, value := "some-cache-key", "some-value"
	ttl := int32(3600)

	val, err := disk.CacheGet(key)
	is.True(strato.IsNotFound(err))
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

	key := "test-key"

	val := &kv.Value{
		Content: []byte("here is some test content"),
	}

	is.NoError(disk.KVPut(key, val))

	fetched, err := disk.KVGet(key)
	is.NoError(err)
	is.Equal(fetched, val)

	is.NoError(disk.KVDelete(key))

	fetched, err = disk.KVGet(key)
	is.Error(err)
	is.Equal(err.Error(), badger.ErrKeyNotFound.Error())
	is.Nil(fetched)

	clean(is)
}

func TestDiskSet(t *testing.T) {
	is := assert.New(t)

	disk := setup(is)

	key, item := "some-set", "some-item"

	set, err := disk.GetSet(key)
	is.NoError(err)
	is.Empty(set)

	set, err = disk.AddToSet(key, item)
	is.NoError(err)
	is.Len(set, 1)
	is.Equal(set[0], item)

	set, err = disk.GetSet(key)
	is.NoError(err)
	is.Len(set, 1)
	is.Equal(set[0], item)

	set, err = disk.RemoveFromSet(key, item)
	is.NoError(err)
	is.Len(set, 1)

	set, err = disk.RemoveFromSet(key, item)
	is.NoError(err)
	is.Empty(set)

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

	testCases := [][]string{
		{},
		{"apple", "orange", "banana"},
	}

	for _, tc := range testCases {
		bs, err := data.SetToBytes(tc)
		is.NoError(err)
		is.NotNil(bs)

		s, err := data.BytesToSet(bs)
		is.NoError(err)
		is.Equal(s, tc)
	}

	clean(is)
}

func setup(is *assert.Assertions) *Disk {
	clean(is)

	disk, err := NewDiskBackend()
	is.NoError(err)
	is.NotNil(disk)

	is.NoError(disk.Flush())

	return disk
}

func clean(is *assert.Assertions) {
	is.NoError(os.RemoveAll("tmp"))
}
