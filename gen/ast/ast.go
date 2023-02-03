package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oringik/crypto-chateau/gen/hash"
	lexem2 "github.com/oringik/crypto-chateau/gen/lexem"
)

type Type int

const (
	Uint64 Type = iota
	Uint32
	Uint16
	Uint8
	Int64
	Int32
	Int16
	Int8
	Int
	Byte
	Bool
	String
	Object
)

var LexerTypeToAstType = map[string]Type{
	"uint64": Uint64,
	"uint32": Uint32,
	"uint16": Uint16,
	"uint8":  Uint8,
	"int64":  Int64,
	"int32":  Int32,
	"int16":  Int16,
	"int8":   Int8,
	"int":    Int,
	"byte":   Byte,
	"bool":   Bool,
	"string": String,
	"object": Object,
}

var AstTypeToGoType = map[Type]string{
	Uint64: "uint64",
	Uint32: "uint32",
	Uint16: "uint16",
	Uint8:  "uint8",
	Int64:  "int64",
	Int32:  "int32",
	Int:    "int",
	Int16:  "int16",
	Int8:   "int8",
	Byte:   "byte",
	Bool:   "bool",
	String: "string",
	Object: "object",
}

var AstTypeToDartType = map[Type]string{
	Uint64: "int",
	Uint32: "int",
	Uint16: "int",
	Uint8:  "int",
	Int64:  "int",
	Int32:  "int",
	Int16:  "int",
	Int8:   "int",
	Int:    "int",
	Byte:   "int",
	Bool:   "bool",
	String: "String",
}

type MethodType string

var (
	Handler MethodType = "Handler"
	Stream  MethodType = "Stream"
)

type Chateau struct {
	PackageName                  string
	Services                     []*Service
	ObjectDefinitions            []*ObjectDefinition
	ObjectDefinitionByObjectName map[string]*ObjectDefinition
}

type ObjectDefinition struct {
	Name   string
	Fields []*Field
}

type Service struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name       string
	Hash       hash.HandlerHash
	Params     []*Param
	Returns    []*Return
	MethodType MethodType
	Tags       []Tag
}

type Tag struct {
	Name  string
	Value string
}

type TypeLink struct {
	Type Type

	ObjectName string
	ObjectLink *ObjectDefinition

	IsArray bool
	ArrSize int
}

type Field struct {
	Name string
	Type TypeLink
}

type Param struct {
	Name string
	Type TypeLink
}

type Return struct {
	Name string
	Type TypeLink
}

type Ast struct {
	Chateau *Chateau
}

var currentLexemIndex int
var lexem *lexem2.Lexem
var lexems []*lexem2.Lexem

func getNextLexem() {
	currentLexemIndex++
	if currentLexemIndex >= len(lexems) {
		lexem = nil
		return
	}

	lexem = lexems[currentLexemIndex]
}

func GenerateAst(lxs []*lexem2.Lexem) *Ast {
	if len(lxs) == 0 {
		panic("0 lexems got")
	}
	lexems = lxs
	lexem = lexems[currentLexemIndex]

	ast := &Ast{}

	chateau := astChateau()

	ast.Chateau = chateau

	return ast
}

func astChateau() *Chateau {
	chateau := &Chateau{}
	chateau.ObjectDefinitionByObjectName = make(map[string]*ObjectDefinition)
	if lexem.Type != lexem2.PackageL {
		panic("expected package name")
	}

	getNextLexem()
	chateau.PackageName = lexem.Value

	getNextLexem()

	for lexem != nil {
		if lexem.Type == lexem2.ServiceL {
			getNextLexem()
			chateau.Services = append(chateau.Services, astService())
			getNextLexem()
		} else if lexem.Type == lexem2.ObjectL {
			getNextLexem()
			chateau.ObjectDefinitions = append(chateau.ObjectDefinitions, astObject())
			getNextLexem()
		} else {
			panic("expected service or object " + lexem.Value)
		}
	}

	for _, object := range chateau.ObjectDefinitions {
		chateau.ObjectDefinitionByObjectName[object.Name] = object
	}

	linkObjectsLinksToFields(chateau.ObjectDefinitionByObjectName, chateau.ObjectDefinitions)

	return chateau
}

func linkObjectsLinksToFields(objectDefinitions map[string]*ObjectDefinition, objects []*ObjectDefinition) {
	for _, object := range objects {
		linkObjectLinksToFields(objectDefinitions, object)
	}
}

func linkObjectLinksToFields(objectDefinitions map[string]*ObjectDefinition, object *ObjectDefinition) {
	for _, field := range object.Fields {
		if field.Type.Type == Object {
			field.Type.ObjectLink = objectDefinitions[field.Type.ObjectName]
			linkObjectLinksToFields(objectDefinitions, field.Type.ObjectLink)
		}
	}
}

func astObject() *ObjectDefinition {
	if lexem.Type != lexem2.IdentifierL {
		panic("expected identifier")
	}

	objectDefinition := &ObjectDefinition{}
	objectDefinition.Name = lexem.Value
	getNextLexem()
	if lexem.Type != lexem2.OpenBraceL {
		panic("expected open brace")
	}
	getNextLexem()
	fields := astFields()

	objectDefinition.Fields = fields

	return objectDefinition
}

func astFields() []*Field {
	if lexem.Type == lexem2.CloseBraceL {
		return nil
	}

	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentifierL {
		panic("expected type or identifier " + lexem.Value)
	}

	var fields []*Field

	for lexem.Type != lexem2.CloseBraceL {
		field := astField()

		fields = append(fields, field)
		getNextLexem()
	}

	return fields
}

func astField() *Field {
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentifierL {
		panic("expected type")
	}

	field := &Field{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentifierL {
		typeLink.Type = Object
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			typeLink.ObjectName = lexem.Value[2+getCountDigits(arrSize):]
		} else {
			typeLink.ObjectName = lexem.Value
		}

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else if lexem.Type == lexem2.TypeL {
		var astType Type
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value]
			if !ok {
				panic("unexpected type " + lexem.Value)
			}

			astType = astTypeLocal
		}

		typeLink.Type = astType

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else {
		panic("unexpected lexem type")
	}

	field.Type = typeLink
	getNextLexem()
	if lexem.Type != lexem2.IdentifierL {
		panic("expected identifier")
	}
	field.Name = lexem.Value

	return field
}

func astService() *Service {
	if lexem.Type != lexem2.IdentifierL {
		panic("expected identifier")
	}

	service := &Service{}

	service.Name = lexem.Value
	getNextLexem()
	if lexem.Type != lexem2.OpenBraceL {
		panic("expected open brace")
	}
	getNextLexem()
	var methods []*Method
	for lexem.Type != lexem2.CloseBraceL {
		method := astMethod()
		method.Hash = hash.GetHandlerHash(service.Name, method.Name)

		methods = append(methods, method)
	}

	service.Methods = methods

	return service
}

func astMethod() *Method {
	if lexem.Type != lexem2.MethodL {
		panic("expected method type")
	}

	method := &Method{}
	method.MethodType = MethodType(lexem.Value)

	getNextLexem()
	if lexem.Type != lexem2.IdentifierL {
		panic("expected method name")
	}
	method.Name = lexem.Value

	getNextLexem()
	params := astParamExpr()
	if lexem.Type != lexem2.CloseParenL {
		panic("expected close paren")
	}
	getNextLexem()

	if lexem.Type != lexem2.ReturnArrowL {
		panic("expected return arrow")
	}

	getNextLexem()

	returns := astReturnExpr()

	method.Params = params
	method.Returns = returns

	tags := astTags()

	method.Tags = tags

	return method
}

func astTags() []Tag {
	if lexem.Type != lexem2.OpenBraceL {
		return nil
	}

	getNextLexem()

	var tags []Tag

	for lexem.Type != lexem2.CloseBraceL {
		if lexem.Type != lexem2.IdentifierL {
			panic("expected tag name (identifier)")
		}

		name := lexem.Value

		getNextLexem()

		if lexem.Type != lexem2.ColonL {
			panic("expected tag key val delimited (colon) " + lexem.Value)
		}

		getNextLexem()

		if lexem.Type != lexem2.IdentifierL {
			panic("expected tag value (identifier)")
		}

		value := lexem.Value

		tags = append(tags, Tag{
			Name:  name,
			Value: value,
		})

		getNextLexem()
		if lexem.Type == lexem2.CommaL {
			getNextLexem()
		}
	}

	getNextLexem()

	return tags
}

func astParamExpr() []*Param {
	if lexem.Type != lexem2.OpenParenL {
		panic("expected open paren")
	}

	getNextLexem()

	var params []*Param

	for lexem.Type != lexem2.CloseParenL {
		param := astParam()

		params = append(params, param)
		getNextLexem()
	}

	return params
}

func astReturnExpr() []*Return {
	if lexem.Type != lexem2.OpenParenL {
		panic("expected open paren")
	}

	getNextLexem()

	var returns []*Return

	for lexem.Type != lexem2.CloseParenL {
		ret := astReturn()

		returns = append(returns, ret)
		getNextLexem()
	}

	getNextLexem()

	return returns
}

func astParam() *Param {
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentifierL {
		panic("expected type")
	}

	param := &Param{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentifierL {
		typeLink.Type = Object
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			typeLink.ObjectName = lexem.Value[2+getCountDigits(arrSize):]
		} else {
			typeLink.ObjectName = lexem.Value
		}

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else if lexem.Type == lexem2.TypeL {
		var astType Type
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		}

		typeLink.Type = astType

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else {
		panic("unexpected lexem type")
	}

	param.Type = typeLink
	getNextLexem()
	if lexem.Type != lexem2.IdentifierL {
		panic("expected identifier")
	}
	param.Name = lexem.Value

	return param
}

func astReturn() *Return {
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentifierL {
		panic("expected type or identefier")
	}

	ret := &Return{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentifierL {
		typeLink.Type = Object
		typeLink.ObjectName = lexem.Value

		isArr, arrSize := getArrExistAndSize(lexem.Value)

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else if lexem.Type == lexem2.TypeL {
		var astType Type
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := LexerTypeToAstType[lexem.Value]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		}

		typeLink.Type = astType

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else {
		panic("unexpected lexem type")
	}

	ret.Type = typeLink

	return ret
}

func getArrExistAndSize(s string) (bool, int) {
	var isArr bool
	var arrSize int
	openArrBracketIndex := strings.Index(s, "[")
	if openArrBracketIndex != -1 {
		isArr = true
		closeArrBracketIndex := strings.Index(s, "]")
		if closeArrBracketIndex-openArrBracketIndex > 1 {
			size, err := strconv.Atoi(s[openArrBracketIndex+1 : closeArrBracketIndex])
			if err != nil {
				panic("unexpected arr size")
			}
			arrSize = size
		}
	}

	return isArr, arrSize
}

func getCountDigits(num int) int {
	var count int
	for num != 0 {
		count++
		num /= 10
	}

	return count
}

func FillDefaultObjectValues(objectDefinitions map[string]*ObjectDefinition, objectName string) string {
	object := objectDefinitions[objectName]
	var fieldsFilled string
	for i, objField := range object.Fields {
		fieldsFilled += FillTypeWithDefaultValue(objField, false, true)
		if i != len(object.Fields)-1 {
			fieldsFilled += ","
		}
	}
	return fieldsFilled
}

func FillTypeWithDefaultValue(field *Field, isArrChecked bool, isWithFieldName bool) string {
	res := strings.ToUpper(field.Name[0:1]) + field.Name[1:] + ": "
	if !isWithFieldName {
		res = ""
	}

	if field.Type.IsArray && !isArrChecked {
		tp := FillTypeWithDefaultValue(field, true, false)
		isArrChecked = false
		isWithFieldName = true

		res += fmt.Sprintf("List.filled(0, %s, growable: true)", tp)
		return res
	}

	res += FillWithDefaultValueType(&field.Type)

	return res
}

func FillWithDefaultValueType(tp *TypeLink) string {
	switch tp.Type {
	case String:
		return `""`
	case Int8, Int16, Int32, Int64, Uint8, Uint16, Uint32, Uint64, Int:
		return "0"
	case Bool:
		return "true"
	case Byte:
		return "0xff"
	case Object:
		var fieldsFilled string
		for i, objField := range tp.ObjectLink.Fields {
			fieldsFilled += FillTypeWithDefaultValue(objField, false, true)
			if i != len(tp.ObjectLink.Fields)-1 {
				fieldsFilled += ","
			}
		}
		return tp.ObjectName + "(" + fieldsFilled + ")"
	}

	return ""
}
