package main

import (
	"context"
	"fmt"
	endpoints "github.com/oringik/crypto-chateau/examples/ping-pong-stream/codegen"
)

func main() {
	client, err := endpoints.NewClientReverse("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}

	count := 8
	peer, err := client.PingPong(context.Background(), &endpoints.PingPongRequest{})
	if err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		msg := &endpoints.Ping{Req: "ping"}
		err := peer.WriteResponse(msg)
		if err != nil {
			panic(err)
		}

		gotMessage := &endpoints.Pong{}

		err = peer.ReadMessage(gotMessage)
		if err != nil {
			panic(err)
		}

		if gotMessage.Resp != "pong" {
			panic(gotMessage.Resp)
		}

		fmt.Println(gotMessage.Resp)
	}

	peer.Close()
}
