package conv

import "encoding/binary"

var (
	ObjectBytesPrefixLength = len(ConvertSizeToBytes(0))
)

func ConvertSizeToBytes(num int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(num))

	return buf
}
