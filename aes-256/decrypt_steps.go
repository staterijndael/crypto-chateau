package aes_256

func InvShiftRows(state [][]uint16) [][]uint16 {
	for row := 1; row < Nb; row++ {
		res := make([]uint16, Nk)
		for col := 0; col < Nk; col++ {
			res[(col+row)%Nk] = state[row][col]
		}

		state[row] = res
	}

	return state
}

func InvMixColumns(state [][]uint16) [][]uint16 {
	for row := 0; row < Nk; row++ {
		s0 := mulBy0e(state[0][row]) ^ mulBy0b(state[1][row]) ^ mulBy0d(state[2][row]) ^ mulBy09(state[3][row])
		s1 := mulBy09(state[0][row]) ^ mulBy0e(state[1][row]) ^ mulBy0b(state[2][row]) ^ mulBy0d(state[3][row])
		s2 := mulBy0d(state[0][row]) ^ mulBy09(state[1][row]) ^ mulBy0e(state[2][row]) ^ mulBy0b(state[3][row])
		s3 := mulBy0b(state[0][row]) ^ mulBy0d(state[1][row]) ^ mulBy09(state[2][row]) ^ mulBy0e(state[3][row])

		state[0][row] = s0
		state[1][row] = s1
		state[2][row] = s2
		state[3][row] = s3
	}

	return state
}

func InvSubBytes(state [][]uint16) [][]uint16 {
	for i, row := range state {
		for j, elem := range row {
			sboxElem := InvSbox[elem]
			state[i][j] = sboxElem
		}
	}

	return state
}
