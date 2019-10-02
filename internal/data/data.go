package data

import (
	"encoding/binary"
	"encoding/json"
	"strconv"
)

func BytesToInt64(bs []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bs))
}

func Int64ToBytes(i int64) []byte {
	bs := make([]byte, 8)

	binary.LittleEndian.PutUint64(bs, uint64(i))

	return bs
}

func BoolAsBytes(b bool) []byte {
	return []byte(strconv.FormatBool(b))
}

func BytesToSet(bs []byte) (*Set, error) {
	var items []string

	if err := json.Unmarshal(bs, &items); err != nil {
		return nil, err
	}

	return &Set{
		items: items,
	}, nil
}
