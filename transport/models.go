package transport

import (
	"errors"
	"math/big"
)

type dhParamsInitMsg struct {
	g *big.Int
	p *big.Int
}

type publicKeyInitMsg struct {
	publicKey *big.Int
}

func formatMsg(fields ...[]byte) []byte {
	if len(fields) == 0 {
		return []byte{}
	}

	result := make([]byte, 0, 1024)

	for i := 0; i < len(fields)-1; i++ {
		result = append(result, fields[i]...)
		result = append(result, '|')
	}

	result = append(result, fields[len(fields)-1]...)

	return result
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
			buf = append(buf, symb)
			lastIndex = i
		} else {
			buf = append(buf, symb)
		}
	}

	if lastIndex == len(buf)-1 {
		return nil, errors.New("incorrect message format")
	}

	result = append(result, buf[lastIndex+1:])

	return result, nil
}
