package gen

import (
	"fmt"
	ast2 "github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/gen/conv"
	"strconv"
	"strings"
	"unicode"
)

var result string
var ast *ast2.Ast

func GenerateDefinitions(astLocal *ast2.Ast) string {
	ast = astLocal

	fillPackage()
	fillImports()
	fillServices()
	fillSqueezes()
	fillObjects()
	fillPeers()
	fillPeersSqueezes()
	fillInitHandlers()
	fillNewServer()

	fillClients()

	return result
}

func fillImports() {
	result += "import \"errors\"\n"
	result += "import \"context\"\n"
	result += "import \"github.com/oringik/crypto-chateau/gen/conv\"\n"
	result += "import \"github.com/oringik/crypto-chateau/peer\"\n"
	result += "import \"github.com/oringik/crypto-chateau/message\"\n"
	result += "import \"github.com/oringik/crypto-chateau/server\"\n"
	result += "import \"go.uber.org/zap\"\n"
	result += "import \"github.com/oringik/crypto-chateau/transport\"\n"
	result += "import \"net\"\n\n"

}

func fillPackage() {
	result += "package " + ast.Chateau.PackageName + "\n\n"
}

func fillClients() {
	for _, service := range ast.Chateau.Services {
		result += "type Client" + service.Name + " struct {\n"
		result += "\tpeer *peer.Peer\n"
		result += "}\n\n"

		result += "func NewClient" + service.Name + "(host string, port int) (*Client" + service.Name + ", error) {\n"
		result += "\tconn, err := net.Dial(\"tcp\", host + \":\" + strconv.Itoa(port))\n"
		result += "\tif err != nil{\n"
		result += "\t\treturn nil, err\n"
		result += "\t}\n"
		result += "\tconn, err = transport.ServerHandshake(conn)\n"
		result += "\tif err != nil{\n"
		result += "\t\treturn nil, err\n"
		result += "\t}\n"
		result += "securedPeer := peer.NewPeer(conn)\n"
		result += "client := &Client" + service.Name + "{peer: securedPeer}\n"
		result += "return client, nil\n"
		result += "}\n\n"
		for _, method := range service.Methods {
			result += "\t" + "func (c *Client" + service.Name + ") " + method.Name + "(ctx context.Context, "
			if method.MethodType == ast2.Stream {
				result += fmt.Sprintf("peer *Peer%s, ", method.Name)
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
			result += ") "
			if method.MethodType != ast2.Stream {
				result += "("
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
				result += ", "
			}
			result += "error"
			if method.MethodType != ast2.Stream {
				result += ")"
			}
			result += "{\n"
			result += "\terr := c.peer.WriteResponse(\"" + method.Name + "\"," + method.Params[0].Name + ")\n\n"
			result += fmt.Sprintf(`msg := make([]byte, 0, 1024)

	for {
		buf := make([]byte, 1024)
		n, err := c.peer.Read(buf)
		if err != nil {
			return nil, err
		}

		if n == 0 {
			break
		}

		if n < len(buf) {
			buf = buf[:n]
			msg = append(msg, buf...)
			break
		}

		msg = append(msg, buf...)
	}

	_, n, err := conv.GetHandlerName(msg)
	if err != nil {
		return nil, err
	}

	if n >= len(msg) {
		return nil, errors.New("incorrect message")
	}

	_, responseMsgParams, err := conv.GetParams(msg[n:])
	if err != nil {
		return nil, err
	}

	respMsg := &%s{}

	err = respMsg.Unmarshal(responseMsgParams)
	if err != nil {
		return nil, err
	}
	
	return respMsg, nil
`, method.Returns[0].Type.ObjectName)

			result += "}\n\n"
		}
	}
}

func fillServices() {
	for _, service := range ast.Chateau.Services {
		result += "type " + service.Name + " interface { \n"
		for _, method := range service.Methods {
			result += "\t" + method.Name + "(ctx context.Context, "
			if method.MethodType == ast2.Stream {
				result += fmt.Sprintf("peer *Peer%s, ", method.Name)
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
			result += ") "
			if method.MethodType != ast2.Stream {
				result += "("
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
				result += ", "
			}
			result += "error"
			if method.MethodType != ast2.Stream {
				result += ")"
			}
			result += "\n"
		}
		result += "}\n\n"
	}
}

func fillSqueezes() {
	for _, service := range ast.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				result += fmt.Sprintf(`func %sSqueeze(fnc func(context.Context, *%s) (*%s, error)) server.HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*%s); ok {
			return fnc(ctx, msg.(*%s))
		} else {
			return nil, errors.New("unknown message type: expected %s")
		}
	}
}`+"\n\n", method.Name, method.Params[0].Type.ObjectName, method.Returns[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName)
			} else {
				result += fmt.Sprintf(`func %sSqueeze(fnc func(context.Context, *Peer%s, *%s) error) server.StreamFunc {
	return func(ctx context.Context, peer interface{}, msg message.Message) error {
		if _, ok := msg.(*%s); ok {
			return fnc(ctx, peer.(*Peer%s), msg.(*%s))
		} else {
			return errors.New("unknown message type: expected %s")
		}
	}
}`+"\n\n", method.Name, method.Name, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName, method.Name, method.Params[0].Type.ObjectName, method.Params[0].Type.ObjectName)
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

			fmt.Println(fmt.Sprintf("fieldName - %s", field.Name))
			fmt.Println()
			fmt.Println(fmt.Sprintf("bytes - %v", []byte(field.Name)))
			fmt.Println()
			fmt.Println()
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

		alreadyWasArraySymb := ":"

		result += "func (o *" + object.Name + ") Unmarshal(params map[string][]byte) error {\n"
		for _, field := range object.Fields {
			convFunction := conv.ConvFunctionUnmarshalByType(field.Type.Type)
			if field.Type.Type == ast2.Object {
				if field.Type.IsArray {
					result += fmt.Sprintf("\t"+`_, arr, err %s= conv.GetArray(params["%s"])`+"\n", alreadyWasArraySymb, field.Name)
					result += "\tif err != nil {\n\t\treturn err\n\t}\n"
					result += "\tfor _, objBytes := range arr {\n"
					result += "\t\tvar curObj *" + field.Type.ObjectName + "\n"
					result += fmt.Sprintf("\t\t"+`conv.%s(curObj,objBytes)`+"\n", convFunction)
					result += fmt.Sprintf("\t\to.%s = append(o.%s, curObj)\n", field.Name, field.Name)
					result += "\t}\n"

					alreadyWasArraySymb = ""
				} else {
					result += fmt.Sprintf("\to.%s = &%s{}\n", field.Name, field.Type.ObjectName)
					result += fmt.Sprintf("\t"+`conv.%s(o.%s,params["%s"])`+"\n", convFunction, field.Name, field.Name)
				}
			} else {
				if field.Type.IsArray {
					result += fmt.Sprintf("\t"+`_, arr, err %s= conv.GetArray(params["%s"])`+"\n", alreadyWasArraySymb, field.Name)
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

					alreadyWasArraySymb = ""
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
				result += "\tpeer *peer.Peer\n"
				result += "}\n\n"

				result += "func (p *Peer" + method.Name + ") WriteResponse(msg *" + method.Returns[0].Type.ObjectName + ") error {\n"
				result += fmt.Sprintf("\t"+`return p.peer.WriteResponse("%s", msg)`+"\n", method.Name)
				result += "}\n\n"

				result += "func (p *Peer" + method.Name + ") WriteError(err error) error {\n"
				result += fmt.Sprintf("\t"+`return p.peer.WriteError("%s", err)`+"\n", method.Name)
				result += "}\n\n"
			}
		}
	}
}

func fillPeersSqueezes() {
	result += "func getPeerByHandlerName(handlerName string, peer *peer.Peer) interface{}{\n"
	for _, service := range ast.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Stream {
				result += "\tif handlerName == \"" + method.Name + "\" {\n"
				result += "\t\treturn Peer" + method.Name + "{peer}\n"
				result += "\t}\n\n"
			}
		}
	}
	result += "\treturn nil\n"
	result += "}\n\n"
}

func fillInitHandlers() {
	var endpointArgs string
	for i, service := range ast.Chateau.Services {
		endpointArgs += string(unicode.ToLower(rune(service.Name[0])))
		if len(service.Name) > 1 {
			endpointArgs += service.Name[1:]
		}
		endpointArgs += " " + service.Name

		if i != len(ast.Chateau.Services)-1 {
			endpointArgs += ","
		}
	}
	result += fmt.Sprintf("func initHandlers(%s) map[string]*server.Handler {\n", endpointArgs)
	result += "\thandlers := make(map[string]*server.Handler)\n\n"
	for _, service := range ast.Chateau.Services {
		for _, method := range service.Methods {
			var methodType string
			if method.MethodType == ast2.Handler {
				methodType = "server.HandlerT"
			} else if method.MethodType == ast2.Stream {
				methodType = "server.StreamT"
			}
			var serviceNameLower string
			serviceNameLower += string(unicode.ToLower(rune(service.Name[0])))
			if len(service.Name) > 1 {
				serviceNameLower += service.Name[1:]
			}
			result += fmt.Sprintf("\t"+`handlers["%s"] = &server.Handler{
		CallFunc%s: %sSqueeze(%s.%s),
		HandlerType:     %s,
		RequestMsgType:  &%s{},
	}`+"\n\n", method.Name, string(method.MethodType), method.Name, serviceNameLower, method.Name, methodType, method.Params[0].Type.ObjectName)
		}
	}
	result += "\treturn handlers\n"
	result += "}\n\n"
}

func fillNewServer() {
	var endpointArgs string
	var endpointNames string
	for i, service := range ast.Chateau.Services {
		endpointArgs += string(unicode.ToLower(rune(service.Name[0])))
		endpointNames += string(unicode.ToLower(rune(service.Name[0])))
		if len(service.Name) > 1 {
			endpointArgs += service.Name[1:]
			endpointNames += service.Name[1:]
		}
		endpointArgs += " " + service.Name

		if i != len(ast.Chateau.Services)-1 {
			endpointArgs += ","
		}
	}

	result += fmt.Sprintf("func NewServer(cfg *server.Config, logger *zap.Logger, %s) *server.Server {\n", endpointArgs)
	result += fmt.Sprintf("\thandlers := initHandlers(%s)\n\n", endpointNames)
	result += "\treturn server.NewServer(cfg, logger, handlers, getPeerByHandlerName)\n"
	result += "}\n\n"
}
