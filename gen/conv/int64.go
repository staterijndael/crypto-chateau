package conv

import "encoding/binary"

func ConvertBytesToUint64(b *BinaryIterator) uint64 {
	return binary.BigEndian.Uint64(b.Bytes)
}

func ConvertUint64ToBytes(num uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, num)

	return buf
}

func ConvertBytesToInt64(b *BinaryIterator) int64 {
	return int64(binary.BigEndian.Uint64(b.Bytes))
}

func ConvertInt64ToBytes(num int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))

	return buf
}
