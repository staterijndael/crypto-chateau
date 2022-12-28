package conv

func ConvertBytesToBool(b *BinaryIterator) bool {
	return b.Bytes[0] == 0x01
}

func ConvertBoolToBytes(b bool) []byte {
	if b {
		return []byte{0x01}
	}

	return []byte{0x00}
}
