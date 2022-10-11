package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConvGetObjectArray(t *testing.T) {
	_, params, err := GetParams([]byte(`[{a: "bjkjkjk"}, {b: 37878}]`))
	assert.NoError(t, err)
	toCompare := []map[string][]byte{
		{"a": []byte("bjkjkjk")},
		{"b": []byte("37878")},
	}
	assert.Equal(t, len(params), len(toCompare))
	for i, _ := range params {
		assert.Equal(t, params[i], toCompare[i], 0)
	}
}

func Test_ConvGetArray(t *testing.T) {
	_, params, err := GetArray([]byte(`[1, 2, 3, 4]`))
	assert.NoError(t, err)
	toCompare := [][]byte{
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
	}
	assert.Equal(t, len(params), len(toCompare))
	for i, _ := range params {
		assert.Equal(t, params[i], toCompare[i], 0)
	}
}
