package generated

import (
	"context"
)

type Endpoint struct {
	UserEndpoint UserEndpoint
}

type UserEndpoint interface {
	InsertUser(context.Context, *InsertUserRequest) (*InsertUserResponse, error)
}

type InsertUserRequest struct {
	UserID uint64
}

type InsertUserResponse struct {
}

func (i *InsertUserRequest) Marshal() ([]byte, error) {
	return nil, nil
}

func (i *InsertUserResponse) Marshal() ([]byte, error) {
	return nil, nil
}
