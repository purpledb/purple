package data

import (
	"encoding/json"
	"strconv"
)

func BytesToInt64(bs []byte) int64 {
	return int64(bs[0])
}

func Int64ToBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
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
