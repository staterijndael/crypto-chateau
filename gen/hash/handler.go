package hash

import (
	"crypto/sha256"
	"fmt"
)

type HandlerHash [4]byte

func (h HandlerHash) Code() string {
	return fmt.Sprintf("hash.HandlerHash{0x%X, 0x%X, 0x%X, 0%X}", h[0], h[1], h[2], h[3])
}

// GetHandlerHash returns first 4 bytes of sha256 hash of serviceName/handlerName
func GetHandlerHash(serviceName string, handlerName string) [4]byte {
	hash := sha256.New().Sum([]byte(serviceName + "/" + handlerName))

	return [4]byte{hash[0], hash[1], hash[2], hash[3]}
}
