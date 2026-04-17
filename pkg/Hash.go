package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func GenerateHash(from string) string {
	salt := os.Getenv("HASH_SALT")
	h := sha256.New()
	h.Write([]byte(from + salt))
	return hex.EncodeToString(h.Sum(nil))
}
