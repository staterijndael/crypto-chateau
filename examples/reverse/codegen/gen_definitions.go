// CODEGEN VERSION: v1.0

package endpoints

import "errors"
import "context"
import "strconv"
import "github.com/oringik/crypto-chateau/gen/conv"
import "github.com/oringik/crypto-chateau/peer"
import "github.com/oringik/crypto-chateau/message"
import "github.com/oringik/crypto-chateau/server"
import "go.uber.org/zap"
import "github.com/oringik/crypto-chateau/transport"
import "net"

type Reverse interface {
	ReverseMagicString(ctx context.Context, peer *peer.Peer, req *ReverseMagicStringRequest) error
}

func ReverseMagicStringSqueeze(fnc func(context.Context, *peer.Peer, *ReverseMagicStringRequest) error) server.StreamFunc {
	return func(ctx context.Context, peer *peer.Peer, msg message.Message) error {
		if _, ok := msg.(*ReverseMagicStringRequest); ok {
			return fnc(ctx, peer, msg.(*ReverseMagicStringRequest))
		} else {
			return errors.New("unknown message type: expected ReverseMagicStringRequest")
		}
	}
}

type ReverseMagicStringRequest struct {
	MagicString string
}

func (o *ReverseMagicStringRequest) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultMagicString []byte
	resultMagicString = append(resultMagicString, []byte("MagicString:")...)
	resultMagicString = append(resultMagicString, conv.ConvertStringToBytes(o.MagicString)...)
	buf = append(buf, resultMagicString...)
	buf = append(buf, '}')
	return buf
}

func (o *ReverseMagicStringRequest) Unmarshal(params map[string][]byte) error {
	o.MagicString = conv.ConvertBytesToString(params["MagicString"])
	return nil
}

type ReverseMagicStringResponse struct {
	ReversedMagicString string
}

func (o *ReverseMagicStringResponse) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultReversedMagicString []byte
	resultReversedMagicString = append(resultReversedMagicString, []byte("ReversedMagicString:")...)
	resultReversedMagicString = append(resultReversedMagicString, conv.ConvertStringToBytes(o.ReversedMagicString)...)
	buf = append(buf, resultReversedMagicString...)
	buf = append(buf, '}')
	return buf
}

func (o *ReverseMagicStringResponse) Unmarshal(params map[string][]byte) error {
	o.ReversedMagicString = conv.ConvertBytesToString(params["ReversedMagicString"])
	return nil
}

func GetHandlers(reverse Reverse) map[string]*server.Handler {
	handlers := make(map[string]*server.Handler)

	var callFuncReverseMagicString server.StreamFunc
	if reverse != nil {
		callFuncReverseMagicString = ReverseMagicStringSqueeze(reverse.ReverseMagicString)
	}

	handlers["ReverseMagicString"] = &server.Handler{
		CallFuncStream:  callFuncReverseMagicString,
		HandlerType:     server.StreamT,
		RequestMsgType:  &ReverseMagicStringRequest{},
		ResponseMsgType: &ReverseMagicStringResponse{},
	}

	return handlers
}

func GetEmptyHandlers() map[string]*server.Handler {
	handlers := make(map[string]*server.Handler)

	handlers["ReverseMagicString"] = &server.Handler{
		HandlerType:     server.StreamT,
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
