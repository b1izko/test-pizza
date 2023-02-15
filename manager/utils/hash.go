package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/b1izko/test-pizza/internal/logger"
)

// Hash ...
func Hash(some string) string {
	data := []byte(some + "pe0C.^(Do=utam50DiHygI0pZk)bx=#")
	hash := sha256.New()
	_, err := hash.Write(data)
	if logger.IsError(err, "Password encryption failed") {
		return ""
	}
	hashedPwd := hex.EncodeToString(hash.Sum(nil))
	return hashedPwd
}
