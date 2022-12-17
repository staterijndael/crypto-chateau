package conv

import "testing"

func TestString(t *testing.T) {
	str := "Hello World!"
	b := ConvertStringToBytes(str)
	str2 := ConvertBytesToString(&BinaryIterator{Bytes: b})
	if str != str2 {
		t.Errorf("Expected %s, got %s", str, str2)
	}
}
