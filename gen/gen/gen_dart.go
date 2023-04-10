package gen

import (
	"fmt"
	ast2 "github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/gen/templates"
	"github.com/pkg/errors"
	"strings"
)

var resultDart string
var astDart *ast2.Ast

func GenerateDefinitionsDart(astLocal *ast2.Ast) string {
	astDart = astLocal

	fillImportsDart()
	fillHandlerHash()
	fillMethodsDart()
	err := fillObjectsDart()
	if err != nil {
		panic(err)
	}

	return resultDart
}

func fillImportsDart() {
	resultDart += `
part of 'client.dart';
`
}

func fillHandlerHash() {
	resultDart += "var handlerHashMap = {\n"
	for _, service := range astDart.Chateau.Services {
		resultDart += "\t\"" + service.Name + "\":{\n"
		for _, method := range service.Methods {
			resultDart += "\t\t\"" + method.Name + "\":[" + method.Hash.Code() + "],\n"
		}
		resultDart += "\t},\n"
	}
	resultDart += "};\n\n"
}

func fillObjectsDart() error {
	ot, err := templates.NewObjectTemplateDart()
	if err != nil {
		return errors.Wrap(err, "failed to create object dart template")
	}

	var genObject string
	for _, object := range astDart.Chateau.ObjectDefinitions {
		genObject, err = ot.GenDart(object, astDart.Chateau.ObjectDefinitionByObjectName)
		if err != nil {
			return errors.Wrap(err, "failed to generate object")
		}

		resultDart += genObject
	}

	return nil
}

func fillMethodsDart() {
	resultDart += `
mixin ClientMixin {
  Peer get _peer;

`
	for _, service := range astDart.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				resultDart += fmt.Sprintf("\tFuture<%s> %s(%s request) => ", method.Returns[0].Type.ObjectName, strings.ToLower(method.Name[:1])+method.Name[1:], method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("_peer.request(HandlerHash(hash:[%s]), request).first.then(%s.fromBytes);\n\n", method.Hash.Code(), method.Returns[0].Type.ObjectName)
			} else if method.MethodType == ast2.Stream {
				resultDart += fmt.Sprintf("\tStream<%s> %s(%s request, Stream<Message> writePipe) => ", method.Returns[0].Type.ObjectName, strings.ToLower(method.Name[:1])+method.Name[1:], method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("_peer.request(HandlerHash(hash:[%s]), request, writePipe).map(%s.fromBytes);\n\n", method.Hash.Code(), method.Returns[0].Type.ObjectName)
			}
		}
	}

	resultDart += "}\n\n"
}
