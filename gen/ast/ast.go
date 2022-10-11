package ast

import (
	lexem2 "github.com/oringik/crypto-chateau/gen/lexem"
	"strconv"
	"strings"
)

type Type int

const (
	Uint32 Type = iota
	Uint64
	Uint8
	Byte
	Bool
	String
	Object
)

var lexerTypeToAstType = map[string]Type{
	"byte":   Byte,
	"uint32": Uint32,
	"uint64": Uint64,
	"uint8":  Uint8,
	"string": String,
	"bool":   Bool,
}

var AstTypeToGoType = map[Type]string{
	Uint32: "uint32",
	Uint64: "uint64",
	Uint8:  "uint8",
	Byte:   "byte",
	String: "string",
	Bool:   "bool",
}

type MethodType string

var (
	Handler MethodType = "Handler"
	Stream  MethodType = "Stream"
)

type Chateau struct {
	PackageName       string
	Services          []*Service
	ObjectDefinitions []*ObjectDefinition
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
	Params     []*Param
	Returns    []*Return
	MethodType MethodType
}

type TypeLink struct {
	Type       Type
	ObjectName string

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
			panic("expected service or object")
		}
	}

	return chateau
}

func astObject() *ObjectDefinition {
	if lexem.Type != lexem2.IdentefierL {
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

	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentefierL {
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
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentefierL {
		panic("expected type")
	}

	field := &Field{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentefierL {
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
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value]
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
	if lexem.Type != lexem2.IdentefierL {
		panic("expected identifier")
	}
	field.Name = lexem.Value

	return field
}

func astService() *Service {
	if lexem.Type != lexem2.IdentefierL {
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
		astMethod := astMethod()

		methods = append(methods, astMethod)
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
	if lexem.Type != lexem2.IdentefierL {
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

	return method
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
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentefierL {
		panic("expected type")
	}

	param := &Param{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentefierL {
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
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value]
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
	if lexem.Type != lexem2.IdentefierL {
		panic("expected identifier")
	}
	param.Name = lexem.Value

	return param
}

func astReturn() *Return {
	if lexem.Type != lexem2.TypeL && lexem.Type != lexem2.IdentefierL {
		panic("expected type or identefier")
	}

	ret := &Return{}
	typeLink := TypeLink{}
	if lexem.Type == lexem2.IdentefierL {
		typeLink.Type = Object
		typeLink.ObjectName = lexem.Value

		isArr, arrSize := getArrExistAndSize(lexem.Value)

		typeLink.IsArray = isArr
		typeLink.ArrSize = arrSize
	} else if lexem.Type == lexem2.TypeL {
		var astType Type
		isArr, arrSize := getArrExistAndSize(lexem.Value)
		if isArr {
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value[2+getCountDigits(arrSize):]]
			if !ok {
				panic("unexpected type")
			}

			astType = astTypeLocal
		} else {
			astTypeLocal, ok := lexerTypeToAstType[lexem.Value]
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
