package conv

import (
	"errors"

	"github.com/oringik/crypto-chateau/gen/hash"
)

func GetHandler(p []byte) (protocol []byte, handlerKey hash.HandlerHash, payloadOffset int, err error) {
	if len(p) < 6 {
		return nil, hash.HandlerHash{}, 0, errors.New("invalid payload: too short")
	}

	protocol = p[:1]
	handlerBytes := p[1:5]
	handlerKey = hash.HandlerHash{handlerBytes[0], handlerBytes[1], handlerBytes[2], handlerBytes[3]}

	return protocol, handlerKey, 5, nil
}
