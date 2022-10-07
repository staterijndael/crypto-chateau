package transport

import (
	"crypto/sha256"
)

func getSha256FromBytes(bytes [32]byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(bytes[:])
	if err != nil {
		return nil, err
	}

	sha256 := hasher.Sum(nil)

	return sha256, nil
}
