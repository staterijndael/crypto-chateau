package conv

func ConvertBytesToUint8(b *BinaryIterator) uint8 {
	return b.Bytes[0]
}

func ConvertUint8ToBytes(num uint8) []byte {
	return []byte{num}
}

func ConvertBytesToInt8(b *BinaryIterator) int8 {
	return int8(b.Bytes[0])
}

func ConvertInt8ToBytes(num int8) []byte {
	return []byte{byte(num)}
}

func ConvertBytesToByte(b *BinaryIterator) byte {
	return b.Bytes[0]
}

func ConvertByteToBytes(num byte) []byte {
	return []byte{num}
}
