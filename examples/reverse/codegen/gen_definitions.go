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

func getPeerByHandlerName(handlerName string, peer *peer.Peer) interface{} {
	return nil
}

func initHandlers(reverse Reverse) map[string]*server.Handler {
	handlers := make(map[string]*server.Handler)

	handlers["ReverseMagicString"] = &server.Handler{
		CallFuncHandler: ReverseMagicStringSqueeze(reverse.ReverseMagicString),
		HandlerType:     server.HandlerT,
		RequestMsgType:  &ReverseMagicStringRequest{},
	}

	return handlers
}

func NewServer(cfg *server.Config, logger *zap.Logger, reverse Reverse) *server.Server {
	handlers := initHandlers(reverse)

	return server.NewServer(cfg, logger, handlers, getPeerByHandlerName)
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
