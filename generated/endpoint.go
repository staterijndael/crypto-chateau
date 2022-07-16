package generated

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
)

type Endpoint struct {
	UserEndpoint UserEndpoint
}

type StreamI interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type UserEndpoint interface {
	SendCode(context.Context, *SendCodeRequest) (*SendCodeResponse, error)
	GetEvents(context.Context, StreamI) error
}

type SendCodeRequest struct {
	Number   string
	PassHash string
}

type SendCodeResponse struct {
}

type User struct {
	Id       uint64
	Nickname string
	Age      int
	Gender   bool
	Status   string
}

func (i *SendCodeRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("SendCode# Number:%s,PassHash:%s", i.Number, i.PassHash))
}

func (i *SendCodeResponse) Marshal() []byte {
	return []byte("SendCode#")
}

func (i *User) Marshal() []byte {
	idBytes, ageBytes := make([]byte, 8), make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, i.Id)
	binary.BigEndian.PutUint64(ageBytes, uint64(i.Age))
	var gender byte
	if i.Gender {
		gender = 1
	} else {
		gender = 0
	}

	marshalBytes := []byte("(Id:")
	marshalBytes = append(marshalBytes, idBytes...)
	marshalBytes = append(marshalBytes, []byte(",Nickname:"+i.Nickname+",Age:")...)
	marshalBytes = append(marshalBytes, ageBytes...)
	marshalBytes = append(marshalBytes, []byte(",Gender:")...)
	marshalBytes = append(marshalBytes, gender)
	marshalBytes = append(marshalBytes, []byte(",Status:")...)
	marshalBytes = append(marshalBytes, []byte(i.Status)...)
	marshalBytes = append(marshalBytes, []byte(")")...)

	return marshalBytes
}

// unmarshal

func (i *SendCodeRequest) Unmarshal(params map[string][]byte) error {
	if len(params["PassHash"]) == 0 || len(params["Number"]) == 0 {
		return errors.New("incorrect number or pass hash")
	}

	i.PassHash = string(params["PassHash"])
	i.Number = string(params["Number"])

	return nil
}

func (i *SendCodeResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

func (i *User) Unmarshal(params map[string][]byte) error {
	i.Id = binary.BigEndian.Uint64(params["Id"])
	i.Age = int(binary.BigEndian.Uint64(params["Age"]))
	if params["Gender"][0] == '1' {
		i.Gender = true
	} else {
		i.Gender = false
	}
	i.Status = string(params["Status"])

	return nil
}
