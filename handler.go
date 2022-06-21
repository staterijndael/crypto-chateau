package crypto_chateau

import (
	"github.com/Oringik/crypto-chateau/generated"
)

type HandlerType int

var HandlerT HandlerType = 0
var StreamT HandlerType = 1

type Handler struct {
	callFunc interface{}
	HandlerType
	requestMsgType Message
}

func initHandlers(endpoint generated.Endpoint, handlers map[string]*Handler) {
	handlers["GetUser"] = &Handler{
		callFunc:       endpoint.UserEndpoint.GetUser,
		HandlerType:    HandlerT,
		requestMsgType: &generated.GetUserRequest{},
	}
}
