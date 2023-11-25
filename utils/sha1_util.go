package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1Data(data []byte) string {
	enc := sha1.Sum(data)
	return hex.EncodeToString(enc[:])
}
