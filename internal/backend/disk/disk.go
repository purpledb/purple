package disk

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lucperkins/strato/internal/services/cache"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/lucperkins/strato/internal/services/set"

	"github.com/lucperkins/strato"

	"github.com/lucperkins/strato/internal/data"

	"github.com/dgraph-io/badger"
)

const rootDataDir = "tmp/strato"

type Disk struct {
	cache, counters, kv, sets *badger.DB
}

var (
	_ cache.Cache     = (*Disk)(nil)
	_ counter.Counter = (*Disk)(nil)
	_ kv.KV           = (*Disk)(nil)
	_ set.Set         = (*Disk)(nil)
)

func NewDiskBackend() (*Disk, error) {
	cache, err := createDb("cache")
	if err != nil {
		return nil, err
	}

	counters, err := createDb("counters")
	if err != nil {
		return nil, err
	}

	kv, err := createDb("kv")
	if err != nil {
		return nil, err
	}

	sets, err := createDb("sets")
	if err != nil {
		return nil, err
	}

	return &Disk{
		cache:    cache,
		counters: counters,
		kv:       kv,
		sets:     sets,
	}, nil
}

func createDb(dir string) (*badger.DB, error) {
	here, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(here, rootDataDir, dir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return badger.Open(badger.DefaultOptions(path))
}

// Interface methods
func (d *Disk) Close() error {
	for _, bk := range []*badger.DB{
		d.cache, d.counters, d.kv, d.sets,
	} {
		if err := bk.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Disk) Flush() error {
	for _, bk := range []*badger.DB{
		d.cache, d.counters, d.kv, d.sets,
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
			return err
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

func dbSetEntry(db *badger.DB, key, value []byte, ttl time.Duration) error {
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
		if err == badger.ErrKeyNotFound {
			return "", strato.NotFound(key)
		} else {
			return "", err
		}
	}

	return string(val), nil
}

func (d *Disk) CacheSet(key string, value string, ttl int32) error {
	k, v := []byte(key), []byte(value)

	t := time.Duration(ttl) * time.Second

	return dbSetEntry(d.cache, k, v, t)
}

// Counter
func (d *Disk) CounterGet(key string) (int64, error) {
	k := []byte(key)

	val, err := dbRead(d.counters, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return data.BytesToInt64(val), nil
}

func (d *Disk) CounterIncrement(key string, increment int64) error {
	k := []byte(key)

	val, err := dbRead(d.counters, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			v := data.Int64ToBytes(increment)

			return dbWrite(d.counters, k, v)
		} else {
			return err
		}
	}

	count := data.BytesToInt64(val)

	count += increment

	newVal := data.Int64ToBytes(count)

	return dbWrite(d.counters, k, newVal)
}

// KV
func (d *Disk) KVGet(key string) (*kv.Value, error) {
	k := []byte(key)

	val, err := dbRead(d.kv, k)
	if err != nil {
		return nil, err
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
func (d *Disk) GetSet(key string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.sets, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	return data.BytesToSet(val)
}

func (d *Disk) AddToSet(key, item string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.sets, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s := []string{item}
			value, err := data.SetToBytes(s)
			if err != nil {
				return nil, err
			}

			if err := dbWrite(d.sets, k, value); err != nil {
				return nil, err
			}

			return s, nil
		} else {
			return nil, err
		}
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return nil, err
	}

	s = append(s, item)

	value, err := data.SetToBytes(s)
	if err != nil {
		return nil, err
	}

	if err := dbWrite(d.sets, k, value); err != nil {
		return nil, err
	}

	return s, nil
}

func (d *Disk) RemoveFromSet(key, item string) ([]string, error) {
	k := []byte(key)

	val, err := dbRead(d.sets, k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return nil, err
	}

	for idx, i := range s {
		if i == item {
			s = append(s[:idx], s[idx+1:]...)
		}
	}

	value, err := data.SetToBytes(s)
	if err != nil {
		return nil, err
	}

	if err := dbWrite(d.sets, k, value); err != nil {
		return nil, err
	}

	return s, nil
}
