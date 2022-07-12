package message

import (
	"errors"
)

type Message interface {
	Marshal() []byte
	Unmarshal(map[string][]byte) error
}

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

func GetParams(p []byte) (map[string][]byte, error) {
	params := make(map[string][]byte)
	paramBuf := make([]byte, 0, len(p))
	valueBuf := make([]byte, 0, len(p))

	paramBufLast := -1
	valueBufLast := -1

	var paramFilled bool
	var stringParsing bool

	for i, b := range p {
		if b == ',' || i == len(p)-1 {
			if (i != len(p)-1) && (p[i+1] == ',') {
				continue
			}

			if i == len(p)-1 {
				valueBuf = append(valueBuf, b)
			}

			if paramBufLast == len(paramBuf)-1 || valueBufLast == len(valueBuf)-1 {
				return nil, errors.New("incorrect message format: null value")
			}

			params[string(paramBuf[paramBufLast+1:])] = valueBuf[valueBufLast+1:]
			paramBufLast = len(paramBuf) - 1
			valueBufLast = len(valueBuf) - 1

			paramFilled = false
		} else if b == ':' && stringParsing == false {
			paramFilled = true
		} else if (b == ' ' && stringParsing == false) || b == '(' || b == ')' {
			continue
		} else if b == '"' {
			stringParsing = !stringParsing
		} else {
			if !paramFilled {
				paramBuf = append(paramBuf, b)
			} else {
				valueBuf = append(valueBuf, b)
			}
		}
	}

	return params, nil
}

func ParseArray(value []byte) ([][]byte, error) {
	if len(value) < 2 || value[0] != '{' || value[len(value)-1] != '}' {
		return nil, errors.New("incorrect array value")
	}

	values := make([][]byte, 0, len(value))

	valueBuf := make([]byte, 0, len(value))
	valueBufLast := -1

	var startedObjectParsing bool
	var finishedObjectParsing bool

	for _, s := range value {
		switch s {
		case '{':
			continue
		case '}':
			values = append(values, valueBuf[valueBufLast+1:])
		case '(':
			startedObjectParsing = true
		case ')':
			finishedObjectParsing = true
		case ',':
			if startedObjectParsing && finishedObjectParsing {
				values = append(values, valueBuf[valueBufLast+1:])
				valueBufLast = len(valueBuf) - 1

				startedObjectParsing = false
				finishedObjectParsing = false
			} else {
				valueBuf = append(valueBuf, s)
			}
		default:
			valueBuf = append(valueBuf, s)
		}
	}

	return values, nil
}
