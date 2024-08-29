package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CheckPasswordHash(password, hashedPassword string) bool {
	real_hash := HashPassword(password)
	return real_hash == hashedPassword
}
