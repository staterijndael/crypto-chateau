package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/peer"
	"time"

	zap2 "go.uber.org/zap"

	endpoints "github.com/oringik/crypto-chateau/examples/ping-pong-stream/codegen"
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
	}, zapInstance, reverseEndpoint)

	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type ReverseEndpoint struct{}

func (r *ReverseEndpoint) PingPong(ctx context.Context, peer *peer.Peer, req *endpoints.PingPongRequest) error {
	for {
		msg := &endpoints.Ping{}
		err := peer.ReadMessage(msg)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		if msg.Req != "ping" {
			return errors.New("expected ping message")
		}

		respMsg := &endpoints.Pong{Resp: "pong"}

		err = peer.WriteResponse(respMsg)
		if err != nil {
			break
		}
	}

	return nil
}
