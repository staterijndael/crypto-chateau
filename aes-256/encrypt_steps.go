package aes_256

import (
	"errors"
)

func subBytes(state [][]uint16) [][]uint16 {
	for i := range state {
		for j := range state[i] {
			row := state[i][j] / 0x10
			col := state[i][j] % 0x10

			sboxElem := Sbox[16*row+col]
			state[i][j] = sboxElem
		}
	}

	return state
}

func keyExpansion(key []byte) ([][]uint16, error) {
	if len(key) != 4*Nk {
		return nil, errors.New("incorrect len of secret key(should be 32(4 * nk))")
	}

	keySchedule := make([][]uint16, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < Nk; c++ {
			keySchedule[r] = append(keySchedule[r], uint16(key[r+4*c]))
		}
	}

	for col := Nk; col < Nb*(Nr+1); col++ {
		if col%Nk == 0 {
			tmpPrevCol := make([]uint16, 4)
			for row := 1; row < 4; row++ {
				tmpPrevCol[row-1] = keySchedule[row][col-1]
			}

			tmpPrevCol[len(tmpPrevCol)-1] = keySchedule[0][col-1]

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

	for col := 0; col < Nb; col++ {
		for row := 0; row < Nb; row++ {
			s := state[row][col] ^ keySchedule[row][Nb*round+uint16(col)]

			state[row][col] = s
		}
	}

	return state
}

func mixColumns(state [][]uint16) [][]uint16 {
	for row := 0; row < Nb; row++ {
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
		res := make([]uint16, 4)
		for col := 0; col < 4; col++ {
			shift := (4 - 1 - col - row) % 4
			if shift < 0 {
				shift = 4 + shift
			}
			res[shift] = state[row][4-1-col]
		}

		state[row] = res
	}

	return state
}
