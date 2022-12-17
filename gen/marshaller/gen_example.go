package marshaller

import (
	"github.com/pkg/errors"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/message"
)

type CommonObject struct {
	MagicString string
}

func (o *CommonObject) Marshal() []byte {
	b := make([]byte, 0, 1024)
	// MagicString
	b = append(b, conv.ConvertSizeToBytes(len([]byte(o.MagicString)))...) // size
	b = append(b, conv.ConvertStringToBytes(o.MagicString)...)
	return b
}

func (o *CommonObject) Unmarshal(b *conv.BinaryIterator) error {
	var err error

	// MagicString
	size, err := b.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read MagicString size")
	}
	stringData, err := b.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicString")
	}
	o.MagicString = conv.ConvertBytesToString(stringData)

	return nil
}

var _ message.Message = (*CommonObject)(nil)

type ReverseRequest struct {
	MagicString     string
	MagicInt64      int64
	MagicBool       bool
	MagicBytes      []byte
	MagicObjectList []CommonObject
}

var _ message.Message = (*ReverseRequest)(nil)

func (o *ReverseRequest) Marshal() []byte {
	b := make([]byte, 0, 1024)
	// MagicString
	b = append(b, conv.ConvertSizeToBytes(len([]byte(o.MagicString)))...) // size
	b = append(b, conv.ConvertStringToBytes(o.MagicString)...)
	// MagicInt64
	b = append(b, conv.ConvertInt64ToBytes(o.MagicInt64)...)
	// MagicBool
	b = append(b, conv.ConvertBoolToBytes(o.MagicBool)...)
	// MagicBytes
	b = append(b, conv.ConvertSizeToBytes(len(o.MagicBytes))...) // size
	b = append(b, o.MagicBytes...)
	// MagicObjectList
	objectBuf := make([]byte, 0, 1024)
	for _, object := range o.MagicObjectList {
		// marshal every object
		objectBuf = append(objectBuf, object.Marshal()...)
	}
	b = append(b, conv.ConvertSizeToBytes(len(objectBuf))...) // size
	b = append(b, objectBuf...)

	return b
}

func (o *ReverseRequest) Unmarshal(b *conv.BinaryIterator) error {
	var (
		err  error
		size int
		buf  *conv.BinaryIterator
	)

	// MagicString
	size, err = b.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read MagicString size")
	}
	buf, err = b.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicString")
	}
	o.MagicString = conv.ConvertBytesToString(buf)
	// MagicInt64
	buf, err = b.Slice(8)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicInt64")
	}
	o.MagicInt64 = conv.ConvertBytesToInt64(buf)
	// MagicBool
	buf, err = b.Slice(1)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicBool")
	}
	o.MagicBool = conv.ConvertBytesToBool(buf)
	// MagicBytes
	size, err = b.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read MagicBytes size")
	}
	buf, err = b.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicBytes")
	}
	o.MagicBytes = buf.Bytes
	// MagicObjectList
	size, err = b.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read MagicObjectList size")
	}
	buf, err = b.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read MagicObjectList")
	}
	var object CommonObject
	for buf.HasNext() {
		if err = object.Unmarshal(buf); err != nil {
			return errors.Wrap(err, "failed to unmarshal MagicObjectList")
		}
		o.MagicObjectList = append(o.MagicObjectList, object)
	}

	return nil
}
