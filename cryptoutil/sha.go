package cryptoutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// SHA256 returns the hex-encoded SHA-256 digest of data.
func SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// SHA512 returns the hex-encoded SHA-512 digest of data.
func SHA512(data []byte) string {
	h := sha512.Sum512(data)
	return hex.EncodeToString(h[:])
}

// SHA256Bytes returns the raw SHA-256 digest bytes.
func SHA256Bytes(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// SHA512Bytes returns the raw SHA-512 digest bytes.
func SHA512Bytes(data []byte) []byte {
	h := sha512.Sum512(data)
	return h[:]
}
