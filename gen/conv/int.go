package conv

import "encoding/binary"

func ConvertBytesToInt(b *BinaryIterator) int {
	return int(binary.BigEndian.Uint64(b.Bytes))
}

func ConvertIntToBytes(num int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))

	return buf
}
