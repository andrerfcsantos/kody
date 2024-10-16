package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hex(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
