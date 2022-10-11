package main

import (
	"context"
	"fmt"
	"github.com/Oringik/crypto-chateau/message"
)

type mockUserEndpoint struct {
}

func (m *mockUserEndpoint) InsertUser(ctx context.Context, msg message.Message) (message.Message, error) {
	fmt.Println("user inserted")

	return nil, nil
}
