package snippet

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID returns a random 8-byte hex string suitable for snippet IDs.
func generateID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic("stashctl: failed to generate snippet ID: " + err.Error())
	}
	return hex.EncodeToString(b)
}
