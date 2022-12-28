package conv

import (
	"encoding/binary"
	"errors"
)

var (
	ErrNotEnoughBytes = errors.New("not enough bytes")
)

type BinaryIterator struct {
	Bytes []byte
	Index int
}

func NewBinaryIterator(b []byte) *BinaryIterator {
	return &BinaryIterator{
		Bytes: b,
		Index: 0,
	}
}

func (b *BinaryIterator) NextSize() (int, error) {
	if b.Index+4 > len(b.Bytes) {
		return 0, ErrNotEnoughBytes
	}
	result := binary.BigEndian.Uint32(b.Bytes[b.Index : b.Index+4])
	b.Index += 4

	return int(result), nil
}

func (b *BinaryIterator) Slice(n int) (*BinaryIterator, error) {
	if b.Index+n > len(b.Bytes) {
		return nil, ErrNotEnoughBytes
	}
	result := &BinaryIterator{
		Bytes: b.Bytes[b.Index : b.Index+n],
		Index: 0,
	}
	b.Index += n

	return result, nil
}

func (b *BinaryIterator) HasNext() bool {
	return b.Index < len(b.Bytes)
}
