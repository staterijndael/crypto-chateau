package aes_256

import (
	"errors"
)

func subBytes(state [][]uint16) [][]uint16 {
	for i, row := range state {
		for j, elem := range row {
			sboxElem := Sbox[elem]
			state[i][j] = sboxElem
		}
	}

	return state
}

func keyExpansion(key []rune) ([][]uint16, error) {
	if len(key) != 4*Nk {
		return nil, errors.New("incorrect len of secret key(should be 32(4 * nk))")
	}

	keySchedule := make([][]uint16, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < Nk; c++ {
			keySchedule[r] = append(keySchedule[r], uint16(key[r+4*c]))
		}
	}

	for col := Nk; col < Nb*(Nk+1); col++ {
		if col%Nk == 0 {
			tmpPrevCol := make([]uint16, 4)
			for row := 0; row < 4; row++ {
				tmpPrevCol[row] = keySchedule[row][col-1]
			}

			tmpPrevCol = append(tmpPrevCol, keySchedule[0][col-1])

			for i, val := range tmpPrevCol {
				sboxElem := Sbox[val]
				tmpPrevCol[i] = sboxElem
			}

			for row := 0; row < 4; row++ {
				s := keySchedule[row][col-4] ^ tmpPrevCol[row] ^ Rcon[row][col/Nk-1]
				keySchedule[row] = append(keySchedule[row], s)
			}
		} else {
			for row := 0; row < 4; row++ {
				s := keySchedule[row][col-4] ^ keySchedule[row][col-1]
				keySchedule[row] = append(keySchedule[row], s)
			}
		}
	}

	return keySchedule, nil
}

func addRoundKey(state [][]uint16, keySchedule [][]uint16, round uint16) [][]uint16 {
	var (
		col uint16
		row uint16
	)

	for col = 0; col < Nk; col++ {
		for row = 0; row < Nb; row++ {
			s := state[row][col] ^ keySchedule[row][Nk*round+col]

			state[row][col] = s
		}
	}

	return state
}

func mixColumns(state [][]uint16) [][]uint16 {
	for row := 0; row < Nk; row++ {
		s0 := mulBy02(state[0][row]) ^ mulBy03(state[1][row]) ^ state[2][row] ^ state[3][row]
		s1 := state[0][row] ^ mulBy02(state[1][row]) ^ mulBy03(state[2][row]) ^ state[3][row]
		s2 := state[0][row] ^ state[1][row] ^ mulBy02(state[2][row]) ^ mulBy03(state[3][row])
		s3 := mulBy03(state[0][row]) ^ state[1][row] ^ state[2][row] ^ mulBy02(state[3][row])

		state[0][row] = s0
		state[1][row] = s1
		state[2][row] = s2
		state[3][row] = s3
	}

	return state
}

func shiftRows(state [][]uint16) [][]uint16 {
	for row := 1; row < Nb; row++ {
		res := make([]uint16, Nk)
		for col := 0; col < Nk; col++ {
			shift := (Nk - 1 - col - row) % Nk
			if shift < 0 {
				shift = Nk + shift
			}
			res[shift] = state[row][Nk-1-col]
		}

		state[row] = res
	}

	return state
}
