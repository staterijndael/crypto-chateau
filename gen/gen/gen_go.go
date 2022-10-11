package gen

import (
	"fmt"
	ast2 "github.com/Oringik/crypto-chateau/gen/ast"
	"github.com/Oringik/crypto-chateau/gen/conv"
	"strconv"
	"strings"
)

var result string
var ast *ast2.Ast

func GenerateDefinitions(astLocal *ast2.Ast) string {
	ast = astLocal

	fillPackage()
	fillImports()
	fillMethodTypes()
	fillServices()
	fillSqueezes()
	fillObjects()
	fillPeers()

	return result
}

func fillImports() {
	result += "import \"errors\"\n"
	result += "import \"context\"\n"
	result += "import \"github.com/Oringik/crypto-chateau/gen/conv\"\n"
	result += "import \"github.com/Oringik/crypto-chateau/peer\"\n"
	result += "import \"github.com/Oringik/crypto-chateau/message\"\n\n"
}

func fillPackage() {
	result += "package " + ast.Chateau.PackageName + "\n\n"
}

func fillMethodTypes() {
	result += `type HandlerFunc func(ctx context.Context, msg message.Message) (message.Message, error)
type StreamFunc func(ctx context.Context, peer *peer.Peer, msg message.Message) error` + "\n\n"
}

func fillServices() {
	for _, service := range ast.Chateau.Services {
		result += "type " + service.Name + " interface { \n"
		for _, method := range service.Methods {
			result += "\t" + method.Name + "(ctx context.Context, "
			if method.MethodType == ast2.Stream {
				result += "peer *peer.Peer, "
			}
			for i, param := range method.Params {
				result += param.Name + " "
				if param.Type.IsArray {
					result += "["
					if param.Type.ArrSize != 0 {
						result += strconv.Itoa(param.Type.ArrSize)
					}
					result += "]"
				}

				if param.Type.Type == ast2.Object {
					result += "*" + param.Type.ObjectName
				} else {
					result += ast2.AstTypeToGoType[param.Type.Type]
				}
				if i != len(method.Params)-1 {
					result += ", "
				}
			}
			result += ") ("
			for i, ret := range method.Returns {
				if ret.Type.IsArray {
					result += "["
					if ret.Type.ArrSize != 0 {
						result += strconv.Itoa(ret.Type.ArrSize)
					}
					result += "]"
				}

				if ret.Type.Type == ast2.Object {
					result += "*" + ret.Type.ObjectName
				} else {
					result += ast2.AstTypeToGoType[ret.Type.Type]
				}
				if i != len(method.Returns)-1 {
					result += ", "
				}
			}
			result += ", error)\n"
		}
		result += "}\n\n"
	}
}

func fillSqueezes() {
	for _, service := range ast.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				result += fmt.Sprintf(`func %sSqueeze(fnc func(context.Context, *%s) (*%s, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*%s); ok {
			return fnc(ctx, msg.(*%s))
		} else {
			return nil, errors.New("unknown message type: expected %s")
		}
	}
}`+"\n\n", method.Name, method.Params[0].Type.ObjectName, method.Returns[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName)
			} else {
				result += fmt.Sprintf(`func %sSqueeze(fnc func(context.Context, *peer.Peer, *%s) error) StreamFunc {
	return func(ctx context.Context, peer *peer.Peer, msg message.Message) error {
		if _, ok := msg.(*%s); ok {
			return fnc(ctx, peer, msg.(*%s))
		} else {
			return errors.New("unknown message type: expected %s")
		}
	}
}`+"\n\n", method.Name, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName)
			}
		}
	}
}

func fillObjects() {
	for _, object := range ast.Chateau.ObjectDefinitions {
		result += "type " + object.Name + " struct {\n"
		for _, field := range object.Fields {
			field.Name = strings.Title(field.Name)
			field.Type.ObjectName = strings.Title(field.Type.ObjectName)

			result += "\t" + field.Name + " "

			if field.Type.IsArray {
				result += "["
				if field.Type.ArrSize != 0 {
					result += strconv.Itoa(field.Type.ArrSize)
				}
				result += "]"
			}

			if field.Type.Type == ast2.Object {
				result += "*" + field.Type.ObjectName + "\n"
			} else {
				result += ast2.AstTypeToGoType[field.Type.Type] + "\n"
			}
		}
		result += "}\n\n"

		// marshal

		result += "func (o *" + object.Name + ") Marshal() []byte {\n"
		result += "\tvar buf []byte\n"
		result += "\t" + `buf = append(buf, '{')` + "\n"
		for i, field := range object.Fields {
			convFunction := conv.ConvFunctionMarhsalByType(field.Type.Type)
			result += fmt.Sprintf("\tvar result%s []byte\n", field.Name)
			result += fmt.Sprintf("\t"+`result%s = append(result%s, []byte("%s:")...)`, field.Name, field.Name, field.Name) + "\n"
			if field.Type.IsArray {
				result += fmt.Sprintf("\t"+`result%s = append(result%s, '[')`, field.Name, field.Name) + "\n"
				result += "\tfor i, val := range o." + field.Name + " {\n"
				result += fmt.Sprintf("\t\tresult%s = append(result%s, conv.%s(val)...)\n", field.Name, field.Name, convFunction)
				result += "\t\tif i != len(o." + field.Name + ") - 1 {\n"
				result += "\t\t\t" + fmt.Sprintf(`result%s = append(result%s, ',')`, field.Name, field.Name) + "\n"
				result += "\t\t}\n"
				result += "\t}\n"
				result += fmt.Sprintf("\t"+`result%s = append(result%s, ']')`, field.Name, field.Name) + "\n\n"
			} else {
				result += fmt.Sprintf("\tresult%s = append(result%s, conv.%s(o.%s)...)\n", field.Name, field.Name, convFunction, field.Name)
			}
			result += fmt.Sprintf("\tbuf = append(buf, result%s...)\n", field.Name)
			if i != len(object.Fields)-1 {
				result += "\tbuf = append(buf, ',')\n"
			}
		}
		result += "\t" + `buf = append(buf, '}')` + "\n"
		result += "\treturn buf\n }\n\n"

		// unmarshal

		result += "func (o *" + object.Name + ") Unmarshal(params map[string][]byte) error {\n"
		for _, field := range object.Fields {
			convFunction := conv.ConvFunctionUnmarshalByType(field.Type.Type)
			if field.Type.Type == ast2.Object {
				if field.Type.IsArray {
					result += fmt.Sprintf("\t"+`_, arr, err := conv.GetArray(params["%s"])`+"\n", field.Name)
					result += "\tif err != nil {\n\t\treturn err\n\t}\n"
					result += "\tfor _, objBytes := range arr {\n"
					result += "\t\tvar curObj " + field.Type.ObjectName + "\n"
					result += fmt.Sprintf("\t\t"+`conv.%s(&curObj,objBytes)`+"\n", convFunction)
					result += fmt.Sprintf("\t\to.%s = append(o.%s, curObj)\n", field.Name, field.Name)
					result += "\t}\n"
				} else {
					result += fmt.Sprintf("\to.%s = &%s{}\n", field.Name, field.Type.ObjectName)
					result += fmt.Sprintf("\t"+`conv.%s(o.%s,params["%s"])`+"\n", convFunction, field.Name, field.Name)
				}
			} else {
				if field.Type.IsArray {
					result += fmt.Sprintf("\t"+`_, arr, err := conv.GetArray(params["%s"])`+"\n", field.Name)
					result += "\tif err != nil {\n\t\treturn err\n\t}\n"
					var iOrMiss string
					if field.Type.ArrSize != 0 {
						iOrMiss = "i"
					} else {
						iOrMiss = "-"
					}
					result += fmt.Sprintf("\tfor %s, valByte := range arr {\n", iOrMiss)
					if field.Type.ArrSize != 0 {
						result += fmt.Sprintf("\t\to.%s[i] = conv.%s(valByte)\n", field.Name, convFunction)
					} else {
						result += "\t\tvar curVal " + ast2.AstTypeToGoType[field.Type.Type] + "\n"
						result += fmt.Sprintf("\t\t"+`curVal = conv.%s(valByte)`+"\n", convFunction)
						result += fmt.Sprintf("\t\to.%s = append(o.%s, curVal)\n", field.Name, field.Name)
					}
					result += "\t}\n"
				} else {
					result += fmt.Sprintf("\t"+`o.%s = conv.%s(params["%s"])`+"\n", field.Name, convFunction, field.Name)
				}
			}
		}
		result += "\treturn nil\n}\n\n"
	}
}

func fillPeers() {
	for _, service := range ast.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Stream {
				result += "type Peer" + method.Name + " struct {\n"
				result += "\thandlerName string\n"
				result += "\tpeer *peer.Peer\n"
				result += "}\n"

				result += "func (p *Peer" + method.Name + ") WriteResponse(ctx context.Context, msg *" + method.Returns[0].Type.ObjectName + ") error {\n"
				result += "\treturn p.peer.WriteResponse(p.handlerName, msg)\n"
				result += "}\n"

				result += "func (p *Peer" + method.Name + ") WriteError(ctx context.Context, err error) error"
				result += "\treturn p.peer.WriteError(p.handlerName, err)\n"
				result += "}\n"
			}
		}
	}
}
