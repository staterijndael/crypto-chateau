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
	fillBinaryCtx()
	fillMethodsDart()
	err := fillObjectsDart()
	if err != nil {
		panic(err)
	}

	return resultDart
}

func fillImportsDart() {
	resultDart += `
import 'dart:async';
import 'dart:typed_data';
import 'package:crypto_chateau_dart/client/models.dart';
import 'package:crypto_chateau_dart/client/conv.dart';
import 'package:crypto_chateau_dart/transport/connection/connection.dart';
import 'package:crypto_chateau_dart/client/binary_iterator.dart';
import 'package:crypto_chateau_dart/transport/handler.dart';
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

func fillBinaryCtx() {
	resultDart += `class BinaryCtx {
  int size;
  int arrSize;
  int pos;
  late BinaryIterator buf;
  late BinaryIterator arrBuf;

  BinaryCtx({
    this.size = 0,
    this.arrSize = 0,
    this.pos = 0,
  }) {
    buf = BinaryIterator(List.empty(growable: true));
    arrBuf = BinaryIterator(List.empty(growable: true));
  }
}` + "\n\n"
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
	resultDart += `extension ExtendList<T> on List<T> {
  void extend(int newLength, T defaultValue) {
    assert(newLength >= 0);

    final lengthDifference = newLength - this.length;
    if (lengthDifference <= 0) {
		this.length = newLength;
        return;
    }

    this.addAll(List.filled(lengthDifference, defaultValue));
  }
}` + "\n\n"

	resultDart += `
class Client {
  final ConnectParams connectParams;
  final MultiplexRequestLoop _pool;

  const Client._({
    required this.connectParams,
    required MultiplexRequestLoop pool,
  }) : _pool = pool;

  factory Client({
    required ConnectParams connectParams,
  }) {
    final encryption = Encryption();
    final connection =
    Connection.root(connectParams).pipe().cipher(encryption).handshake(encryption).multiplex().pipe();

    return Client._(
      connectParams: connectParams,
      pool: MultiplexRequestLoop(connection),
    );
  }
`
	resultDart += "// handlers\n\n"
	for _, service := range astDart.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				resultDart += fmt.Sprintf("\tFuture<%s> %s(%s request) async {\n", method.Returns[0].Type.ObjectName, strings.ToLower(method.Name[:1])+method.Name[1:], method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("\t\t_pool.sendRequest(HandlerHash(hash:[%s]), request, %s(%s));\n", method.Hash.Code(), method.Returns[0].Type.ObjectName, ast2.FillDefaultObjectValues(astDart.Chateau.ObjectDefinitionByObjectName, method.Returns[0].Type.ObjectName))
				resultDart += fmt.Sprintf("\t}\n\n")
			} else if method.MethodType == ast2.Stream {
				resultDart += fmt.Sprintf("\tPeer %s() {\n", strings.ToLower(method.Name[:1])+method.Name[1:])
				resultDart += fmt.Sprintf("\t\treturn peer;\n")
				resultDart += "\t}\n\n"
			}
		}
	}

	resultDart += "}\n\n"
}
