package dh

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDh(t *testing.T) {
	clientAliceStore := &KeyStore{}
	clientBobStore := &KeyStore{}

	clientAliceStore.GeneratePrivateKey()
	clientBobStore.GeneratePrivateKey()

	assert.NotEqual(t, clientAliceStore.PrivateKey, clientBobStore.PrivateKey)
	log.Println(clientAliceStore.PrivateKey, clientBobStore.PrivateKey)

	err := clientAliceStore.GeneratePublicKey()
	assert.NoError(t, err)
	err = clientBobStore.GeneratePublicKey()
	assert.NoError(t, err)

	assert.NotEqual(t, t, clientAliceStore.PublicKey, clientBobStore.PublicKey)
	log.Println(clientAliceStore.PublicKey, clientBobStore.PublicKey)

	err = clientAliceStore.GenerateSharedKey(clientBobStore.PublicKey)
	err = clientBobStore.GenerateSharedKey(clientAliceStore.PublicKey)

	assert.Equal(t, clientAliceStore.SharedKey, clientBobStore.SharedKey)
	log.Println(clientAliceStore.SharedKey, clientBobStore.SharedKey)
}
