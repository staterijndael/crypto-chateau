package aes_256

import "errors"

func keyExpansion(key []rune) ([][]uint8, error) {
	if len(key) != 4*Nk {
		return nil, errors.New("incorrect len of secret key(should be 32(4 * nk))")
	}

	keySchedule := make([][]uint8, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < Nk; c++ {
			keySchedule[r] = append(keySchedule[r], uint8(key[r+4*c]))
		}
	}

	for col := Nk; col < Nb*(Nk+1); col++ {
		if col%Nk == 0 {
			tmpPrevCol := make([]uint8, 4)
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
