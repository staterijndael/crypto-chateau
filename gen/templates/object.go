package templates

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/oringik/crypto-chateau/gen/ast"
)

type ObjectTemplate struct {
	tpl *template.Template
}

func NewObjectTemplate() (*ObjectTemplate, error) {
	tpl := template.New("object")
	tpl = tpl.Funcs(objectTemplateFunc)
	tpl, err := tpl.ParseFS(embFS, "object.go.tpl")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse object template")
	}

	return &ObjectTemplate{
		tpl: tpl,
	}, nil
}

func (t *ObjectTemplate) Gen(definition *ast.ObjectDefinition) (string, error) {
	b := bytes.NewBuffer(nil)

	err := t.tpl.ExecuteTemplate(b, "object.go.tpl", definition)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute object template")
	}

	return b.String(), nil
}

var objectTemplateFunc = template.FuncMap{
	"mul": func(a, b int) int { return a * b },
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	},
	"eqType": func(a ast.Type, b string) bool {
		return strings.EqualFold(ast.AstTypeToGoType[a], b)
	},
	"GoType":  GoType,
	"ToCamel": strcase.ToCamel,
}

func GoType(t *ast.TypeLink) (string, error) {
	var textType string

	switch t.Type {
	case ast.Uint64:
		textType = "uint64"
	case ast.Uint32:
		textType = "uint32"
	case ast.Uint16:
		textType = "uint16"
	case ast.Uint8:
		textType = "uint8"
	case ast.Int64:
		textType = "int64"
	case ast.Int32:
		textType = "int32"
	case ast.Int16:
		textType = "int16"
	case ast.Int8:
		textType = "int8"
	case ast.Byte:
		textType = "byte"
	case ast.Bool:
		textType = "bool"
	case ast.String:
		textType = "string"
	case ast.Object:
		textType = strcase.ToCamel(t.ObjectName)
	default:
		return "", errors.New("unknown type: " + strconv.Itoa(int(t.Type)))
	}

	if !t.IsArray {
		return textType, nil
	}

	if t.ArrSize == 0 {
		return "[]" + textType, nil
	}

	return "[" + strconv.Itoa(t.ArrSize) + "]" + textType, nil
}
