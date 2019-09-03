package data

import (
	"encoding/binary"
	"encoding/json"
)

func BytesToInt64(bs []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bs))
}

func Int64ToBytes(i int64) []byte {
	bs := make([]byte, 8)

	binary.LittleEndian.PutUint64(bs, uint64(i))

	return bs
}

func BytesToSet(bs []byte) ([]string, error) {
	var set []string

	if err := json.Unmarshal(bs, &set); err != nil {
		return nil, err
	}

	return set, nil
}

func SetToBytes(set []string) ([]byte, error) {
	return json.Marshal(set)
}
