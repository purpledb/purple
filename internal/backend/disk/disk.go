package disk

import (
	"github.com/lucperkins/strato/internal/services/flag"
	"github.com/lucperkins/strato/internal/util"
	"os"
	"path/filepath"
	"time"

	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/services/cache"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/lucperkins/strato/internal/services/set"

	"github.com/lucperkins/strato/internal/data"

	"github.com/dgraph-io/badger"
)

const rootDataDir = "tmp/strato"

type Disk struct {
	cache, counter, flag, kv, set *badger.DB
}

func (d *Disk) Name() string {
	return "disk"
}

var (
	_ cache.Cache     = (*Disk)(nil)
	_ counter.Counter = (*Disk)(nil)
	_ flag.Flag       = (*Disk)(nil)
	_ kv.KV           = (*Disk)(nil)
	_ set.Set         = (*Disk)(nil)
)

func NewDiskBackend() (*Disk, error) {
	cacheDb, err := createDb("cache")
	if err != nil {
		return nil, err
	}

	counterDb, err := createDb("counter")
	if err != nil {
		return nil, err
	}

	flagDb, err := createDb("flag")
	if err != nil {
		return nil, err
	}

	kvDb, err := createDb("kv")
	if err != nil {
		return nil, err
	}

	setDb, err := createDb("set")
	if err != nil {
		return nil, err
	}

	return &Disk{
		cache:   cacheDb,
		counter: counterDb,
		flag:    flagDb,
		kv:      kvDb,
		set:     setDb,
	}, nil
}

func createDb(subDir string) (*badger.DB, error) {
	here, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(here, rootDataDir, subDir)

	if err := util.MkDirIfNotExists(path); err != nil {
		return nil, err
	}

	return badger.Open(badger.DefaultOptions(path))
}

// Service methods
func (d *Disk) Close() error {
	for _, bk := range []*badger.DB{
		d.cache, d.counter, d.kv, d.set,
	} {
		if err := bk.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Disk) Flush() error {
	for _, bk := range []*badger.DB{
		d.cache, d.counter, d.flag, d.kv, d.set,
	} {
		if err := bk.DropAll(); err != nil {
			return err
		}
	}

	return nil
}

// Generic functions
func dbRead(db *badger.DB, key []byte) ([]byte, error) {
	var value []byte

	if err := db.View(func(tx *badger.Txn) error {
		it, err := tx.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return strato.NotFound(string(key))
			} else {
				return err
			}
		}

		val, err := it.ValueCopy(nil)
		if err != nil {
			return err
		}

		value = val

		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}

func dbWrite(db *badger.DB, key, value []byte) error {
	return db.Update(func(tx *badger.Txn) error {
		return tx.Set(key, value)
	})
}

func dbDelete(db *badger.DB, key []byte) error {
	return db.Update(func(tx *badger.Txn) error {
		return tx.Delete(key)
	})
}

func dbSetCacheEntry(db *badger.DB, key, value []byte, ttl time.Duration) error {
	entry := badger.NewEntry(key, value).WithTTL(ttl)

	return db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(entry)
	})
}

// Cache
func (d *Disk) CacheGet(key string) (string, error) {
	k := []byte(key)

	val, err := dbRead(d.cache, k)
	if err != nil {
		if strato.IsNotFound(err) {
			return "", strato.NotFound(key)
		} else {
			return "", err
		}
	}

	return string(val), nil
}

func (d *Disk) CacheSet(key string, value string, ttl int32) error {
	if key == "" {
		return strato.ErrNoKey
	}

	if value == "" {
		return strato.ErrNoValue
	}

	k, v := []byte(key), []byte(value)

	t := time.Duration(ttl) * time.Second

	return dbSetCacheEntry(d.cache, k, v, t)
}

// Counter
func (d *Disk) CounterGet(key string) (int64, error) {
	k := []byte(key)

	val, err := dbRead(d.counter, k)
	if err != nil {
		if strato.IsNotFound(err) {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return data.BytesToInt64(val), nil
}

func (d *Disk) CounterIncrement(key string, increment int64) error {
	k := []byte(key)

	val, err := dbRead(d.counter, k)
	if err != nil {
		if strato.IsNotFound(err) {
			v := data.Int64ToBytes(increment)

			return dbWrite(d.counter, k, v)
		} else {
			return err
		}
	}

	count := data.BytesToInt64(val)

	count += increment

	newVal := data.Int64ToBytes(count)

	return dbWrite(d.counter, k, newVal)
}

// Flag
func (d *Disk) FlagGet(key string) (bool, error) {
	k := []byte(key)

	val, err := dbRead(d.flag, k)
	if err != nil {
		if strato.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return data.BoolFromBytes(val)
}

func (d *Disk) FlagSet(key string, value bool) error {
	k := []byte(key)

	val := data.BoolAsBytes(value)

	return dbWrite(d.flag, k, val)
}

// KV
func (d *Disk) KVGet(key string) (*kv.Value, error) {
	k := []byte(key)

	val, err := dbRead(d.kv, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, strato.NotFound(key)
		} else {
			return nil, err
		}
	}

	return &kv.Value{
		Content: val,
	}, nil
}

func (d *Disk) KVPut(key string, value *kv.Value) error {
	k := []byte(key)

	return dbWrite(d.kv, k, value.Content)
}

func (d *Disk) KVDelete(key string) error {
	k := []byte(key)
	return dbDelete(d.kv, k)
}

// Set
func (d *Disk) SetGet(key string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.set, k)
	if err != nil {
		return nil, err
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return nil, err
	}

	return s.Get(), nil
}

func (d *Disk) SetAdd(key, item string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.set, k)
	if err != nil {
		if strato.IsNotFound(err) {
			s := data.NewSet(item)
			value, err := s.AsBytes()
			if err != nil {
				return nil, err
			}

			if err := dbWrite(d.set, k, value); err != nil {
				return nil, err
			}

			return s.Get(), nil
		} else {
			return nil, err
		}
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return nil, err
	}

	s.Add(item)

	value, err := s.AsBytes()
	if err != nil {
		return nil, err
	}

	if err := dbWrite(d.set, k, value); err != nil {
		return nil, err
	}

	return s.Get(), nil
}

func (d *Disk) SetRemove(key, item string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.set, k)
	if err != nil {
		return nil, err
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return nil, err
	}

	s.Remove(item)

	value, err := s.AsBytes()
	if err != nil {
		return nil, err
	}

	if err := dbWrite(d.set, k, value); err != nil {

		return nil, err
	}

	return s.Get(), nil
}
