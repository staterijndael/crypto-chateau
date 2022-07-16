package crypto_chateau

import (
	"context"
	"errors"
	"github.com/Oringik/crypto-chateau/generated"
	"github.com/Oringik/crypto-chateau/message"
)

type HandlerType int

var HandlerT HandlerType = 0
var StreamT HandlerType = 1

type Handler struct {
	callFunc interface{}
	HandlerType
	requestMsgType message.Message
}

func initHandlers(endpoint generated.Endpoint, handlers map[string]*Handler) {
	handlers["SendCode"] = &Handler{
		callFunc:       endpoint.UserEndpoint.SendCode,
		HandlerType:    HandlerT,
		requestMsgType: &generated.SendCodeRequest{},
	}
	handlers["GetEvents"] = &Handler{
		callFunc:    endpoint.UserEndpoint.GetEvents,
		HandlerType: StreamT,
	}
}

func callFuncToHandlerFunc(fnc interface{}) (func(context.Context, message.Message) (message.Message, error), error) {
	switch fnc.(type) {
	case func(ctx context.Context, request *generated.SendCodeRequest) (*generated.SendCodeResponse, error):
		callFunc := func(ctx context.Context, message message.Message) (message.Message, error) {
			convertedMessage, ok := message.(*generated.SendCodeRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.SendCodeRequest) (*generated.SendCodeResponse, error))(ctx, convertedMessage)
			return resp, err
		}

		return callFunc, nil
	default:
		return nil, errors.New("incorrect handler func type")
	}
}
