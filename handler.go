package crypto_chateau

import (
	"context"
	"errors"
	"fmt"
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

func callFuncToHandlerFunc(fnc interface{}) (func(context.Context, Message) (Message, error), error) {
	switch fnc.(type) {
	case func(context.Context, *generated.GetUserRequest) (*generated.GetUserResponse, error):
		callFunc := func(ctx context.Context, message Message) (Message, error) {
			convertedMessage, ok := message.(*generated.GetUserRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				fmt.Println(err)
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.GetUserRequest) (*generated.GetUserResponse, error))(ctx, convertedMessage)
			return resp, err
		}

		return callFunc, nil
	default:
		return nil, errors.New("incorrect handler func type")
	}
}
