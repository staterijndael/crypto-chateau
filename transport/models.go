package transport

import (
	"errors"
)

type publicKeyInitMsg struct {
	publicKey [32]byte
}

func parseMsg(msg []byte) ([][]byte, error) {
	if len(msg) == 0 {
		return nil, errors.New("empty message")
	}
	result := make([][]byte, 0, 5)

	buf := make([]byte, 0, 1024)
	lastIndex := -1

	for i, symb := range msg {
		if symb == '|' {
			if lastIndex+1 == i {
				return nil, errors.New("incorrect message format")
			}

			result = append(result, buf[lastIndex+1:])
			lastIndex = i
		}

		buf = append(buf, symb)
	}

	if lastIndex == len(buf)-1 {
		return nil, errors.New("incorrect message format")
	}

	result = append(result, buf[lastIndex+1:])

	return result, nil
}
