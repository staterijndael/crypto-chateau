package server

import (
	"context"
	"fmt"
	"github.com/oringik/crypto-chateau/message"
)

type mockUserEndpoint struct {
}

func (m *mockUserEndpoint) InsertUser(ctx context.Context, msg message.Message) (message.Message, error) {
	fmt.Println("user inserted")

	return nil, nil
}
