package aes_256

import "errors"

func CBCXor(b []byte, cPrev []byte) ([]byte, error) {
	if len(b) != Nb*4 || len(cPrev) != Nb*4 {
		return nil, errors.New("expected Nb * 4 as a input")
	}

	for i := 0; i < Nb*4; i++ {
		b[i] = b[i] ^ cPrev[i]
	}

	return b, nil
}
