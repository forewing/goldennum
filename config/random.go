package config

import (
	"crypto/rand"
	"encoding/base64"
)

// randomBytes return random bytes
func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// randomString returns a URL-safe,
// base64 encoded securely generated random string.
func randomString(s int) (string, error) {
	b, err := randomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
