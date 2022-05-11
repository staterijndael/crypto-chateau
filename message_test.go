package crypto_chateau

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestMessage struct {
	SomeInt    int64
	SomeBool   bool
	SomeString string
}

func (t *TestMessage) Marshal() {
	// do something useful
}

func Test_ParseMessage(t *testing.T) {
	num := []byte{123, 0, 0, 0, 0, 0, 0, 0}
	msg := []byte("handlerName# someBool: 1, someInt : ")
	msg = append(msg, num...)
	msg = append(msg, []byte(", someString :asdmamsd")...)
	handlerName, msgStruct, err := ParseMessage(msg, &TestMessage{})
	assert.NoError(t, err)
	assert.Equal(t, "handlerName", string(handlerName))

	t.Log(msgStruct)
}
