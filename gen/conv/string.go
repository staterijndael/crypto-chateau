package conv

func ConvertBytesToString(b *BinaryIterator) string {
	return string(b.Bytes)
}

func ConvertStringToBytes(str string) []byte {
	return []byte(str)
}
