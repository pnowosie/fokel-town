package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

type UserData struct {
	// Id is 6-nibble hex string. I picked this for simplicity and to better illustrate the idea of a Merkle-Patricia tree.
	// In general Id would be a 32-byte hash of rest of the UserData fields.
	Id string `json:"id"`
	// FirstName is arbitrary string.
	FirstName string `json:"firstname"`
	// LastName is arbitrary string.
	LastName string `json:"lastname"`
}

// IsValid checks only Id is a valid hex string.
func (ud *UserData) IsValid() bool {
	bytes, err := hex.DecodeString(ud.Id)
	return err == nil && len(bytes) == 6/2
}

func (ud *UserData) Hash() Hash {
	dgst := Hash{}
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s:%s:%s", ud.Id, ud.FirstName, ud.LastName)))
	copy(dgst[:], hash.Sum(nil))
	return dgst
}
