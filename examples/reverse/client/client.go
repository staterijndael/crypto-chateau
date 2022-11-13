package main

import (
	"context"
	endpoints "github.com/oringik/crypto-chateau/examples/reverse/codegen"
	"log"
	"time"
)

func main() {
	t := time.Now()

	client, err := endpoints.NewClientReverse("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}

	resp, err := client.ReverseMagicString(context.Background(), &endpoints.ReverseMagicStringRequest{
		MagicString: "privet kotik",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("time to send - recive is %s", time.Since(t))

	log.Printf("recived response is: %s", resp.ReversedMagicString)
}
