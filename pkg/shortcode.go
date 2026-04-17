package pkg

import (
	"encoding/base64"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateShortCode() string {
	rand.NewSource(time.Now().UnixNano())
	code := make([]byte, 4)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return base64.RawURLEncoding.EncodeToString(code)
}
