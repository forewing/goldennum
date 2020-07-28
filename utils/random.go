package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomBytes return random bytes
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RandomString returns a URL-safe,
// base64 encoded securely generated random string.
func RandomString(s int) (string, error) {
	b, err := RandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
