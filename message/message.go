package message

import "github.com/oringik/crypto-chateau/gen/conv"

type Message interface {
	Marshal() []byte
	Unmarshal(iterator *conv.BinaryIterator) error
}
