package types

import (
	"encoding/binary"
)

var _ binary.ByteOrder

const (
	// NFTKeyPrefix is the prefix to retrieve all NFT
	NFTKeyPrefix        = "NFT/value/"
	NFTByOwnerKeyPrefix = "NFT/owner/value/"
)

// NFTKey returns the store key to retrieve an NFT from the index fields
func NFTKey(
	address string,
) []byte {
	var key []byte

	key = append(key, []byte(address)...)
	key = append(key, []byte("/")...)

	return key
}

// NFTOwnerKey returns the store key to retrieve an NFT owner from the index fields
func NFTOwnerKey(
	address string,
) []byte {
	var key []byte

	key = append(key, []byte(address)...)
	key = append(key, []byte("/")...)

	return key
}
