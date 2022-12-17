package conv

import "encoding/binary"

func ConvertBytesToUint32(b *BinaryIterator) uint32 {
	return binary.BigEndian.Uint32(b.Bytes)
}

func ConvertUint32ToBytes(num uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, num)

	return buf
}

func ConvertBytesToInt32(b *BinaryIterator) int32 {
	return int32(binary.BigEndian.Uint32(b.Bytes))
}

func ConvertInt32ToBytes(num int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(num))

	return buf
}
