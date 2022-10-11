package generated

import (
	"context"
	"github.com/Oringik/crypto-chateau/message"
	"github.com/Oringik/crypto-chateau/peer"
)

type HandlerType int

var HandlerT HandlerType = 0
var StreamT HandlerType = 1

type Handler struct {
	CallFuncHandler HandlerFunc
	CallFuncStream  StreamFunc
	HandlerType
	RequestMsgType message.Message
}

type StreamReq interface {
	Init(ctx context.Context, peer *peer.Peer, initMessage message.Message) error
}

func InitHandlers(endpoint Endpoint, handlers map[string]*Handler) {
	handlers["SendCode"] = &Handler{
		CallFuncHandler: SendCodeSqueeze(endpoint.UserEndpoint.SendCode),
		HandlerType:     HandlerT,
		RequestMsgType:  &SendCodeRequest{},
	}
	handlers["HandleCode"] = &Handler{
		CallFuncHandler: HandleCodeSqueeze(endpoint.UserEndpoint.HandleCode),
		HandlerType:     HandlerT,
		RequestMsgType:  &HandleCodeRequest{},
	}
	handlers["Register"] = &Handler{
		CallFuncHandler: RegisterSqueeze(endpoint.UserEndpoint.Register),
		HandlerType:     HandlerT,
		RequestMsgType:  &RegisterRequest{},
	}
	handlers["AuthToken"] = &Handler{
		CallFuncHandler: AuthTokenSqueeze(endpoint.UserEndpoint.AuthToken),
		HandlerType:     HandlerT,
		RequestMsgType:  &AuthTokenRequest{},
	}
	handlers["AuthCreds"] = &Handler{
		CallFuncHandler: AuthCredentialsSqueeze(endpoint.UserEndpoint.AuthCredentials),
		HandlerType:     HandlerT,
		RequestMsgType:  &AuthCredentialsRequest{},
	}
	handlers["RequiredOPKSqueeze"] = &Handler{
		CallFuncHandler: RequiredOPKSqueeze(endpoint.UserEndpoint.RequiredOPK),
		HandlerType:     HandlerT,
		RequestMsgType:  &RequiredOPKRequest{},
	}
	handlers["FindUsersByPartNickname"] = &Handler{
		CallFuncHandler: FindUsersByPartNicknameSqueeze(endpoint.UserEndpoint.FindUsersByPartNickname),
		HandlerType:     HandlerT,
		RequestMsgType:  &FindUsersByPartNicknameRequest{},
	}
	handlers["GetInitMsgKeys"] = &Handler{
		CallFuncHandler: GetInitMsgKeysSqueeze(endpoint.UserEndpoint.GetInitMsgKeys),
		HandlerType:     HandlerT,
		RequestMsgType:  &GetInitMsgKeysRequest{},
	}
	handlers["LoadOPK"] = &Handler{
		CallFuncHandler: LoadOPKSqueeze(endpoint.UserEndpoint.LoadOPK),
		HandlerType:     HandlerT,
		RequestMsgType:  &LoadOPKRequest{},
	}
	handlers["GetEvents"] = &Handler{
		CallFuncStream: GetEventsSqueeze(endpoint.UserEndpoint.GetEvents),
		HandlerType:    StreamT,
		RequestMsgType: &EventStreamInitMessage{},
	}

}
