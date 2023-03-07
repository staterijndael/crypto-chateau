
import 'dart:async';
import 'dart:typed_data';
import 'package:crypto_chateau_dart/client/models.dart';
import 'package:crypto_chateau_dart/client/conv.dart';
import 'package:crypto_chateau_dart/transport/connection/connection.dart';
import 'package:crypto_chateau_dart/client/binary_iterator.dart';
import 'package:crypto_chateau_dart/transport/handler.dart';
var handlerHashMap = {
	"Reverse":{
		"ReverseMagicString":[0x90, 0xA, 0xDC, 0x45],
		"Rasd":[0xCB, 0xB1, 0x2D, 0x3D],
	},
	"Reverse2":{
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

extension ExtendList<T> on List<T> {
  void extend(int newLength, T defaultValue) {
    assert(newLength >= 0);

    final lengthDifference = newLength - this.length;
    if (lengthDifference <= 0) {
		this.length = newLength;
        return;
    }

    this.addAll(List.filled(lengthDifference, defaultValue));
  }
}


class Client {
  final ConnectParams connectParams;
  final Peer _peer;

  const Client._({
    required this.connectParams,
    required Peer peer,
  }) : _peer = peer;

  factory Client({
    required ConnectParams connectParams,
  }) {
    final encryption = Encryption();
    final connection = Connection.root(connectParams).pipe().cipher(encryption);

    return Client._(
      connectParams: connectParams,
      peer: Peer(
        MultiplexConnection(
          connection,
        ),
      ),
    );
  }

// handlers

	Future<ReverseMagicStringResponse> reverseMagicString(ReverseMagicStringRequest request) => _peer.sendRequest(HandlerHash(hash:[0x90, 0xA, 0xDC, 0x45]), request).then(ReverseMagicStringResponse.fromBytes);

	Stream<ReverseMagicStringResponse> rasd(ReverseMagicStringRequest request) => _peer.sendStreamRequest(HandlerHash(hash:[0xCB, 0xB1, 0x2D, 0x3D]), request).map(ReverseMagicStringResponse.fromBytes);

}



class ReverseCommonObject implements Message {
  List<int>? Key;
  List<String>? Value;
  ReverseCommonObject({
    required this.Key,
    required this.Value,
  });
  

  ReverseCommonObject Copy() => ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true));

  static ReverseCommonObject fromBytes(Uint8List bytes) => ReverseCommonObject(
        Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true))..Unmarshal(BinaryIterator(bytes));

  Uint8List Marshal() {
      List<int> b = [];

      List<int> size = ConvertSizeToBytes(0);
      b.addAll(size);
      List<int> arrBufKey = [];
      for (var elKey in Key!) {
	arrBufKey.addAll(ConvertByteToBytes(elKey));
      }
      b.addAll(ConvertSizeToBytes(arrBufKey.length));
      b.addAll(arrBufKey);
      List<int> arrBufValue = [];
      for (var elValue in Value!) {
	arrBufValue.addAll(ConvertSizeToBytes(elValue.codeUnits.length));
	arrBufValue.addAll(ConvertStringToBytes(elValue));
      }
      b.addAll(ConvertSizeToBytes(arrBufValue.length));
      b.addAll(arrBufValue);
      size = ConvertSizeToBytes(b.length - size.length);
      for (int i = 0; i < size.length; i++) {
      	b[i] = size[i];
      }

      return Uint8List.fromList(b);
  }

  

  void Unmarshal(BinaryIterator b) {
  	BinaryCtx binaryCtx = BinaryCtx();
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	binaryCtx.pos = 0;

    
    bool isEmptyKey = true;
    
  	while (binaryCtx.arrBuf.hasNext()) {
  	       
  	       isEmptyKey = false;
  	       
  	  	   
  	  	   
  	  	       int elKey;
  	  	   
  
      binaryCtx.buf = binaryCtx.arrBuf.slice(1);
      elKey = ConvertBytesToByte(binaryCtx.buf);
  
  
           Key!.add(elKey);
  	}

    
  	if (isEmptyKey){
  	    Key = null;
  	}
  	
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	binaryCtx.pos = 0;

    
    bool isEmptyValue = true;
    
  	while (binaryCtx.arrBuf.hasNext()) {
  	       
  	       isEmptyValue = false;
  	       
  	  	   
  	  	   
  	  	       String elValue;
  	  	   
  
      binaryCtx.size = binaryCtx.arrBuf.nextSize();
  	binaryCtx.buf = binaryCtx.arrBuf.slice(binaryCtx.size);
      elValue = ConvertBytesToString(binaryCtx.buf);
  
  
           Value!.add(elValue);
  	}

    
  	if (isEmptyValue){
  	    Value = null;
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
    required this.MagicString,
    required this.MagicInt8,
    required this.MagicInt16,
    required this.MagicInt32,
    required this.MagicInt64,
    required this.MagicUInt8,
    required this.MagicUInt16,
    required this.MagicUInt32,
    required this.MagicUInt64,
    required this.MagicBool,
    required this.MagicBytes,
    required this.MagicObject,
    required this.MagicObjectArray,
  });
  

  ReverseMagicStringRequest Copy() => ReverseMagicStringRequest(MagicString: "",MagicInt8: 0,MagicInt16: 0,MagicInt32: 0,MagicInt64: 0,MagicUInt8: 0,MagicUInt16: 0,MagicUInt32: 0,MagicUInt64: 0,MagicBool: true,MagicBytes: List.filled(0, 0xff, growable: true),MagicObject: ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)),MagicObjectArray: List.filled(0, ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)), growable: true));

  static ReverseMagicStringRequest fromBytes(Uint8List bytes) => ReverseMagicStringRequest(
        MagicString: "",MagicInt8: 0,MagicInt16: 0,MagicInt32: 0,MagicInt64: 0,MagicUInt8: 0,MagicUInt16: 0,MagicUInt32: 0,MagicUInt64: 0,MagicBool: true,MagicBytes: List.filled(0, 0xff, growable: true),MagicObject: ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)),MagicObjectArray: List.filled(0, ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)), growable: true))..Unmarshal(BinaryIterator(bytes));

  Uint8List Marshal() {
      List<int> b = [];

      List<int> size = ConvertSizeToBytes(0);
      b.addAll(size);
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
      b.addAll(ConvertSizeToBytes(arrBufMagicBytes.length));
      b.addAll(arrBufMagicBytes);
		b.addAll(MagicObject!.Marshal());
      List<int> arrBufMagicObjectArray = [];
      for (var elMagicObjectArray in MagicObjectArray!) {
		arrBufMagicObjectArray.addAll(elMagicObjectArray.Marshal());
      }
      b.addAll(ConvertSizeToBytes(arrBufMagicObjectArray.length));
      b.addAll(arrBufMagicObjectArray);
      size = ConvertSizeToBytes(b.length - size.length);
      for (int i = 0; i < size.length; i++) {
      	b[i] = size[i];
      }

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
  	binaryCtx.pos = 0;

    
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
  	binaryCtx.pos = 0;

    
  	while (binaryCtx.arrBuf.hasNext()) {
  	       
  	  	   
  	  	   ReverseCommonObject elMagicObjectArray = ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true));
  	  	   
  
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
    required this.ReversedMagicString,
    required this.MagicInt8,
    required this.MagicInt16,
    required this.MagicInt32,
    required this.MagicInt64,
    required this.MagicUInt8,
    required this.MagicUInt16,
    required this.MagicUInt32,
    required this.MagicUInt64,
    required this.MagicBool,
    required this.MagicBytes,
    required this.MagicObject,
    required this.MagicObjectArray,
  });
  

  ReverseMagicStringResponse Copy() => ReverseMagicStringResponse(ReversedMagicString: "",MagicInt8: 0,MagicInt16: 0,MagicInt32: 0,MagicInt64: 0,MagicUInt8: 0,MagicUInt16: 0,MagicUInt32: 0,MagicUInt64: 0,MagicBool: true,MagicBytes: List.filled(0, 0xff, growable: true),MagicObject: ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)),MagicObjectArray: List.filled(0, ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)), growable: true));

  static ReverseMagicStringResponse fromBytes(Uint8List bytes) => ReverseMagicStringResponse(
        ReversedMagicString: "",MagicInt8: 0,MagicInt16: 0,MagicInt32: 0,MagicInt64: 0,MagicUInt8: 0,MagicUInt16: 0,MagicUInt32: 0,MagicUInt64: 0,MagicBool: true,MagicBytes: List.filled(0, 0xff, growable: true),MagicObject: ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)),MagicObjectArray: List.filled(0, ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true)), growable: true))..Unmarshal(BinaryIterator(bytes));

  Uint8List Marshal() {
      List<int> b = [];

      List<int> size = ConvertSizeToBytes(0);
      b.addAll(size);
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
      b.addAll(ConvertSizeToBytes(arrBufMagicBytes.length));
      b.addAll(arrBufMagicBytes);
		b.addAll(MagicObject!.Marshal());
      List<int> arrBufMagicObjectArray = [];
      for (var elMagicObjectArray in MagicObjectArray!) {
		arrBufMagicObjectArray.addAll(elMagicObjectArray.Marshal());
      }
      b.addAll(ConvertSizeToBytes(arrBufMagicObjectArray.length));
      b.addAll(arrBufMagicObjectArray);
      size = ConvertSizeToBytes(b.length - size.length);
      for (int i = 0; i < size.length; i++) {
      	b[i] = size[i];
      }

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
  	binaryCtx.pos = 0;

    
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
  	binaryCtx.pos = 0;

    
  	while (binaryCtx.arrBuf.hasNext()) {
  	       
  	  	   
  	  	   ReverseCommonObject elMagicObjectArray = ReverseCommonObject(Key: List.filled(0, 0xff, growable: true),Value: List.filled(0, "", growable: true));
  	  	   
  
      binaryCtx.size = binaryCtx.arrBuf.nextSize();
  	binaryCtx.buf = binaryCtx.arrBuf.slice(binaryCtx.size);
      elMagicObjectArray!.Unmarshal(binaryCtx.buf);
  
  
           MagicObjectArray!.add(elMagicObjectArray);
  	}

    
  }

}