package dh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type KeyStore struct {
	PrivateKey *big.Int
	PublicKey  *big.Int
	SharedKey  *big.Int
}

func (k *KeyStore) GenerateSharedKey(receivedPublicKey *big.Int) error {
	sharedKey := new(big.Int)

	if !IsKeyValid(receivedPublicKey) {
		return errors.New("incorrect received public key")
	}
	if !IsKeyValid(k.PrivateKey) {
		return errors.New("incorrect private key")
	}

	sharedKey.Exp(receivedPublicKey, k.PrivateKey, Prime)

	k.SharedKey = sharedKey

	fmt.Println(sharedKey.String())

	return nil
}

func (k *KeyStore) GeneratePublicKey() error {
	if !IsKeyValid(k.PrivateKey) {
		return errors.New("incorrect private key")
	}

	publicKey := new(big.Int)
	publicKey.Exp(Generator, k.PrivateKey, Prime)

	k.PublicKey = publicKey

	return nil
}

func (k *KeyStore) GeneratePrivateKey() {
	privateKey, _ := rand.Int(rand.Reader, Prime)

	k.PrivateKey = privateKey
}

func IsValidParams(g *big.Int, p *big.Int) error {
	if g.Cmp(Generator) != 0 {
		return errors.New("incorrect generator")
	}

	if p.Cmp(Prime) != 0 {
		return errors.New("incorrect prime")
	}

	return nil
}

func IsKeyValid(key *big.Int) bool {
	return key.BitLen() != 0
}
