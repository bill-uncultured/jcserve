// Package jchash implements a password hashing algorithm
package jchash

import (
	"crypto/sha512"
	"encoding/base64"
)

// HashPassword returns a base64 encoded sha512 hash of the string passed to it.
func HashPassword(password string) string {
	pwHash := sha512.Sum512([]byte(password))
	return base64.StdEncoding.EncodeToString(pwHash[:])
}
