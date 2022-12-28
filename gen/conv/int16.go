package conv

import "encoding/binary"

func ConvertBytesToUint16(b *BinaryIterator) uint16 {
	return binary.BigEndian.Uint16(b.Bytes)
}

func ConvertUint16ToBytes(num uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, num)

	return buf
}

func ConvertBytesToInt16(b *BinaryIterator) int16 {
	return int16(binary.BigEndian.Uint16(b.Bytes))
}

func ConvertInt16ToBytes(num int16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(num))

	return buf
}
