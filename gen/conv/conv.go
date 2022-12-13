package conv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/message"
)

func ConvFunctionMarhsalByType(t ast.Type) string {
	if t == ast.Object {
		return "ConvertObjectToBytes"
	}

	if t == ast.Uint8 {
		return "ConvertUint8ToBytes"
	}

	if t == ast.Uint32 {
		return "ConvertUint32ToBytes"
	}

	if t == ast.Uint64 {
		return "ConvertUint64ToBytes"
	}

	if t == ast.String {
		return "ConvertStringToBytes"
	}

	if t == ast.Bool {
		return "ConvertBoolToBytes"
	}

	if t == ast.Byte {
		return "ConvertByteToBytes"
	}

	if t == ast.Uint16 {
		return "ConvertUint16ToBytes"
	}

	return ""
}

func ConvFunctionUnmarshalByType(t ast.Type) string {
	if t == ast.Object {
		return "ConvertBytesToObject"
	}

	if t == ast.Uint8 {
		return "ConvertBytesToUint8"
	}

	if t == ast.Uint32 {
		return "ConvertBytesToUint32"
	}

	if t == ast.Uint64 {
		return "ConvertBytesToUint64"
	}

	if t == ast.String {
		return "ConvertBytesToString"
	}

	if t == ast.Bool {
		return "ConvertBytesToBool"
	}

	if t == ast.Byte {
		return "ConvertBytesToByte"
	}

	if t == ast.Uint16 {
		return "ConvertBytesToUint16"
	}

	return ""
}

func ConvertInt8ToBytes(num int8) []byte {
	return []byte{byte(num)}
}

func ConvertBytesToInt8(b []byte) int8 {
	return int8(b[0])
}

func ConvertInt32ToBytes(num int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(num))

	return buf
}

func ConvertBytesToInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func ConvertInt64ToBytes(num int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))

	return buf
}

func ConvertBytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func ConvertUint16ToBytes(num uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, num)

	return buf
}

func ConvertBytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func ConvertByteToBytes(b byte) []byte {
	return []byte{b}
}

func ConvertBytesToObject(msg message.Message, b []byte) {
	_, params, err := GetParams(b)
	if err != nil {
		fmt.Println(err)
	}
	err = msg.Unmarshal(params)
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertBytesToUint8(b []byte) uint8 {
	return b[0]
}

func ConvertBytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func ConvertBytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func ConvertBytesToString(b []byte) string {
	return string(b)
}

func ConvertBoolToString(b []byte) bool {
	if b[0] == '1' {
		return true
	}

	return false
}

func ConvertUint8ToBytes(num uint8) []byte {
	return []byte{num}
}

func ConvertBytesToByte(b []byte) byte {
	return b[0]
}

func ConvertUint32ToBytes(num uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, num)

	return buf
}

func ConvertUint64ToBytes(num uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, num)

	return buf
}

func ConvertStringToBytes(str string) []byte {
	return []byte(str)
}

func ConvertBoolToBytes(b bool) []byte {
	if b {
		return []byte{'1'}
	}

	return []byte{'0'}
}

func ConvertObjectToBytes(msg message.Message) []byte {
	return msg.Marshal()
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

func GetParams(p []byte) (int, map[string][]byte, error) {
	params := make(map[string][]byte)
	paramBuf := make([]byte, 0, len(p))
	valueBuf := make([]byte, 0, len(p))

	paramBufLast := -1
	valueBufLast := -1

	var paramFilled bool
	var stringParsing bool

	var openBraceCount int
	var closeBraceCount int

	var openSquareBracketCount int
	var closeSquareBracketCount int

	var isArrParsing bool

	for i, b := range p {
		if (b == ',' && paramBufLast != len(paramBuf)-1 && valueBufLast != len(valueBuf)-1 && openBraceCount == closeBraceCount+1 && (!isArrParsing || openSquareBracketCount == closeSquareBracketCount)) || (b == '}' && openBraceCount == closeBraceCount+1) {
			if b == '}' && i != len(p)-1 {
				valueBuf = append(valueBuf, b)
				closeBraceCount++
			}
			if paramBufLast == len(paramBuf)-1 || valueBufLast == len(valueBuf)-1 {
				return 0, nil, errors.New("incorrect message format: null value")
			}

			params[string(paramBuf[paramBufLast+1:])] = valueBuf[valueBufLast+1:]
			paramBufLast = len(paramBuf) - 1
			valueBufLast = len(valueBuf) - 1

			paramFilled = false
			isArrParsing = false
		} else if b == '[' {
			valueBuf = append(valueBuf, '[')
			isArrParsing = true
			openSquareBracketCount++
		} else if b == ']' {
			valueBuf = append(valueBuf, ']')
			closeSquareBracketCount++
		} else if b == '{' {
			if paramFilled {
				valueBuf = append(valueBuf, b)
			}
			openBraceCount++
		} else if b == '}' {
			valueBuf = append(valueBuf, b)
			closeBraceCount++
		} else if b == ':' && stringParsing == false && !paramFilled {
			paramFilled = true
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

	return len(p), params, nil
}
