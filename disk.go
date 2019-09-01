package strato

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"time"
)

type Disk struct {
	db *badger.DB
}

var (
	_ Cache   = (*Disk)(nil)
	_ KV      = (*Disk)(nil)
)

func NewDisk(file string) (*Disk, error) {
	db, err := badger.Open(badger.DefaultOptions(file))
	if err != nil {
		return nil, err
	}
	return &Disk{
		db: db,
	}, nil
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

// KV
func (d *Disk) KVGet(location *Location) (*Value, error) {
	var content []byte

	if err := location.validate(); err != nil {
		return nil, err
	}

	if err := d.db.View(func(tx *badger.Txn) error {
		key := locationToKey(location)

		it, err := tx.Get(key)
		if err != nil {
			return err
		}

		if err := it.Value(func(val []byte) error {
			content = val

			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Value{
		Content: content,
	}, nil
}

func (d *Disk) KVPut(location *Location, value *Value) error {
	return d.db.Update(func(tx *badger.Txn) error {
		key := locationToKey(location)
		val := value.Content

		if err := tx.Set(key, val); err != nil {
			return err
		}

		return tx.Commit()
	})
}

func (d *Disk) KVDelete(location *Location) error {
	return d.db.Update(func(tx *badger.Txn) error {
		key := locationToKey(location)

		return tx.Delete(key)
	})
}

func (d *Disk) Close() error {
	return d.db.Close()
}

// Helpers
func locationToKey(location *Location) []byte {
	return []byte(fmt.Sprintf("%s__%s", location.Bucket, location.Key))
}

func counterKey(key string) []byte {
	return []byte(fmt.Sprintf("counter__%s", key))
}

func cacheKey(key string) []byte {
	return []byte(fmt.Sprintf("cache__%s", key))
}