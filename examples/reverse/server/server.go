package main

import (
	"context"
	endpoints "github.com/oringik/crypto-chateau/examples/reverse/codegen"
	server2 "github.com/oringik/crypto-chateau/server"
	zap2 "go.uber.org/zap"
	"log"
	"time"
)

func main() {
	log.Println("server started")

	zapInstance, _ := zap2.NewProduction()
	reverseEndpoint := &ReverseEndpoint{}

	connWriteDeadline := 500 * time.Second
	connReadDeadline := 500 * time.Second

	server := endpoints.NewServer(&server2.Config{
		IP:                "0.0.0.0",
		Port:              8080,
		ConnReadDeadline:  &connReadDeadline,
		ConnWriteDeadline: &connWriteDeadline,
	}, zapInstance, reverseEndpoint)

	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type ReverseEndpoint struct{}

func (r *ReverseEndpoint) ReverseMagicString(ctx context.Context, req *endpoints.ReverseMagicStringRequest) (*endpoints.ReverseMagicStringResponse, error) {
	log.Printf("recived %s", req.MagicString)
	magicStringRunes := []rune(req.MagicString)
	left := 0
	right := len(magicStringRunes) - 1

	for left < right {
		magicStringRunes[left] = magicStringRunes[left] ^ magicStringRunes[right]
		magicStringRunes[right] = magicStringRunes[left] ^ magicStringRunes[right]
		magicStringRunes[left] = magicStringRunes[left] ^ magicStringRunes[right]

		left++
		right--
	}

	reversedMsg := string(magicStringRunes)

	return &endpoints.ReverseMagicStringResponse{
		ReversedMagicString: reversedMsg,
	}, nil
}
