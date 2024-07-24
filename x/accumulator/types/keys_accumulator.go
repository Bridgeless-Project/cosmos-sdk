package types

import (
	"encoding/binary"
)

var _ binary.ByteOrder

const (
	// AdminKeyPrefix is the prefix to retrieve all admin
	AdminKeyPrefix = "admin/value/"
)

// AdminKey returns the store key to retrieve an admin from the index fields
func AdminKey(
	address string,
) []byte {
	var key []byte

	key = append(key, []byte(address)...)
	key = append(key, []byte("/")...)

	return key
}
