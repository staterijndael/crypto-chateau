import 'dart:convert';
import 'dart:async';
import 'dart:typed_data';
import 'package:crypto_chateau_dart/client/models.dart';
import 'package:crypto_chateau_dart/client/conv.dart';
import 'package:crypto_chateau_dart/transport/peer.dart';
import 'package:crypto_chateau_dart/transport/pipe.dart';
import 'dart:io';
import 'package:crypto_chateau_dart/client/binary_iterator.dart';import 'package:crypto_chateau_dart/transport/conn.dart';
import 'package:crypto_chateau_dart/transport/handler.dart';

var handlerHashMap = {
	"Reverse":{
		"ReverseMagicString":[0x90, 0xA, 0xDC, 0x45],
		"Rasd":[0xCB, 0xB1, 0x2D, 0x3D],
	},
};

class BinaryCtx {
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
}

class ConnectParams {
  String host;
  int port;
  bool isEncryptionEnabled;

  ConnectParams(
      {required this.host,
      required this.port,
      required this.isEncryptionEnabled});
}

class Client {
	ConnectParams connectParams;

	late Peer peer;
	Completer<void>? _completer;	Client({required this.connectParams}){
		_completer = _createCompleter();
	}
  Completer<void> _createCompleter() {
    _connect();
    return Completer<void>();
  }

  Future<void> _connect() async {
    Socket tcpConn =
        await Socket.connect(connectParams.host, connectParams.port);
    peer = Peer(Pipe(Conn(tcpConn)));
	await peer.establishSecureConn();
    _completer!.complete();
  }

  Future<void> get connected => _completer!.future;// handlers

	Future<ReverseMagicStringResponse> reverseMagicString(ReverseMagicStringRequest request) async {
			peer.sendRequestClient(HandlerHash(hash:[0x90, 0xA, 0xDC, 0x45]), request);
			ReverseMagicStringResponse resp = await peer.readMessage(ReverseMagicStringResponse()) as ReverseMagicStringResponse;
			return resp;
	}

	Future<ReverseMagicStringResponse> rasd(ReverseMagicStringRequest request) async {
			peer.sendRequestClient(HandlerHash(hash:[0xCB, 0xB1, 0x2D, 0x3D]), request);
			ReverseMagicStringResponse resp = await peer.readMessage(ReverseMagicStringResponse()) as ReverseMagicStringResponse;
			return resp;
	}

}



class ReverseCommonObject implements Message {
  List<int>? Key;
  List<String>? Value;

  ReverseCommonObject({
    this.Key,
    this.Value,
  });

  Uint8List Marshal() {
      List<int> b = [];
      int len = 0;
      List<int> arrBufKey = [];
      for (var elKey in Key!) {
	arrBufKey.addAll(ConvertByteToBytes(elKey));
      }
      b.addAll(arrBufKey);
      List<int> arrBufValue = [];
      for (var elValue in Value!) {
	arrBufValue.addAll(ConvertSizeToBytes(elValue.codeUnits.length));
	arrBufValue.addAll(ConvertStringToBytes(elValue));
      }
      b.addAll(arrBufValue);

      return Uint8List.fromList(b);
  }

  

  void Unmarshal(BinaryIterator b) {
  	BinaryCtx binaryCtx = BinaryCtx();
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	binaryCtx.pos = 0;
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   
  	  	       int elKey;
  	  	   
  
      binaryCtx.buf = binaryCtx.arrBuf.slice(1);
      elKey = ConvertBytesToByte(binaryCtx.buf);
  
  
          Key![binaryCtx.pos] = elKey;
  		binaryCtx.pos++;
  	}
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	binaryCtx.pos = 0;
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   
  	  	       String elValue;
  	  	   
  
      binaryCtx.size = binaryCtx.arrBuf.nextSize();
  	binaryCtx.buf = binaryCtx.arrBuf.slice(binaryCtx.size);
      elValue = ConvertBytesToString(binaryCtx.buf);
  
  
          Value![binaryCtx.pos] = elValue;
  		binaryCtx.pos++;
  	}
  }

}

class ReverseMagicStringRequest implements Message {
  String? MagicString;
  int? MagicInt8;
  int? MagicInt16;
  int? MagicInt32;
  int? MagicInt64;
  int? MagicUInt8;
  int? MagicUInt16;
  int? MagicUInt32;
  int? MagicUInt64;
  bool? MagicBool;
  List<int>? MagicBytes;
  ReverseCommonObject? MagicObject;
  List<ReverseCommonObject>? MagicObjectArray;

  ReverseMagicStringRequest({
    this.MagicString,
    this.MagicInt8,
    this.MagicInt16,
    this.MagicInt32,
    this.MagicInt64,
    this.MagicUInt8,
    this.MagicUInt16,
    this.MagicUInt32,
    this.MagicUInt64,
    this.MagicBool,
    this.MagicBytes,
    this.MagicObject,
    this.MagicObjectArray,
  });

  Uint8List Marshal() {
      List<int> b = [];
      int len = 0;
	b.addAll(ConvertSizeToBytes(MagicString!.codeUnits.length));
	b.addAll(ConvertStringToBytes(MagicString!));
	b.addAll(ConvertInt8ToBytes(MagicInt8!));
	b.addAll(ConvertInt16ToBytes(MagicInt16!));
	b.addAll(ConvertInt32ToBytes(MagicInt32!));
	b.addAll(ConvertInt64ToBytes(MagicInt64!));
	b.addAll(ConvertUint8ToBytes(MagicUInt8!));
	b.addAll(ConvertUint16ToBytes(MagicUInt16!));
	b.addAll(ConvertUint32ToBytes(MagicUInt32!));
	b.addAll(ConvertUint64ToBytes(MagicUInt64!));
	b.addAll(ConvertBoolToBytes(MagicBool!));
      List<int> arrBufMagicBytes = [];
      for (var elMagicBytes in MagicBytes!) {
	arrBufMagicBytes.addAll(ConvertByteToBytes(elMagicBytes));
      }
      b.addAll(arrBufMagicBytes);
		b.addAll(MagicObject!.Marshal());
      List<int> arrBufMagicObjectArray = [];
      for (var elMagicObjectArray in MagicObjectArray!) {
		arrBufMagicObjectArray.addAll(elMagicObjectArray.Marshal());
      }
      b.addAll(arrBufMagicObjectArray);

      return Uint8List.fromList(b);
  }

  

  void Unmarshal(BinaryIterator b) {
  	BinaryCtx binaryCtx = BinaryCtx();
      
  
      binaryCtx.size = b.nextSize();
  	binaryCtx.buf = b.slice(binaryCtx.size);
      MagicString = ConvertBytesToString(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicInt8 = ConvertBytesToInt8(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(2);
      MagicInt16 = ConvertBytesToInt16(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(4);
      MagicInt32 = ConvertBytesToInt32(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(8);
      MagicInt64 = ConvertBytesToInt64(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicUInt8 = ConvertBytesToUint8(binaryCtx.buf);
  
      
  
      binaryCtx.buf = b.slice(2);
      MagicUInt16 = ConvertBytesToUint16(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(4);
      MagicUInt32 = ConvertBytesToUint32(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(8);
  	MagicUInt64 = ConvertBytesToUint64(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicBool = ConvertBytesToBool(binaryCtx.buf);
  
  
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   
  	  	       int elMagicBytes;
  	  	   
  
      binaryCtx.buf = binaryCtx.arrBuf.slice(1);
      elMagicBytes = ConvertBytesToByte(binaryCtx.buf);
  
  
          MagicBytes!.add(elMagicBytes);
  		
  	}
      
  
      binaryCtx.size = b.nextSize();
  	binaryCtx.buf = b.slice(binaryCtx.size);
      MagicObject!.Unmarshal(binaryCtx.buf);
  
  
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   ReverseCommonObject elMagicObjectArray = ReverseCommonObject();
  	  	   
  
      binaryCtx.size = binaryCtx.arrBuf.nextSize();
  	binaryCtx.buf = binaryCtx.arrBuf.slice(binaryCtx.size);
      elMagicObjectArray!.Unmarshal(binaryCtx.buf);
  
  
          MagicObjectArray!.add(elMagicObjectArray);
  		
  	}
  }

}

class ReverseMagicStringResponse implements Message {
  String? ReversedMagicString;
  int? MagicInt8;
  int? MagicInt16;
  int? MagicInt32;
  int? MagicInt64;
  int? MagicUInt8;
  int? MagicUInt16;
  int? MagicUInt32;
  int? MagicUInt64;
  bool? MagicBool;
  List<int>? MagicBytes;
  ReverseCommonObject? MagicObject;
  List<ReverseCommonObject>? MagicObjectArray;

  ReverseMagicStringResponse({
    this.ReversedMagicString,
    this.MagicInt8,
    this.MagicInt16,
    this.MagicInt32,
    this.MagicInt64,
    this.MagicUInt8,
    this.MagicUInt16,
    this.MagicUInt32,
    this.MagicUInt64,
    this.MagicBool,
    this.MagicBytes,
    this.MagicObject,
    this.MagicObjectArray,
  });

  Uint8List Marshal() {
      List<int> b = [];
      int len = 0;
	b.addAll(ConvertSizeToBytes(ReversedMagicString!.codeUnits.length));
	b.addAll(ConvertStringToBytes(ReversedMagicString!));
	b.addAll(ConvertInt8ToBytes(MagicInt8!));
	b.addAll(ConvertInt16ToBytes(MagicInt16!));
	b.addAll(ConvertInt32ToBytes(MagicInt32!));
	b.addAll(ConvertInt64ToBytes(MagicInt64!));
	b.addAll(ConvertUint8ToBytes(MagicUInt8!));
	b.addAll(ConvertUint16ToBytes(MagicUInt16!));
	b.addAll(ConvertUint32ToBytes(MagicUInt32!));
	b.addAll(ConvertUint64ToBytes(MagicUInt64!));
	b.addAll(ConvertBoolToBytes(MagicBool!));
      List<int> arrBufMagicBytes = [];
      for (var elMagicBytes in MagicBytes!) {
	arrBufMagicBytes.addAll(ConvertByteToBytes(elMagicBytes));
      }
      b.addAll(arrBufMagicBytes);
		b.addAll(MagicObject!.Marshal());
      List<int> arrBufMagicObjectArray = [];
      for (var elMagicObjectArray in MagicObjectArray!) {
		arrBufMagicObjectArray.addAll(elMagicObjectArray.Marshal());
      }
      b.addAll(arrBufMagicObjectArray);

      return Uint8List.fromList(b);
  }

  

  void Unmarshal(BinaryIterator b) {
  	BinaryCtx binaryCtx = BinaryCtx();
      
  
      binaryCtx.size = b.nextSize();
  	binaryCtx.buf = b.slice(binaryCtx.size);
      ReversedMagicString = ConvertBytesToString(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicInt8 = ConvertBytesToInt8(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(2);
      MagicInt16 = ConvertBytesToInt16(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(4);
      MagicInt32 = ConvertBytesToInt32(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(8);
      MagicInt64 = ConvertBytesToInt64(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicUInt8 = ConvertBytesToUint8(binaryCtx.buf);
  
      
  
      binaryCtx.buf = b.slice(2);
      MagicUInt16 = ConvertBytesToUint16(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(4);
      MagicUInt32 = ConvertBytesToUint32(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(8);
  	MagicUInt64 = ConvertBytesToUint64(binaryCtx.buf);
  
  
      
  
      binaryCtx.buf = b.slice(1);
      MagicBool = ConvertBytesToBool(binaryCtx.buf);
  
  
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   
  	  	       int elMagicBytes;
  	  	   
  
      binaryCtx.buf = binaryCtx.arrBuf.slice(1);
      elMagicBytes = ConvertBytesToByte(binaryCtx.buf);
  
  
          MagicBytes!.add(elMagicBytes);
  		
  	}
      
  
      binaryCtx.size = b.nextSize();
  	binaryCtx.buf = b.slice(binaryCtx.size);
      MagicObject!.Unmarshal(binaryCtx.buf);
  
  
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   
  	  	   ReverseCommonObject elMagicObjectArray = ReverseCommonObject();
  	  	   
  
      binaryCtx.size = binaryCtx.arrBuf.nextSize();
  	binaryCtx.buf = binaryCtx.arrBuf.slice(binaryCtx.size);
      elMagicObjectArray!.Unmarshal(binaryCtx.buf);
  
  
          MagicObjectArray!.add(elMagicObjectArray);
  		
  	}
  }

}