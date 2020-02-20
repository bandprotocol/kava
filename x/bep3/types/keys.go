package types

import (
	"encoding/binary"
	"encoding/hex"

	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	// ModuleName is the name of the module
	ModuleName = "bep3"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute is the querier route for bep3
	QuerierRoute = ModuleName

	// DefaultParamspace default namestore
	DefaultParamspace = ModuleName
)

// Key prefixes
var (
	HTLTKeyPrefix       = []byte{0x00} // prefix for keys that store HTLTs
	HTLTByTimeKeyPrefix = []byte{0x01} // prefix for keys of the HTLTByTime index
)

// GetHTLTByTimeKey returns the key for iterating HTLTs by time
func GetHTLTByTimeKey(expirationTime uint64, htltID []byte) []byte {
	return append(Uint64ToBytes(expirationTime), htltID...)
}

// Uint64ToBytes converts a uint64 into fixed length bytes for use in store keys.
func Uint64ToBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(id))
	return bz
}

// Uint64FromBytes converts some fixed length bytes back into a uint64.
func Uint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// BytesToHexEncodedString converts data from []byte to a hex-encoded string
func BytesToHexEncodedString(data []byte) string {
	encodedData := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(encodedData, data)
	return string(encodedData)
}

// HexEncodedStringToBytes converts data from a hex-encoded string to []bytes
func HexEncodedStringToBytes(data string) ([]byte, error) {
	decodedData, err := hex.DecodeString(data)
	if err != nil {
		return []byte{}, err
	}
	return decodedData, nil
}

// StringToHexBytes converts data from a string to Tendermint's cmn.HexBytes
func StringToHexBytes(data string) (cmn.HexBytes, error) {
	dataRawBytes, err := hex.DecodeString(data)
	if err != nil {
		return cmn.HexBytes{}, err
	}

	return cmn.HexBytes(dataRawBytes), nil
}
