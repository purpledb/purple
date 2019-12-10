package data

import (
	"encoding/binary"
	"encoding/json"
	"strconv"
)

func BoolAsBytes(b bool) []byte {
	return []byte(strconv.FormatBool(b))
}

func BoolFromBytes(bs []byte) (bool, error) {
	return strconv.ParseBool(string(bs))
}

func BytesToInt64(bs []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bs))
}

func Int64ToBytes(i int64) []byte {
	bs := make([]byte, 8)

	binary.LittleEndian.PutUint64(bs, uint64(i))

	return bs
}

func BytesToSet(bs []byte) (*Set, error) {
	items := make([]string, 0)

	if err := json.Unmarshal(bs, &items); err != nil {
		return nil, err
	}

	return &Set{
		items: items,
	}, nil
}
