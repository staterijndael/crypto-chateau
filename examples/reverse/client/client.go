package main

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	endpoints "github.com/oringik/crypto-chateau/examples/reverse/codegen"
)

func main() {
	client, err := endpoints.NewClientReverse("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100000; i++ {

		resp, err := client.ReverseMagicString(context.Background(), &endpoints.ReverseMagicStringRequest{
			MagicString: "privet kotik",
			MagicInt8:   10,
			MagicInt16:  20,
			MagicInt32:  30,
			MagicInt64:  40,
			MagicUInt8:  50,
			MagicUInt16: 60,
			MagicUInt32: 70,
			MagicUInt64: 80,
			MagicBool:   true,
			MagicBytes:  []byte{1, 2, 3, 4, 5},
			MagicObject: endpoints.ReverseCommonObject{
				Key:   [16]byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115},
				Value: [32]string{"hello", "world"},
			},
			MagicObjectArray: []endpoints.ReverseCommonObject{
				{
					Key:   [16]byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115},
					Value: [32]string{"sub", "object"},
				},
			},
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.ReversedMagicString + " " + strconv.Itoa(i+1))

		excepted := &endpoints.ReverseMagicStringResponse{
			ReversedMagicString: "kitok tevirp",
			MagicInt8:           110,
			MagicInt16:          120,
			MagicInt32:          130,
			MagicInt64:          140,
			MagicUInt8:          150,
			MagicUInt16:         160,
			MagicUInt32:         170,
			MagicUInt64:         180,
			MagicBool:           false,
			MagicBytes:          []byte{1, 2, 3, 4, 5},
			MagicObject: endpoints.ReverseCommonObject{
				Key:   [16]byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115},
				Value: [32]string{"hello", "world"},
			},
			MagicObjectArray: []endpoints.ReverseCommonObject{
				{
					Key:   [16]byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115},
					Value: [32]string{"sub", "object"},
				},
			},
		}

		ok := reflect.DeepEqual(resp, excepted)

		if !ok {
			fmt.Printf("expected:\t%+v\ngot:\t\t%+v\n", excepted, resp)
			panic("not equal")
		}
	}
}
