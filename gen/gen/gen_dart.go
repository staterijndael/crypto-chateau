package gen

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	ast2 "github.com/oringik/crypto-chateau/gen/ast"
)

var resultDart string
var astDart *ast2.Ast

func GenerateDefinitionsDart(astLocal *ast2.Ast) string {
	astDart = astLocal

	fillImportsDart()
	fillMethodsDart()
	fillObjectsDart()

	return resultDart
}

func fillImportsDart() {
	resultDart += "import 'package:crypto_chateau_dart/transport/conn_bloc.dart';\n"
	resultDart += "import 'dart:convert';\n"
	resultDart += "import 'dart:async';\n"
	resultDart += "import 'dart:typed_data';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/models.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/conv.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/client.dart';\n\n"
}

func fillObjectsDart() {
	for _, object := range astDart.Chateau.ObjectDefinitions {
		resultDart += "class " + object.Name + " extends Message { \n"
		var constructorArgs string
		for i, field := range object.Fields {
			fieldNameLowed := string(unicode.ToLower(rune(field.Name[0])))
			if len(field.Name) > 1 {
				fieldNameLowed += field.Name[1:]
			}

			field.Name = fieldNameLowed
			field.Type.ObjectName = strings.Title(field.Type.ObjectName)

			if field.Type.IsArray {
				if field.Type.ArrSize != 0 {
					resultDart += "\t// arr max elements count: " + strconv.Itoa(field.Type.ArrSize) + "\n"
				}
				resultDart += "\tList<"
			}

			if !field.Type.IsArray {
				resultDart += "\t"
			}

			if field.Type.Type == ast2.Object {
				resultDart += field.Type.ObjectName
			} else {
				resultDart += ast2.AstTypeToDartType[field.Type.Type]
			}

			if field.Type.IsArray {
				resultDart += ">"
			}

			resultDart += "? " + field.Name + "; \n"
			constructorArgs += "this." + field.Name
			if i != len(object.Fields)-1 {
				constructorArgs += ", "
			}
		}

		resultDart += "\n"

		// constructor

		if constructorArgs != "" {
			resultDart += "\t" + object.Name + "({" + constructorArgs + "});\n\n"
		} else {
			resultDart += "\t" + object.Name + "();\n\n"
		}

		// marshal

		resultDart += "\tUint8List Marshal() {\n"
		resultDart += "\t\tList<int> buf = List.empty(growable: true);\n"
		resultDart += "\t\t" + `buf.addAll('{'.codeUnits);` + "\n"
		for i, field := range object.Fields {
			//convFunction := conv.ConvFunctionMarhsalByType(field.Type.Type)
			convFunction := "TODO_CHANGE"
			resultDart += fmt.Sprintf("\t\tList<int> resultDart%s = List.empty(growable: true);\n", field.Name)
			if field.Type.IsArray {
				resultDart += fmt.Sprintf("\t\t"+`resultDart%s.addAll('['.codeUnits);`, field.Name) + "\n"
				resultDart += "\t\tfor (int i = 0; i < " + field.Name + "!.length; i++) {\n"
				resultDart += "\t\t\tvar val = " + field.Name + "![i];\n"
				resultDart += fmt.Sprintf("\t\tresultDart%s.addAll('%s:'.codeUnits);\n", field.Name, strings.Title(field.Name))
				resultDart += fmt.Sprintf("\t\t\tresultDart%s.addAll(%s(val));\n", field.Name, convFunction)
				resultDart += "\t\t\tif (i != " + field.Name + "!.length - 1) {\n"
				resultDart += "\t\t\t\t" + fmt.Sprintf(`resultDart%s.addAll(','.codeUnits);`, field.Name) + "\n"
				resultDart += "\t\t\t}\n"
				resultDart += "\t\t}\n"
				resultDart += fmt.Sprintf("\t\t"+`resultDart%s.addAll(']'.codeUnits);`, field.Name) + "\n\n"
			} else {
				resultDart += fmt.Sprintf("\t\tresultDart%s.addAll('%s:'.codeUnits);\n", field.Name, strings.Title(field.Name))
				resultDart += fmt.Sprintf("\t\tresultDart%s.addAll(%s(%s!));\n", field.Name, convFunction, field.Name)
			}
			resultDart += fmt.Sprintf("\t\tbuf.addAll(resultDart%s);\n", field.Name)
			if i != len(object.Fields)-1 {
				resultDart += "\t\tbuf.addAll(','.codeUnits);\n"
			}
		}
		resultDart += "\t\t" + `buf.addAll('}'.codeUnits);` + "\n"
		resultDart += "\t\treturn Uint8List.fromList(buf);\n }\n\n"

		// unmarshal

		resultDart += "\tUnmarshal(Map<String, Uint8List> params) {\n"
		for _, field := range object.Fields {
			//convFunction := conv.ConvFunctionUnmarshalByType(field.Type.Type)
			convFunction := "TODO_CHANGE"
			if field.Type.Type == ast2.Object {
				if field.Type.IsArray {
					resultDart += fmt.Sprintf("\t\t\t"+`var arr = GetArray(params["%s"]!)[1];`+"\n", strings.Title(field.Name))
					resultDart += "\t\t\tfor (int i = 0; i < arr.length; i++) {\n"
					resultDart += "\t\t\tUint8List objBytes = arr[i];\n"
					resultDart += fmt.Sprintf("\t\t\t"+"%s curObj = new %s();\n", field.Type.ObjectName, field.Type.ObjectName)
					resultDart += fmt.Sprintf("\t\t\t"+`%s(curObj,objBytes);`+"\n", convFunction)
					resultDart += fmt.Sprintf("\t\t\t%s!.add(curObj);\n", field.Name)
					resultDart += "\t}\n"
				} else {
					resultDart += fmt.Sprintf("\t\t%s = new %s();\n", field.Name, field.Type.ObjectName)
					resultDart += fmt.Sprintf("\t\t"+`%s(%s!,params["%s"]!);`+"\n", convFunction, field.Name, strings.Title(field.Name))
				}
			} else {
				if field.Type.IsArray {
					resultDart += fmt.Sprintf("\t\t"+`var arr%s = GetArray(params["%s"]!)[1];`+"\n", field.Name, strings.Title(field.Name))
					resultDart += fmt.Sprintf("\t\t%s = List.generate(arr%s.length, (index) => %s(arr%s[index]));\n", field.Name, field.Name, convFunction, field.Name)
				} else {
					resultDart += fmt.Sprintf("\t\t"+`%s = %s(params["%s"]!);`+"\n", field.Name, convFunction, strings.Title(field.Name))
				}
			}

		}

		resultDart += "\t}\n\n"
		resultDart += "}\n\n"
	}
}

func fillMethodsDart() {
	resultDart += `class ConnectParams {
  String host;
  int port;
  bool isEncryptionEnabled;

  ConnectParams(
      {required this.host,
      required this.port,
      required this.isEncryptionEnabled});
}` + "\n\n"

	resultDart += "class Client {\n"
	resultDart += "\tConnectParams connectParams;\n\n"
	resultDart += "\tlate InternalClient internalClient;\n"
	resultDart += "\tClient({required this.connectParams}){\n"
	resultDart += "\t\tinternalClient = InternalClient(host: connectParams.host,port: connectParams.port,isEncryptionEnabled: connectParams.isEncryptionEnabled);\n"
	resultDart += "\t}\n"
	resultDart += "// handlers\n\n"
	for _, service := range astDart.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				resultDart += fmt.Sprintf("\tFuture<%s> %s(%s request) async {\n", method.Returns[0].Type.ObjectName, strings.ToLower(method.Name[:1])+method.Name[1:], method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("\t\t\t%s res = %s();\n", method.Returns[0].Type.ObjectName, method.Returns[0].Type.ObjectName)
				resultDart += fmt.Sprintf("\t\t\tUint8List decoratedMsg = decorateRawDataByHandlerName(\"%s\", request.Marshal());\n", method.Name)
				resultDart += fmt.Sprintf("\t\t\tUint8List rawResponse = await internalClient.handleMessage(\"%s\", decoratedMsg);\n", method.Name)
				resultDart += fmt.Sprintf("\t\t\tMap<String, Uint8List> params = GetParams(rawResponse)[1];\n")
				resultDart += fmt.Sprintf("\t\t\tres.Unmarshal(params);\n")
				resultDart += fmt.Sprintf("\t\t\treturn res;\n")
				resultDart += fmt.Sprintf("\t}\n\n")
			} else if method.MethodType == ast2.Stream {
				resultDart += fmt.Sprintf("\tFuture<void Function(SendMessage msg)> %s(void Function() onEncryptEnabled, void Function(%s msg) onGotMessage, %s initMessage) {\n", strings.ToLower(method.Name[:1])+method.Name[1:], method.Returns[0].Type.ObjectName, method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("\t\treturn internalClient.listenUpdates(\"%s\", onEncryptEnabled, %s(), onGotMessage, initMessage);\n", method.Name, method.Returns[0].Type.ObjectName)
				resultDart += "\t}\n\n"
			}
		}
	}

	resultDart += "}\n\n"
}
