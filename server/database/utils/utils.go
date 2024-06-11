package utils

import (
	"crypto/rand"
	"encoding/hex"
)

const SESSION_SIZE = 32

// generate secure 32 byte token
func GenerateSecureToken() (string, error) {
    var b [SESSION_SIZE]byte
    _, err := rand.Read(b[:])
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b[:]), nil
}
