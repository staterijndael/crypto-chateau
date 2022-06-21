package generated

import (
	"context"
	"strconv"
)

type Endpoint struct {
	UserEndpoint UserEndpoint
}

type GetUserFunc func(context.Context, *GetUserRequest) (*GetUserResponse, error)

type UserEndpoint interface {
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
}

type GetUserRequest struct {
	UserID uint64
}

type GetUserResponse struct {
	UserName string
}

func (i *GetUserRequest) Marshal() ([]byte, error) {
	return []byte("GetUser# UserID:" + strconv.Itoa(int(i.UserID))), nil
}

func (i *GetUserResponse) Marshal() ([]byte, error) {
	return []byte("GetUser# UserName:" + i.UserName), nil
}
