package conv

import (
	"bytes"
	"errors"
)

func GetHandlerName(p []byte) ([]byte, int, error) {
	buf := make([]byte, 0, 50)
	for i, b := range p {
		if b == '#' {
			return buf, i + 1, nil
		}

		buf = append(buf, b)
	}

	return nil, 0, errors.New("incorrect message format: handler name not found")
}

func GetArray(p []byte) (int, [][]byte, error) {
	if len(p) == 0 {
		return 0, nil, errors.New("array is zero length")
	}

	if p[0] != '[' {
		return 0, nil, errors.New("expected open brace")
	}

	openSquareBracketCount := 1
	var closeSquareBracketCount int

	i := 1
	for openSquareBracketCount != closeSquareBracketCount && i < len(p) {
		if p[i] == '[' {
			openSquareBracketCount++
		}

		if p[i] == ']' {
			closeSquareBracketCount++
		}

		i++
	}

	if openSquareBracketCount != closeSquareBracketCount {
		return 0, nil, errors.New("expected end of array")
	}

	values := bytes.Split(p[1:i], []byte(","))
	for i, value := range values {
		values[i] = bytes.TrimSpace(value)
	}

	return i, values, nil
}
