package crypto_chateau

import (
	"context"
	"crypto-chateau/generated"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockUserEndpoint struct {
}

func (m *mockUserEndpoint) InsertUser(ctx context.Context, msg Message) (Message, error) {
	fmt.Println("user inserted")

	return nil, nil
}

func Test_InitHandlers(t *testing.T) {
	s := &Server{
		Handlers: make(map[string]interface{}),
	}

	mock := &mockUserEndpoint{}

	err := s.initHandlers(generated.Endpoint{UserEndpoint: mock})
	assert.NoError(t, err)

	fmt.Println(s.Handlers)
}
