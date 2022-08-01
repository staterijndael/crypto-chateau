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
	handlers["HandleCode"] = &Handler{
		callFunc:       endpoint.UserEndpoint.HandleCode,
		HandlerType:    HandlerT,
		requestMsgType: &generated.HandleCodeRequest{},
	}
	handlers["Register"] = &Handler{
		callFunc:       endpoint.UserEndpoint.Register,
		HandlerType:    HandlerT,
		requestMsgType: &generated.RegisterRequest{},
	}
	handlers["AuthToken"] = &Handler{
		callFunc:       endpoint.UserEndpoint.AuthToken,
		HandlerType:    HandlerT,
		requestMsgType: &generated.AuthTokenRequest{},
	}
	handlers["AuthCredentials"] = &Handler{
		callFunc:       endpoint.UserEndpoint.AuthCredentials,
		HandlerType:    HandlerT,
		requestMsgType: &generated.AuthCredentialsRequest{},
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
	case func(ctx context.Context, request *generated.HandleCodeRequest) (*generated.HandleCodeResponse, error):
		callFunc := func(ctx context.Context, message message.Message) (message.Message, error) {
			convertedMessage, ok := message.(*generated.HandleCodeRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.HandleCodeRequest) (*generated.HandleCodeResponse, error))(ctx, convertedMessage)
			return resp, err
		}
		return callFunc, nil
	case func(ctx context.Context, request *generated.AuthTokenRequest) (*generated.AuthTokenResponse, error):
		callFunc := func(ctx context.Context, message message.Message) (message.Message, error) {
			convertedMessage, ok := message.(*generated.AuthTokenRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.AuthTokenRequest) (*generated.AuthTokenResponse, error))(ctx, convertedMessage)
			return resp, err
		}
		return callFunc, nil
	case func(ctx context.Context, request *generated.RegisterRequest) (*generated.RegisterResponse, error):
		callFunc := func(ctx context.Context, message message.Message) (message.Message, error) {
			convertedMessage, ok := message.(*generated.RegisterRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.RegisterRequest) (*generated.RegisterResponse, error))(ctx, convertedMessage)
			return resp, err
		}
		return callFunc, nil
	case func(ctx context.Context, request *generated.AuthCredentialsRequest) (*generated.AuthCredentialsResponse, error):
		callFunc := func(ctx context.Context, message message.Message) (message.Message, error) {
			convertedMessage, ok := message.(*generated.AuthCredentialsRequest)
			if !ok {
				err := errors.New("error converting message to GetUserRequest")
				return nil, err
			}

			resp, err := fnc.(func(context.Context, *generated.AuthCredentialsRequest) (*generated.AuthCredentialsResponse, error))(ctx, convertedMessage)
			return resp, err
		}
		return callFunc, nil
	default:
		return nil, errors.New("incorrect handler func type")
	}
}
