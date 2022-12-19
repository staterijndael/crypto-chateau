package conv

import (
	"errors"
)

func GetHandler(p []byte) (protocol []byte, handlerKey []byte, payloadOffset int, err error) {
	if len(p) < 6 {
		return nil, nil, 0, errors.New("invalid payload: too short")
	}

	protocol = p[:1]
	handlerKey = p[1:5]

	return protocol, handlerKey, 5, nil
}
