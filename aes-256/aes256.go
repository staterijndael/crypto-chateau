package aes_256

func encrypt(inputBytes []rune, key []rune) ([]rune, error) {
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

	output := make([]rune, len(inputBytes))

	for row := range state {
		for col := range state[row] {
			output[row+4*col] = rune(state[row][col])
		}
	}

	return output, nil
}

func decrypt(cipher []rune, key []rune) ([]rune, error) {
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

	output := make([]rune, len(cipher))

	for row := range state {
		for col := range state[row] {
			output[row+4*col] = rune(state[row][col])
		}
	}

	return output, nil
}
