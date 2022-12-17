package templates

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oringik/crypto-chateau/gen/ast"
)

const exceptedObjectCode = ``

func TestObjectTemplate_Gen(t *testing.T) {
	ot, err := NewObjectTemplate()

	require.NoError(t, err, "failed to create object template")

	def := &ast.ObjectDefinition{
		Name: "MagicRequest",
		Fields: []*ast.Field{
			{
				Name: "MagicString",
				Type: ast.TypeLink{
					Type: ast.String,
				},
			},
			{
				Name: "MagicUint32",
				Type: ast.TypeLink{
					Type: ast.Uint32,
				},
			},
			{
				Name: "MagicBytes",
				Type: ast.TypeLink{
					Type:    ast.Byte,
					IsArray: true,
					ArrSize: 16,
				},
			},
			{
				Name: "MagicObjects",
				Type: ast.TypeLink{
					Type:       ast.Object,
					ObjectName: "WonderObject",
					IsArray:    true,
				},
			},
		},
	}

	code, err := ot.Gen(def)

	require.NoError(t, err, "failed to generate code")

	err = os.WriteFile("object_example.go", []byte(code), 0644)
	require.NoError(t, err, "failed to save generated code")

	require.Equal(t, exceptedObjectCode, code, "generated code is not as expected")
}
