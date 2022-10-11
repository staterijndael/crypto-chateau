package endpoints

import (
	"errors"
	"go.uber.org/zap"
)
import "context"
import "github.com/Oringik/crypto-chateau/gen/conv"
import "github.com/Oringik/crypto-chateau/peer"
import "github.com/Oringik/crypto-chateau/message"
import crypto_chateau "github.com/Oringik/crypto-chateau"

type UserEndpoint interface {
	GetUserAHAAHAHHAHAH(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	GetUserArar(ctx context.Context, peer *PeerGetUserArar, asdasd *GetUserRequest) error
}

func GetUserAHAAHAHHAHAHSqueeze(fnc func(context.Context, *GetUserRequest) (*GetUserResponse, error)) crypto_chateau.HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*GetUserRequest); ok {
			return fnc(ctx, msg.(*GetUserRequest))
		} else {
			return nil, errors.New("unknown message type: expected GetUserRequest")
		}
	}
}

func GetUserArarSqueeze(fnc func(context.Context, *PeerGetUserArar, *GetUserRequest) error) crypto_chateau.StreamFunc {
	return func(ctx context.Context, peer interface{}, msg message.Message) error {
		if _, ok := msg.(*GetUserRequest); ok {
			return fnc(ctx, peer.(*PeerGetUserArar), msg.(*GetUserRequest))
		} else {
			return errors.New("unknown message type: expected GetUserRequest")
		}
	}
}

type GetUserRequest struct {
	IdentityKey [32]byte
	Aa          uint32
	Oo          uint64
	User        *User
}

func (o *GetUserRequest) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultIdentityKey []byte
	resultIdentityKey = append(resultIdentityKey, []byte("IdentityKey:")...)
	resultIdentityKey = append(resultIdentityKey, '[')
	for i, val := range o.IdentityKey {
		resultIdentityKey = append(resultIdentityKey, conv.ConvertByteToBytes(val)...)
		if i != len(o.IdentityKey)-1 {
			resultIdentityKey = append(resultIdentityKey, ',')
		}
	}
	resultIdentityKey = append(resultIdentityKey, ']')

	buf = append(buf, resultIdentityKey...)
	buf = append(buf, ',')
	var resultAa []byte
	resultAa = append(resultAa, []byte("Aa:")...)
	resultAa = append(resultAa, conv.ConvertUint32ToBytes(o.Aa)...)
	buf = append(buf, resultAa...)
	buf = append(buf, ',')
	var resultOo []byte
	resultOo = append(resultOo, []byte("Oo:")...)
	resultOo = append(resultOo, conv.ConvertUint64ToBytes(o.Oo)...)
	buf = append(buf, resultOo...)
	buf = append(buf, ',')
	var resultUser []byte
	resultUser = append(resultUser, []byte("User:")...)
	resultUser = append(resultUser, conv.ConvertObjectToBytes(o.User)...)
	buf = append(buf, resultUser...)
	buf = append(buf, '}')
	return buf
}

func (o *GetUserRequest) Unmarshal(params map[string][]byte) error {
	_, arr, err := conv.GetArray(params["IdentityKey"])
	if err != nil {
		return err
	}
	for i, valByte := range arr {
		o.IdentityKey[i] = conv.ConvertBytesToByte(valByte)
	}
	o.Aa = conv.ConvertBytesToUint32(params["Aa"])
	o.Oo = conv.ConvertBytesToUint64(params["Oo"])
	o.User = &User{}
	conv.ConvertBytesToObject(o.User, params["User"])
	return nil
}

type User struct {
	Resp *GetUserResponse
}

func (o *User) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultResp []byte
	resultResp = append(resultResp, []byte("Resp:")...)
	resultResp = append(resultResp, conv.ConvertObjectToBytes(o.Resp)...)
	buf = append(buf, resultResp...)
	buf = append(buf, '}')
	return buf
}

func (o *User) Unmarshal(params map[string][]byte) error {
	o.Resp = &GetUserResponse{}
	conv.ConvertBytesToObject(o.Resp, params["Resp"])
	return nil
}

type GetUserResponse struct {
	SessionToken string
	IdentityKey  [32]byte
}

func (o *GetUserResponse) Marshal() []byte {
	var buf []byte
	buf = append(buf, '{')
	var resultSessionToken []byte
	resultSessionToken = append(resultSessionToken, []byte("SessionToken:")...)
	resultSessionToken = append(resultSessionToken, conv.ConvertStringToBytes(o.SessionToken)...)
	buf = append(buf, resultSessionToken...)
	buf = append(buf, ',')
	var resultIdentityKey []byte
	resultIdentityKey = append(resultIdentityKey, []byte("IdentityKey:")...)
	resultIdentityKey = append(resultIdentityKey, '[')
	for i, val := range o.IdentityKey {
		resultIdentityKey = append(resultIdentityKey, conv.ConvertByteToBytes(val)...)
		if i != len(o.IdentityKey)-1 {
			resultIdentityKey = append(resultIdentityKey, ',')
		}
	}
	resultIdentityKey = append(resultIdentityKey, ']')

	buf = append(buf, resultIdentityKey...)
	buf = append(buf, '}')
	return buf
}

func (o *GetUserResponse) Unmarshal(params map[string][]byte) error {
	o.SessionToken = conv.ConvertBytesToString(params["SessionToken"])
	_, arr, err := conv.GetArray(params["IdentityKey"])
	if err != nil {
		return err
	}
	for i, valByte := range arr {
		o.IdentityKey[i] = conv.ConvertBytesToByte(valByte)
	}
	return nil
}

type PeerGetUserArar struct {
	peer *peer.Peer
}

func (p *PeerGetUserArar) WriteResponse(msg *GetUserResponse) error {
	return p.peer.WriteResponse("GetUserArar", msg)
}

func (p *PeerGetUserArar) WriteError(err error) error {
	return p.peer.WriteError("GetUserArar", err)
}

func getPeerByHandlerName(handlerName string, peer *peer.Peer) interface{} {
	if handlerName == "GetUserArar" {
		return PeerGetUserArar{peer}
	}

	return nil
}

func initHandlers(userEndpoint UserEndpoint) map[string]*crypto_chateau.Handler {
	handlers := make(map[string]*crypto_chateau.Handler)

	handlers["GetUserAHAAHAHHAHAH"] = &crypto_chateau.Handler{
		CallFuncHandler: GetUserAHAAHAHHAHAHSqueeze(userEndpoint.GetUserAHAAHAHHAHAH),
		HandlerType:     crypto_chateau.HandlerT,
		RequestMsgType:  &GetUserRequest{},
	}

	handlers["GetUserArar"] = &crypto_chateau.Handler{
		CallFuncStream: GetUserArarSqueeze(userEndpoint.GetUserArar),
		HandlerType:    crypto_chateau.StreamT,
		RequestMsgType: &GetUserRequest{},
	}

	return handlers
}

func NewServer(cfg *crypto_chateau.Config, logger *zap.Logger, userEndpoint UserEndpoint) *crypto_chateau.Server {
	handlers := initHandlers(userEndpoint)

	return crypto_chateau.NewServer(cfg, logger, handlers)
}
