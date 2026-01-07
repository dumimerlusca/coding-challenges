package main

import (
	"crypto/md5"
	"encoding/hex"
)

func generateHash(input string, length int) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:length]
}
