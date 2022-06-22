package crypto_chateau

import (
	"errors"
	"reflect"
	"strings"
)

type Message interface {
	Marshal() ([]byte, error)
}

func ParseMessage(p []byte, msgType Message) (Message, error) {
	var msg Message
	params, err := getParams(p)
	if err != nil {
		return nil, err
	}

	msg, err = generateMessage(params, msgType)
	if err != nil {
		return nil, err
	}

	return msg, nil
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

func generateMessage(params map[string][]byte, msgType Message) (Message, error) {
	elem := reflect.ValueOf(msgType).Elem()

	for param, value := range params {
		f := elem.FieldByName(strings.Title(param))

		if f.IsValid() {
			switch f.Kind() {
			// todo: can be optimized with monomorphism
			// etc binary.BigEndian not suitable
			case reflect.Int8:
				var x int8
				x |= int8(value[0])
				f.SetInt(int64(x))
			case reflect.Int16:
				var x int16
				x |= int16(value[1])
				x |= int16(value[0]) << 8
				f.SetInt(int64(x))
			case reflect.Int32:
				var x int32
				x |= int32(value[3])
				x |= int32(value[2]) << 8
				x |= int32(value[1]) << 16
				x |= int32(value[0]) << 24
				f.SetInt(int64(x))
			case reflect.Int64:
				var x int64
				x |= int64(value[7])
				x |= int64(value[6]) << 8
				x |= int64(value[5]) << 16
				x |= int64(value[4]) << 24
				x |= int64(value[3]) << 32
				x |= int64(value[2]) << 40
				x |= int64(value[1]) << 48
				x |= int64(value[0]) << 56
				f.SetInt(x)
			case reflect.Int:
				// exact same as int64
				var x int
				x |= int(value[7])
				x |= int(value[6]) << 8
				x |= int(value[5]) << 16
				x |= int(value[4]) << 24
				x |= int(value[3]) << 32
				x |= int(value[2]) << 40
				x |= int(value[1]) << 48
				x |= int(value[0]) << 56
				f.SetInt(int64(x))
			case reflect.Uint:
				var x uint
				x |= uint(value[7])
				x |= uint(value[6]) << 8
				x |= uint(value[5]) << 16
				x |= uint(value[4]) << 24
				x |= uint(value[3]) << 32
				x |= uint(value[2]) << 40
				x |= uint(value[1]) << 48
				x |= uint(value[0]) << 56
				f.SetUint(uint64(x))
			case reflect.Uint16:
				var x uint16
				x |= uint16(value[1])
				x |= uint16(value[0]) << 8
				f.SetUint(uint64(x))
			case reflect.Uint32:
				var x uint32
				x |= uint32(value[3])
				x |= uint32(value[2]) << 8
				x |= uint32(value[1]) << 16
				x |= uint32(value[0]) << 24
				f.SetUint(uint64(x))
			case reflect.Uint64:
				var x uint64
				x |= uint64(value[7])
				x |= uint64(value[6]) << 8
				x |= uint64(value[5]) << 16
				x |= uint64(value[4]) << 24
				x |= uint64(value[3]) << 32
				x |= uint64(value[2]) << 40
				x |= uint64(value[1]) << 48
				x |= uint64(value[0]) << 56
				f.SetUint(x)
			case reflect.String:
				f.SetString(string(value))
			case reflect.Bool:
				var x bool
				if value[0] == '1' {
					x = true
				} else if value[0] == '0' {
					x = false
				} else {
					return nil, errors.New("incorrect value for bool type " + string(value))
				}
				f.SetBool(x)
			default:
				return nil, errors.New("unknown field type")
			}
		} else {
			return nil, errors.New("unknown param " + param)
		}
	}

	return msgType, nil
}

func getParams(p []byte) (map[string][]byte, error) {
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

			if len(paramBuf) != 0 && len(valueBuf) != 0 {
				params[string(paramBuf[paramBufLast+1:])] = valueBuf[valueBufLast+1:]
				paramBufLast = len(paramBuf) - 1
				valueBufLast = len(valueBuf) - 1

				paramFilled = false
			}
		} else if b == ':' && stringParsing == false {
			paramFilled = true
		} else if b == ' ' && stringParsing == false {
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

	if len(paramBuf) != 0 || len(valueBuf) != 0 {
		return nil, errors.New("incorrect message format")
	}

	return params, nil
}
