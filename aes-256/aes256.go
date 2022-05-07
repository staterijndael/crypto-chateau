package aes_256

import (
	"errors"
)

func Encrypt(inputBytes []byte, key []byte) ([]byte, error) {
	if len(inputBytes) == 0 {
		return nil, errors.New("incorrect input bytes length")
	}
	if len(key) != 4*Nk {
		return nil, errors.New("incorrect len of secret key(should be 32(4 * nk))")
	}

	for len(inputBytes)%(Nb*4) != 0 {
		inputBytes = append(inputBytes, ' ')
	}

	result := make([]byte, 0, len(inputBytes))

	for batch := 1; batch <= len(inputBytes)/(Nb*4); batch++ {
		offset := (batch - 1) * Nb * 4
		limit := offset + Nb*4

		state := inputBytes[offset:limit]

		encryptedData, err := encrypt(state, key)
		if err != nil {
			return nil, err
		}

		result = append(result, encryptedData...)
	}

	return result, nil
}

func Decrypt(cipher []byte, key []byte) ([]byte, error) {
	if len(cipher) == 0 || len(cipher)%Nb*4 != 0 {
		return nil, errors.New("incorrect input bytes length")
	}
	if len(key) != 4*Nk {
		return nil, errors.New("incorrect len of secret key(should be 32(4 * nk))")
	}

	result := make([]byte, 0, len(cipher))

	for batch := 1; batch <= len(cipher)/(Nb*4); batch++ {
		offset := (batch - 1) * Nb * 4
		limit := offset + Nb*4

		state := cipher[offset:limit]

		decryptedData, err := decrypt(state, key)
		if err != nil {
			return nil, err
		}

		result = append(result, decryptedData...)
	}

	finalIndex := len(result) - 1
	for i := len(result) - 1; i >= 0; i-- {
		if result[i] != ' ' {
			finalIndex = i + 1
			break
		}
	}

	result = result[:finalIndex]

	return result, nil
}

func encrypt(inputBytes []byte, key []byte) ([]byte, error) {
	if len(inputBytes) != 4*Nb {
		return nil, errors.New("incorrect input bytes length")
	}

	state := make([][]uint16, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < Nb; c++ {
			state[r] = append(state[r], uint16(inputBytes[r+4*c]))
		}
	}

	keySchedule, err := keyExpansion(key)
	if err != nil {
		return nil, err
	}

	state = addRoundKey(state, keySchedule, 0)

	var rnd uint16 = 1

	for ; rnd < Nr; rnd++ {
		state = subBytes(state)
		state = shiftRows(state)
		state = mixColumns(state)
		state = addRoundKey(state, keySchedule, rnd)
	}

	state = subBytes(state)
	state = shiftRows(state)
	state = addRoundKey(state, keySchedule, rnd)

	output := make([]byte, len(inputBytes))

	for row := range state {
		for col := range state[row] {
			output[row+4*col] = byte(state[row][col])
		}
	}

	return output, nil
}

func decrypt(cipher []byte, key []byte) ([]byte, error) {
	state := make([][]uint16, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < Nb; c++ {
			state[r] = append(state[r], uint16(cipher[r+4*c]))
		}
	}

	keySchedule, err := keyExpansion(key)
	if err != nil {
		return nil, err
	}

	state = addRoundKey(state, keySchedule, Nr)

	var rnd uint16 = Nr - 1

	for ; rnd > 0; rnd-- {
		state = InvShiftRows(state)
		state = InvSubBytes(state)
		state = addRoundKey(state, keySchedule, rnd)
		state = InvMixColumns(state)
	}

	state = InvShiftRows(state)
	state = InvSubBytes(state)
	state = addRoundKey(state, keySchedule, rnd)

	output := make([]byte, len(cipher))

	for row := range state {
		for col := range state[row] {
			output[row+4*col] = byte(state[row][col])
		}
	}

	return output, nil
}
