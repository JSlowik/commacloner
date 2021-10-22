package threecommas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSignature computes the HMAC signature based on path and api secret
func ComputeSignature(path, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(path))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
