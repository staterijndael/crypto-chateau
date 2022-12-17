// CODEGEN VERSION: v1.0

package endpoints

import (
	"context"
	"errors"
	"net"
	"strconv"

	"go.uber.org/zap"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/peer"
	"github.com/oringik/crypto-chateau/server"
	"github.com/oringik/crypto-chateau/transport"
)

var tagsByHandlerName = map[string]map[string]string{
	"ReverseMagicString": {"key": "val", "ajsdajsd": "asdasd"},
}

type Reverse interface {
	ReverseMagicString(ctx context.Context, req *ReverseMagicStringRequest) (*ReverseMagicStringResponse, error)
}

func ReverseMagicStringSqueeze(fnc func(context.Context, *ReverseMagicStringRequest) (*ReverseMagicStringResponse, error)) server.HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*ReverseMagicStringRequest); ok {
			return fnc(ctx, msg.(*ReverseMagicStringRequest))
		} else {
			return nil, errors.New("unknown message type: expected ReverseMagicStringRequest")
		}
	}
}

type ReverseCommonObject struct {
	Key   [16]byte
	Value [32]string
}

func (o *ReverseCommonObject) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultKey []byte
	resultKey = append(resultKey, []byte("Key:")...)
	resultKey = append(resultKey, '[')
	for i, val := range o.Key {
		resultKey = append(resultKey, conv.ConvertByteToBytes(val)...)
		if i != len(o.Key)-1 {
			resultKey = append(resultKey, ',')
		}
	}
	resultKey = append(resultKey, ']')

	buf = append(buf, resultKey...)
	buf = append(buf, ',')
	var resultValue []byte
	resultValue = append(resultValue, []byte("Value:")...)
	resultValue = append(resultValue, '[')
	for i, val := range o.Value {
		resultValue = append(resultValue, conv.ConvertStringToBytes(val)...)
		if i != len(o.Value)-1 {
			resultValue = append(resultValue, ',')
		}
	}
	resultValue = append(resultValue, ']')

	buf = append(buf, resultValue...)
	buf = append(buf, '}')
	return buf
}

func (o *ReverseCommonObject) Unmarshal(params map[string][]byte) error {
	_, arr, err := conv.GetArray(params["Key"])
	if err != nil {
		return err
	}
	for i, valByte := range arr {
		o.Key[i] = conv.ConvertBytesToByte(valByte)
	}
	_, arr, err = conv.GetArray(params["Value"])
	if err != nil {
		return err
	}
	for i, valByte := range arr {
		o.Value[i] = conv.ConvertBytesToString(valByte)
	}
	return nil
}

type ReverseMagicStringRequest struct {
	MagicString      string
	MagicInt8        int8
	MagicInt16       int16
	MagicInt32       int32
	MagicInt64       int64
	MagicUInt8       uint8
	MagicUInt16      uint16
	MagicUInt32      uint32
	MagicUInt64      uint64
	MagicBool        bool
	MagicBytes       []byte
	MagicObject      *ReverseCommonObject
	MagicObjectArray []*ReverseCommonObject
}

func (o *ReverseMagicStringRequest) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultMagicString []byte
	resultMagicString = append(resultMagicString, []byte("MagicString:")...)
	resultMagicString = append(resultMagicString, conv.ConvertStringToBytes(o.MagicString)...)
	buf = append(buf, resultMagicString...)
	buf = append(buf, ',')
	var resultMagicInt8 []byte
	resultMagicInt8 = append(resultMagicInt8, []byte("MagicInt8:")...)
	resultMagicInt8 = append(resultMagicInt8, conv.ConvertInt8ToBytes(o.MagicInt8)...)
	buf = append(buf, resultMagicInt8...)
	buf = append(buf, ',')
	var resultMagicInt16 []byte
	resultMagicInt16 = append(resultMagicInt16, []byte("MagicInt16:")...)
	resultMagicInt16 = append(resultMagicInt16, conv.ConvertInt16ToBytes(o.MagicInt16)...)
	buf = append(buf, resultMagicInt16...)
	buf = append(buf, ',')
	var resultMagicInt32 []byte
	resultMagicInt32 = append(resultMagicInt32, []byte("MagicInt32:")...)
	resultMagicInt32 = append(resultMagicInt32, conv.ConvertInt32ToBytes(o.MagicInt32)...)
	buf = append(buf, resultMagicInt32...)
	buf = append(buf, ',')
	var resultMagicInt64 []byte
	resultMagicInt64 = append(resultMagicInt64, []byte("MagicInt64:")...)
	resultMagicInt64 = append(resultMagicInt64, conv.ConvertInt64ToBytes(o.MagicInt64)...)
	buf = append(buf, resultMagicInt64...)
	buf = append(buf, ',')
	var resultMagicUInt8 []byte
	resultMagicUInt8 = append(resultMagicUInt8, []byte("MagicUInt8:")...)
	resultMagicUInt8 = append(resultMagicUInt8, conv.ConvertUint8ToBytes(o.MagicUInt8)...)
	buf = append(buf, resultMagicUInt8...)
	buf = append(buf, ',')
	var resultMagicUInt16 []byte
	resultMagicUInt16 = append(resultMagicUInt16, []byte("MagicUInt16:")...)
	resultMagicUInt16 = append(resultMagicUInt16, conv.ConvertUint16ToBytes(o.MagicUInt16)...)
	buf = append(buf, resultMagicUInt16...)
	buf = append(buf, ',')
	var resultMagicUInt32 []byte
	resultMagicUInt32 = append(resultMagicUInt32, []byte("MagicUInt32:")...)
	resultMagicUInt32 = append(resultMagicUInt32, conv.ConvertUint32ToBytes(o.MagicUInt32)...)
	buf = append(buf, resultMagicUInt32...)
	buf = append(buf, ',')
	var resultMagicUInt64 []byte
	resultMagicUInt64 = append(resultMagicUInt64, []byte("MagicUInt64:")...)
	resultMagicUInt64 = append(resultMagicUInt64, conv.ConvertUint64ToBytes(o.MagicUInt64)...)
	buf = append(buf, resultMagicUInt64...)
	buf = append(buf, ',')
	var resultMagicBool []byte
	resultMagicBool = append(resultMagicBool, []byte("MagicBool:")...)
	resultMagicBool = append(resultMagicBool, conv.ConvertBoolToBytes(o.MagicBool)...)
	buf = append(buf, resultMagicBool...)
	buf = append(buf, ',')
	var resultMagicBytes []byte
	resultMagicBytes = append(resultMagicBytes, []byte("MagicBytes:")...)
	resultMagicBytes = append(resultMagicBytes, '[')
	for i, val := range o.MagicBytes {
		resultMagicBytes = append(resultMagicBytes, conv.ConvertByteToBytes(val)...)
		if i != len(o.MagicBytes)-1 {
			resultMagicBytes = append(resultMagicBytes, ',')
		}
	}
	resultMagicBytes = append(resultMagicBytes, ']')

	buf = append(buf, resultMagicBytes...)
	buf = append(buf, ',')
	var resultMagicObject []byte
	resultMagicObject = append(resultMagicObject, []byte("MagicObject:")...)
	resultMagicObject = append(resultMagicObject, conv.ConvertObjectToBytes(o.MagicObject)...)
	buf = append(buf, resultMagicObject...)
	buf = append(buf, ',')
	var resultMagicObjectArray []byte
	resultMagicObjectArray = append(resultMagicObjectArray, []byte("MagicObjectArray:")...)
	resultMagicObjectArray = append(resultMagicObjectArray, '[')
	for i, val := range o.MagicObjectArray {
		resultMagicObjectArray = append(resultMagicObjectArray, conv.ConvertObjectToBytes(val)...)
		if i != len(o.MagicObjectArray)-1 {
			resultMagicObjectArray = append(resultMagicObjectArray, ',')
		}
	}
	resultMagicObjectArray = append(resultMagicObjectArray, ']')

	buf = append(buf, resultMagicObjectArray...)
	buf = append(buf, '}')
	return buf
}

func (o *ReverseMagicStringRequest) Unmarshal(params map[string][]byte) error {
	o.MagicString = conv.ConvertBytesToString(params["MagicString"])
	o.MagicInt8 = conv.ConvertBytesToInt8(params["MagicInt8"])
	o.MagicInt16 = conv.ConvertBytesToInt16(params["MagicInt16"])
	o.MagicInt32 = conv.ConvertBytesToInt32(params["MagicInt32"])
	o.MagicInt64 = conv.ConvertBytesToInt64(params["MagicInt64"])
	o.MagicUInt8 = conv.ConvertBytesToUint8(params["MagicUInt8"])
	o.MagicUInt16 = conv.ConvertBytesToUint16(params["MagicUInt16"])
	o.MagicUInt32 = conv.ConvertBytesToUint32(params["MagicUInt32"])
	o.MagicUInt64 = conv.ConvertBytesToUint64(params["MagicUInt64"])
	o.MagicBool = conv.ConvertBytesToBool(params["MagicBool"])
	_, arr, err := conv.GetArray(params["MagicBytes"])
	if err != nil {
		return err
	}
	for _, valByte := range arr {
		var curVal byte
		curVal = conv.ConvertBytesToByte(valByte)
		o.MagicBytes = append(o.MagicBytes, curVal)
	}
	o.MagicObject = &ReverseCommonObject{}
	conv.ConvertBytesToObject(o.MagicObject, params["MagicObject"])
	_, arr, err = conv.GetArray(params["MagicObjectArray"])
	if err != nil {
		return err
	}
	for _, objBytes := range arr {
		var curObj *ReverseCommonObject
		conv.ConvertBytesToObject(curObj, objBytes)
		o.MagicObjectArray = append(o.MagicObjectArray, curObj)
	}
	return nil
}

type ReverseMagicStringResponse struct {
	ReversedMagicString string
	MagicInt8           int8
	MagicInt16          int16
	MagicInt32          int32
	MagicInt64          int64
	MagicUInt8          uint8
	MagicUInt16         uint16
	MagicUInt32         uint32
	MagicUInt64         uint64
	MagicBool           bool
	MagicBytes          []byte
	MagicObject         *ReverseCommonObject
	MagicObjectArray    []*ReverseCommonObject
}

func (o *ReverseMagicStringResponse) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultReversedMagicString []byte
	resultReversedMagicString = append(resultReversedMagicString, []byte("ReversedMagicString:")...)
	resultReversedMagicString = append(resultReversedMagicString, conv.ConvertStringToBytes(o.ReversedMagicString)...)
	buf = append(buf, resultReversedMagicString...)
	buf = append(buf, ',')
	var resultMagicInt8 []byte
	resultMagicInt8 = append(resultMagicInt8, []byte("MagicInt8:")...)
	resultMagicInt8 = append(resultMagicInt8, conv.ConvertInt8ToBytes(o.MagicInt8)...)
	buf = append(buf, resultMagicInt8...)
	buf = append(buf, ',')
	var resultMagicInt16 []byte
	resultMagicInt16 = append(resultMagicInt16, []byte("MagicInt16:")...)
	resultMagicInt16 = append(resultMagicInt16, conv.ConvertInt16ToBytes(o.MagicInt16)...)
	buf = append(buf, resultMagicInt16...)
	buf = append(buf, ',')
	var resultMagicInt32 []byte
	resultMagicInt32 = append(resultMagicInt32, []byte("MagicInt32:")...)
	resultMagicInt32 = append(resultMagicInt32, conv.ConvertInt32ToBytes(o.MagicInt32)...)
	buf = append(buf, resultMagicInt32...)
	buf = append(buf, ',')
	var resultMagicInt64 []byte
	resultMagicInt64 = append(resultMagicInt64, []byte("MagicInt64:")...)
	resultMagicInt64 = append(resultMagicInt64, conv.ConvertInt64ToBytes(o.MagicInt64)...)
	buf = append(buf, resultMagicInt64...)
	buf = append(buf, ',')
	var resultMagicUInt8 []byte
	resultMagicUInt8 = append(resultMagicUInt8, []byte("MagicUInt8:")...)
	resultMagicUInt8 = append(resultMagicUInt8, conv.ConvertUint8ToBytes(o.MagicUInt8)...)
	buf = append(buf, resultMagicUInt8...)
	buf = append(buf, ',')
	var resultMagicUInt16 []byte
	resultMagicUInt16 = append(resultMagicUInt16, []byte("MagicUInt16:")...)
	resultMagicUInt16 = append(resultMagicUInt16, conv.ConvertUint16ToBytes(o.MagicUInt16)...)
	buf = append(buf, resultMagicUInt16...)
	buf = append(buf, ',')
	var resultMagicUInt32 []byte
	resultMagicUInt32 = append(resultMagicUInt32, []byte("MagicUInt32:")...)
	resultMagicUInt32 = append(resultMagicUInt32, conv.ConvertUint32ToBytes(o.MagicUInt32)...)
	buf = append(buf, resultMagicUInt32...)
	buf = append(buf, ',')
	var resultMagicUInt64 []byte
	resultMagicUInt64 = append(resultMagicUInt64, []byte("MagicUInt64:")...)
	resultMagicUInt64 = append(resultMagicUInt64, conv.ConvertUint64ToBytes(o.MagicUInt64)...)
	buf = append(buf, resultMagicUInt64...)
	buf = append(buf, ',')
	var resultMagicBool []byte
	resultMagicBool = append(resultMagicBool, []byte("MagicBool:")...)
	resultMagicBool = append(resultMagicBool, conv.ConvertBoolToBytes(o.MagicBool)...)
	buf = append(buf, resultMagicBool...)
	buf = append(buf, ',')
	var resultMagicBytes []byte
	resultMagicBytes = append(resultMagicBytes, []byte("MagicBytes:")...)
	resultMagicBytes = append(resultMagicBytes, '[')
	for i, val := range o.MagicBytes {
		resultMagicBytes = append(resultMagicBytes, conv.ConvertByteToBytes(val)...)
		if i != len(o.MagicBytes)-1 {
			resultMagicBytes = append(resultMagicBytes, ',')
		}
	}
	resultMagicBytes = append(resultMagicBytes, ']')

	buf = append(buf, resultMagicBytes...)
	buf = append(buf, ',')
	var resultMagicObject []byte
	resultMagicObject = append(resultMagicObject, []byte("MagicObject:")...)
	resultMagicObject = append(resultMagicObject, conv.ConvertObjectToBytes(o.MagicObject)...)
	buf = append(buf, resultMagicObject...)
	buf = append(buf, ',')
	var resultMagicObjectArray []byte
	resultMagicObjectArray = append(resultMagicObjectArray, []byte("MagicObjectArray:")...)
	resultMagicObjectArray = append(resultMagicObjectArray, '[')
	for i, val := range o.MagicObjectArray {
		resultMagicObjectArray = append(resultMagicObjectArray, conv.ConvertObjectToBytes(val)...)
		if i != len(o.MagicObjectArray)-1 {
			resultMagicObjectArray = append(resultMagicObjectArray, ',')
		}
	}
	resultMagicObjectArray = append(resultMagicObjectArray, ']')

	buf = append(buf, resultMagicObjectArray...)
	buf = append(buf, '}')
	return buf
}

func (o *ReverseMagicStringResponse) Unmarshal(params map[string][]byte) error {
	o.ReversedMagicString = conv.ConvertBytesToString(params["ReversedMagicString"])
	o.MagicInt8 = conv.ConvertBytesToInt8(params["MagicInt8"])
	o.MagicInt16 = conv.ConvertBytesToInt16(params["MagicInt16"])
	o.MagicInt32 = conv.ConvertBytesToInt32(params["MagicInt32"])
	o.MagicInt64 = conv.ConvertBytesToInt64(params["MagicInt64"])
	o.MagicUInt8 = conv.ConvertBytesToUint8(params["MagicUInt8"])
	o.MagicUInt16 = conv.ConvertBytesToUint16(params["MagicUInt16"])
	o.MagicUInt32 = conv.ConvertBytesToUint32(params["MagicUInt32"])
	o.MagicUInt64 = conv.ConvertBytesToUint64(params["MagicUInt64"])
	o.MagicBool = conv.ConvertBytesToBool(params["MagicBool"])
	_, arr, err := conv.GetArray(params["MagicBytes"])
	if err != nil {
		return err
	}
	for _, valByte := range arr {
		var curVal byte
		curVal = conv.ConvertBytesToByte(valByte)
		o.MagicBytes = append(o.MagicBytes, curVal)
	}
	o.MagicObject = &ReverseCommonObject{}
	conv.ConvertBytesToObject(o.MagicObject, params["MagicObject"])
	_, arr, err = conv.GetArray(params["MagicObjectArray"])
	if err != nil {
		return err
	}
	for _, objBytes := range arr {
		var curObj *ReverseCommonObject
		conv.ConvertBytesToObject(curObj, objBytes)
		o.MagicObjectArray = append(o.MagicObjectArray, curObj)
	}
	return nil
}

func GetHandlers(reverse Reverse) map[string]*server.Handler {
	handlers := make(map[string]*server.Handler)

	var callFuncReverseMagicString server.HandlerFunc
	if reverse != nil {
		callFuncReverseMagicString = ReverseMagicStringSqueeze(reverse.ReverseMagicString)
	}

	handlers["ReverseMagicString"] = &server.Handler{
		CallFuncHandler: callFuncReverseMagicString,
		HandlerType:     server.HandlerT,
		RequestMsgType:  &ReverseMagicStringRequest{},
		ResponseMsgType: &ReverseMagicStringResponse{},
		Tags:            tagsByHandlerName["ReverseMagicString"],
	}

	return handlers
}

func GetEmptyHandlers() map[string]*server.Handler {
	handlers := make(map[string]*server.Handler)

	handlers["ReverseMagicString"] = &server.Handler{
		HandlerType:     server.HandlerT,
		RequestMsgType:  &ReverseMagicStringRequest{},
		ResponseMsgType: &ReverseMagicStringResponse{},
	}

	return handlers
}

func NewServer(cfg *server.Config, logger *zap.Logger, reverse Reverse) *server.Server {
	handlers := GetHandlers(reverse)

	return server.NewServer(cfg, logger, handlers)
}

func CallClientMethod(ctx context.Context, host string, port int, serviceName string, methodName string, req message.Message) (message.Message, error) {
	if serviceName == "Reverse" {
		if methodName == "ReverseMagicString" {
			client, err := NewClientReverse(host, port)
			if err != nil {
				return nil, err
			}
			return client.ReverseMagicString(ctx, req.(*ReverseMagicStringRequest))
		}
	}

	return nil, errors.New("unknown service or method")
}

type ClientReverse struct {
	peer *peer.Peer
}

func NewClientReverse(host string, port int) (*ClientReverse, error) {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err = transport.ServerHandshake(conn)
	if err != nil {
		return nil, err
	}
	securedPeer := peer.NewPeer(conn)
	client := &ClientReverse{peer: securedPeer}
	return client, nil
}

func (c *ClientReverse) ReverseMagicString(ctx context.Context, req *ReverseMagicStringRequest) (*ReverseMagicStringResponse, error) {
	err := c.peer.WriteResponse("ReverseMagicString", req)

	msg := make([]byte, 0, 1024)

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

	respMsg := &ReverseMagicStringResponse{}

	err = respMsg.Unmarshal(responseMsgParams)
	if err != nil {
		return nil, err
	}

	return respMsg, nil
}
