package dh

import "golang.org/x/crypto/curve25519"

func DH(priv [32]byte, pub [32]byte) ([32]byte, error) {
	sharedKey, err := curve25519.X25519(priv[:], pub[:])
	if err != nil {
		return [32]byte{}, err
	}

	var buf [32]byte
	copy(buf[:], sharedKey)
	return buf, nil
}
