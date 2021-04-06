package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func MatchPass(pass, encodePass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	sum := h.Sum(nil)
	encoded := hex.EncodeToString(sum)
	return encoded == encodePass
}
