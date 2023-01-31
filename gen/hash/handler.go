package hash

import (
	"crypto/sha256"
	"fmt"
)

type HandlerHash [4]byte

func (h HandlerHash) Code() string {
	return fmt.Sprintf("0x%X, 0x%X, 0x%X, 0x%X", h[0], h[1], h[2], h[3])
}

// GetHandlerHash returns first 4 bytes of sha256 hash of serviceName/handlerName
func GetHandlerHash(serviceName string, handlerName string) [4]byte {
	hash := sha256.Sum256([]byte(serviceName + "/" + handlerName))

	var result [4]byte
	copy(result[:], hash[len(hash)-4:])

	return result
}
