package main

import (
	"context"
	"time"

	zap2 "go.uber.org/zap"

	endpoints "github.com/oringik/crypto-chateau/examples/reverse/codegen"
	server2 "github.com/oringik/crypto-chateau/server"
)

func main() {
	zapInstance, _ := zap2.NewProduction()
	reverseEndpoint := &ReverseEndpoint{}

	connWriteDeadline := 500 * time.Second
	connReadDeadline := 500 * time.Second

	server := endpoints.NewServer(&server2.Config{
		IP:                "0.0.0.0",
		Port:              8080,
		ConnReadDeadline:  &connReadDeadline,
		ConnWriteDeadline: &connWriteDeadline,
	}, zapInstance, reverseEndpoint, nil)

	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type ReverseEndpoint struct{}

func (r *ReverseEndpoint) Rasd(ctx context.Context, req *endpoints.ReverseMagicStringRequest) (*endpoints.ReverseMagicStringResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReverseEndpoint) ReverseMagicString(ctx context.Context, req *endpoints.ReverseMagicStringRequest) (*endpoints.ReverseMagicStringResponse, error) {
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
		MagicInt8:           req.MagicInt8 + 100,
		MagicInt16:          req.MagicInt16 + 100,
		MagicInt32:          req.MagicInt32 + 100,
		MagicInt64:          req.MagicInt64 + 100,
		MagicUInt8:          req.MagicUInt8 + 100,
		MagicUInt16:         req.MagicUInt16 + 100,
		MagicUInt32:         req.MagicUInt32 + 100,
		MagicUInt64:         req.MagicUInt64 + 100,
		MagicBool:           !req.MagicBool,
		MagicBytes:          req.MagicBytes,
		MagicObject:         req.MagicObject,
		MagicObjectArray:    req.MagicObjectArray,
	}, nil
}
