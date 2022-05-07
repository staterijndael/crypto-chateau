package transport

import (
	"crypto/sha256"
	"math/big"
)

func getSha256FromBigInt(num *big.Int) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(num.Bytes())
	if err != nil {
		return nil, err
	}

	sha256 := hasher.Sum(nil)

	return sha256, nil
}
