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
	resultDart += "import 'dart:convert';\n"
	resultDart += "import 'dart:async';\n"
	resultDart += "import 'dart:typed_data';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/models.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/conv.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/transport/peer.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/transport/pipe.dart';\n"
	resultDart += "import 'dart:io';\n"
	resultDart += "import 'package:crypto_chateau_dart/client/binary_iterator.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/transport/conn.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/transport/multiplex_conn.dart';\n"
	resultDart += "import 'package:crypto_chateau_dart/transport/handler.dart';\n\n"
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
	resultDart += "\tlate Peer peer;\n"
	resultDart += "\tlate MultiplexConnPool pool;\n"
	resultDart += "\tCompleter<void>? _completer;"
	resultDart += "\tClient({required this.connectParams}){\n"
	resultDart += "\t\t_completer = _createCompleter();\n"
	resultDart += "\t}\n"
	resultDart += `  Completer<void> _createCompleter() {
    _connect();
    return Completer<void>();
  }

  Future<void> _connect() async {
     Socket tcpConn =
        await Socket.connect(connectParams.host, connectParams.port);
    Peer peer = Peer(Pipe(Conn(tcpConn)));
    await peer.establishSecureConn();
    pool = MultiplexConnPool(peer.pipe.tcpConn, true);
    pool.run();
    _completer!.complete();
  }

  Future<void> get connected => _completer!.future;`
	resultDart += "// handlers\n\n"
	for _, service := range astDart.Chateau.Services {
		for _, method := range service.Methods {
			if method.MethodType == ast2.Handler {
				resultDart += fmt.Sprintf("\tFuture<%s> %s(%s request) async {\n", method.Returns[0].Type.ObjectName, strings.ToLower(method.Name[:1])+method.Name[1:], method.Params[0].Type.ObjectName)
				resultDart += fmt.Sprintf("MultiplexConn multiplexConn = pool.newMultiplexConn();\n")
				resultDart += fmt.Sprintf("Peer peer = Peer(Pipe(multiplexConn));\n\n")
				resultDart += fmt.Sprintf("\t\t\tpeer.sendRequestClient(HandlerHash(hash:[%s]), request);\n", method.Hash.Code())
				resultDart += fmt.Sprintf("\t\t\t%s resp = await peer.readMessage(%s(%s)) as %s;\n", method.Returns[0].Type.ObjectName, method.Returns[0].Type.ObjectName, ast2.FillDefaultObjectValues(astDart.Chateau.ObjectDefinitionByObjectName, method.Returns[0].Type.ObjectName), method.Returns[0].Type.ObjectName)
				resultDart += "\t\t\treturn resp;\n"
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
