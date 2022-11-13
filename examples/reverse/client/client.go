package main

import (
	"context"
	"fmt"
	endpoints "github.com/oringik/crypto-chateau/examples/reverse/codegen"
)

func main() {
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

	fmt.Println(resp.ReversedMagicString)
}
