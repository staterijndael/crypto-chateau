package conv

import "encoding/binary"

func ConvertBytesToInt(b *BinaryIterator) uint64 {
	return binary.BigEndian.Uint64(b.Bytes)
}

func ConvertIntToBytes(num uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, num)

	return buf
}
