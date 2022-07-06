package generated

import (
	"context"
	"fmt"
	"strconv"
)

type Endpoint struct {
	UserEndpoint UserEndpoint
}

type StreamI interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type UserEndpoint interface {
	SendCode(context.Context, *SendCodeRequest) (*SendCodeResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	GetUserUpdates(context.Context, StreamI) error
}

type SendCodeRequest struct {
	Number   string
	PassHash string
}

type SendCodeResponse struct {
}

type GetUserRequest struct {
	UserID uint64
}

type GetUserResponse struct {
	UserName string
}

func (i *SendCodeRequest) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("SendCode# Number:%s,PassHash:%s", i.Number, i.PassHash)), nil
}

func (i *SendCodeResponse) Marshal() ([]byte, error) {
	return nil, nil
}

func (i *GetUserRequest) Marshal() ([]byte, error) {
	return []byte("GetUser# UserID:" + strconv.Itoa(int(i.UserID))), nil
}

func (i *GetUserResponse) Marshal() ([]byte, error) {
	return []byte("GetUser# UserName:\"" + i.UserName + "\""), nil
}
