package strato

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"strconv"
	"time"
)

type Disk struct {
	db *badger.DB
}

var (
	_ Cache   = (*Disk)(nil)
	_ Counter = (*Disk)(nil)
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

	return int64(val[0]), nil
}

func (d *Disk) CounterIncrement(key string, increment int64) error {
	k := counterKey(key)

	val, err := d.read(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			v := []byte{byte(increment)}

			return d.write(k, v)
		} else {
			return err
		}
	}

	count := int64(val[0])

	count += increment

	newVal := intToBytes(count)

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

func (d *Disk) Close() error {
	return d.db.Close()
}

// Set

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

func intToBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}

func bytesToSet(bs []byte) ([]string, error) {
	var set []string

	if err := json.Unmarshal(bs, &set); err != nil {
		return nil, err
	}

	return set, nil
}

func setToBytes(set []string) ([]byte, error) {
	return json.Marshal(set)
}