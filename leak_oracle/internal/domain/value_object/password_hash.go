package value_object

import (
	"crypto/sha256"
	"encoding/hex"
)

type PasswordHash string

func NewPasswordHashFromPassword(password []byte) (PasswordHash, error) {
	sha256Hash := sha256.Sum256(password)
	return PasswordHash(hex.EncodeToString(sha256Hash[:])), nil
}
