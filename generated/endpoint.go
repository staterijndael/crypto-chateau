package generated

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oringik/crypto-chateau/message"
	"strconv"
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
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	GetUserUpdates(context.Context, StreamI) error
}

type SendCodeRequest struct {
	Number   string
	PassHash string
}

type SendCodeResponse struct {
}

type GetUserRequest struct {
	UserID uint64
}

type User struct {
	Id       uint64
	Nickname string
	Age      int
	Gender   bool
	Status   string
}

type GetUserResponse struct {
	User User
}

type GetUsersRequest struct {
	Offset int
	Limit  int
}

type GetUsersResponse struct {
	Users []*User
}

type CreateUserRequest struct {
	User *User
}

type CreateUserResponse struct {
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

func (i *GetUserRequest) Marshal() []byte {
	return []byte("GetUser# UserID:" + strconv.Itoa(int(i.UserID)))
}

func (i *GetUserResponse) Marshal() []byte {
	strBytes := []byte("GetUser# User:")
	strBytes = append(strBytes, i.User.Marshal()...)
	return strBytes
}

func (i *GetUsersRequest) Marshal() []byte {
	marshalStr := fmt.Sprintf("GetUsersRequest# Offset: %d, Limit: %d", i.Offset, i.Limit)
	return []byte(marshalStr)
}

func (i *GetUsersResponse) Marshal() []byte {
	var usersParam []byte
	for j := 0; j < len(i.Users); j++ {
		usersParam = append(usersParam, i.Users[j].Marshal()...)
		if j < len(i.Users)-1 {
			usersParam = append(usersParam, byte(','))
		}
	}
	marshalStr := fmt.Sprintf("GetUsersResponse# Users: {%s}", string(usersParam))
	return []byte(marshalStr)
}

func (i *CreateUserRequest) Marshal() []byte {
	marshalStr := fmt.Sprintf("GetUser# User: %s", string(i.User.Marshal()))
	return []byte(marshalStr)
}

func (i *CreateUserResponse) Marshal() []byte {
	return nil
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

func (i *GetUserRequest) Unmarshal(params map[string][]byte) error {
	i.UserID = binary.BigEndian.Uint64(params["UserID"])

	return nil
}

func (i *GetUserResponse) Unmarshal(params map[string][]byte) error {
	return i.User.Unmarshal(params)
}

func (i *GetUsersRequest) Unmarshal(params map[string][]byte) error {
	i.Offset = int(binary.BigEndian.Uint64(params["Offset"]))
	i.Limit = int(binary.BigEndian.Uint64(params["Limit"]))

	return nil
}

func (i *GetUsersResponse) Unmarshal(params map[string][]byte) error {
	rawUsers, err := message.ParseArray(params["Users"])
	if err != nil {
		return err
	}

	users := make([]*User, 0, len(rawUsers))

	for _, rawUser := range rawUsers {
		userParams, err := message.GetParams(rawUser)
		if err != nil {
			return err
		}

		user := &User{}
		err = user.Unmarshal(userParams)
		if err != nil {
			return err
		}

		users = append(users, user)
	}

	i.Users = users

	return nil
}

func (i *CreateUserRequest) Unmarshal(params map[string][]byte) error {
	return i.User.Unmarshal(params)
}

func (i *CreateUserResponse) Unmarshal(params map[string][]byte) error {
	return nil
}
