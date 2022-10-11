package endpoints

import (
	"fmt"
	"github.com/Oringik/crypto-chateau/gen/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenDeifinition(t *testing.T) {
	getUserReq := &GetUserRequest{
		IdentityKey: [32]byte{'a', 'a', 'b', 'd'},
		Aa:          12,
		Oo:          32,
		User: &User{
			Resp: &GetUserResponse{
				SessionToken: "asdasd",
				IdentityKey:  [32]byte{'a', 'b', 'u', 'f'},
			},
		},
	}

	result := getUserReq.Marshal()
	fmt.Println(string(result))

	getUserReqUn := &GetUserRequest{}
	_, params, err := conv.GetParams(result)
	assert.NoError(t, err)
	err = getUserReqUn.Unmarshal(params)
	assert.NoError(t, err)

	fmt.Println(getUserReqUn)
}
