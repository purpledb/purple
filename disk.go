package strato

import (
	"fmt"
	"github.com/lucperkins/strato/internal/data"
	"os"
	"time"

	"github.com/dgraph-io/badger"
)

const dataDir = "tmp/strato"

type Disk struct {
	db      *badger.DB
}

var (
	_ Cache   = (*Disk)(nil)
	_ Counter = (*Disk)(nil)
	_ KV      = (*Disk)(nil)
	_ Set     = (*Disk)(nil)
)

func NewDisk() (*Disk, error) {
	if err := createDataDir(dataDir); err != nil {
		return nil, err
	}

	db, err := badger.Open(badger.DefaultOptions(dataDir))
	if err != nil {
		return nil, err
	}

	return &Disk{
		db:      db,
	}, nil
}

func createDataDir(dataDir string) error {
	return os.MkdirAll(dataDir, os.ModePerm)
}

// Backend methods
func (d *Disk) Close() error {
	return d.db.Close()
}

func (d *Disk) Flush() error {
	return d.db.DropAll()
}

// Generic functions
func (d *Disk) read(key []byte) ([]byte, error) {
	var value []byte

	if err := d.db.View(func(tx *badger.Txn) error {
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

func (d *Disk) write(key, value []byte) error {
	return d.db.Update(func(tx *badger.Txn) error {
		return tx.Set(key, value)
	})
}

func (d *Disk) delete(key []byte) error {
	return d.db.Update(func(tx *badger.Txn) error {
		return tx.Delete(key)
	})
}

func (d *Disk) setEntry(key, value []byte, ttl time.Duration) error {
	entry := badger.NewEntry(key, value).WithTTL(ttl)

	return d.db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(entry)
	})
}

// Cache
func (d *Disk) CacheGet(key string) (string, error) {
	k := cacheKey(key)
	val, err := d.read(k)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (d *Disk) CacheSet(key string, value string, ttl int32) error {
	k, v := cacheKey(key), []byte(value)

	t := time.Duration(ttl) * time.Second

	return d.setEntry(k, v, t)
}

// Counter
func (d *Disk) CounterGet(key string) (int64, error) {
	k := counterKey(key)

	val, err := d.read(k)
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
	k := counterKey(key)

	val, err := d.read(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			v := data.Int64ToBytes(increment)

			return d.write(k, v)
		} else {
			return err
		}
	}

	count := data.BytesToInt64(val)

	count += increment

	newVal := data.Int64ToBytes(count)

	return d.write(k, newVal)
}

// KV
func (d *Disk) KVGet(location *Location) (*Value, error) {
	key := kvKey(location)

	val, err := d.read(key)
	if err != nil {
		return nil, err
	}

	return &Value{
		Content: val,
	}, nil
}

func (d *Disk) KVPut(location *Location, value *Value) error {
	key := kvKey(location)
	return d.write(key, value.Content)
}

func (d *Disk) KVDelete(location *Location) error {
	key := kvKey(location)
	return d.delete(key)
}

// Set
func (d *Disk) GetSet(key string) ([]string, error) {
	k := setKey(key)

	val, err := d.read(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	return data.BytesToSet(val)
}

func (d *Disk) AddToSet(key, item string) error {
	k := setKey(key)

	val, err := d.read(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s := []string{item}
			value, err := data.SetToBytes(s)
			if err != nil {
				return err
			}
			return d.write(k, value)
		} else {
			return err
		}
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return err
	}

	s = append(s, item)

	value, err := data.SetToBytes(s)
	if err != nil {
		return err
	}

	return d.write(k, value)
}

func (d *Disk) RemoveFromSet(key, item string) error {
	k := setKey(key)

	val, err := d.read(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrNoSet
		} else {
			return err
		}
	}

	s, err := data.BytesToSet(val)
	if err != nil {
		return err
	}

	for idx, i := range s {
		if i == item {
			s = append(s[:idx], s[idx+1:]...)
		}
	}

	value, err := data.SetToBytes(s)
	if err != nil {
		return err
	}

	return d.write(k, value)
}

// Helpers
func cacheKey(key string) []byte {
	return []byte(fmt.Sprintf("cache__%s", key))
}

func counterKey(key string) []byte {
	return []byte(fmt.Sprintf("counter__%s", key))
}

func kvKey(location *Location) []byte {
	return []byte(fmt.Sprintf("%s__%s", location.Bucket, location.Key))
}

func setKey(key string) []byte {
	return []byte(fmt.Sprintf("set__%s", key))
}
