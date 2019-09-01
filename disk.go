package strato

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

type Disk struct {
	db *badger.DB
}

var _ KV = (*Disk)(nil)

func locationToKey(location *Location) []byte {
	return []byte(fmt.Sprintf("%s__%s", location.Bucket, location.Key))
}

func NewDisk(file string) (*Disk, error) {
	db, err := badger.Open(badger.DefaultOptions(file))
	if err != nil {
		return nil, err
	}
	return &Disk{
		db: db,
	}, nil
}

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

		return tx.Set(key, val)
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