package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseMsg(t *testing.T) {
	msg := []byte("aa|bbb|cccc|myadaa")

	fields, err := parseMsg(msg)
	assert.NoError(t, err)

	for _, field := range fields {
		t.Log(string(field))
	}
}

func Test_FormatMsg(t *testing.T) {
	msg := formatMsg([]byte("aaa"), []byte("bbbb"), []byte("ddddd"))

	t.Log(string(msg))
}
