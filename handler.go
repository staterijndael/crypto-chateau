package crypto_chateau

import (
	"github.com/Oringik/crypto-chateau/generated"
)

type Handler struct {
	callFunc       interface{}
	requestMsgType Message
}

func initHandlers(endpoint generated.Endpoint, handlers map[string]*Handler) {
	handlers["GetUser"] = &Handler{
		callFunc:       endpoint.UserEndpoint.GetUser,
		requestMsgType: &generated.GetUserRequest{},
	}
}
