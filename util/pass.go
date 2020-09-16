package util

import (
	"crypto/sha256"
	"encoding/base64"
)

func Hash(val string) string {
	passHash := sha256.Sum256([]byte(val))
	return base64.StdEncoding.EncodeToString(passHash[:])
}

func HashValid(pass, hash string) bool {
	if hash == Hash(pass) {
		return true
	}
	return false
}
