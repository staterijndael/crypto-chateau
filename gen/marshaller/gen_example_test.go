package marshaller

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oringik/crypto-chateau/gen/conv"
)

func TestCommonObject(t *testing.T) {
	req := &ReverseRequest{
		MagicString: "Hello, world!",
		MagicInt64:  1,
		MagicBool:   true,
		MagicBytes:  []byte{1, 2, 3, 4, 5},
		MagicObjectList: []CommonObject{
			{
				MagicString: "Hello, world!",
			},
		},
	}

	buf := req.Marshal()

	actual := &ReverseRequest{}
	err := actual.Unmarshal(conv.NewBinaryIterator(buf))

	require.NoError(t, err)
	require.EqualValues(t, req, actual)
}

func BenchmarkCommonObject(b *testing.B) {
	for i := 0; i < b.N; i++ {

		req := &ReverseRequest{
			MagicString: "Hello, world!",
			MagicInt64:  1,
			MagicBool:   true,
			MagicBytes:  []byte{1, 2, 3, 4, 5},
			MagicObjectList: []CommonObject{
				{
					MagicString: "Hello, world!",
				},
			},
		}

		buf := req.Marshal()

		actual := &ReverseRequest{}
		_ = actual.Unmarshal(conv.NewBinaryIterator(buf))
	}
}
