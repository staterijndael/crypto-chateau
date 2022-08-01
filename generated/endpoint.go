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
	HandleCode(context.Context, *HandleCodeRequest) (*HandleCodeResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	AuthToken(context.Context, *AuthTokenRequest) (*AuthTokenResponse, error)
	AuthCredentials(context.Context, *AuthCredentialsRequest) (*AuthCredentialsResponse, error)
	GetEvents(context.Context, StreamI) error
}

type AuthTokenRequest struct {
	SessionToken string
}

type AuthTokenResponse struct {
}

type AuthCredentialsRequest struct {
	Number   string
	PassHash string
}

type AuthCredentialsResponse struct {
	SessionToken string
}

type RegisterRequest struct {
	Number   string
	Code     string
	Nickname string
	Bio      string
	PassHash string
}

type RegisterResponse struct {
	SessionToken string
}

type SendCodeRequest struct {
	Number string
}

type SendCodeResponse struct {
}

type HandleCodeRequest struct {
	Number string
	Code   uint8
}

type HandleCodeResponse struct {
}

func (a AuthCredentialsResponse) Marshal() []byte {
	return []byte(fmt.Sprintf("AuthCreds# SessionToken:%s", a.SessionToken))
}

func (i *AuthCredentialsRequest) Marshal() []byte {
	return nil
}

func (i *AuthTokenRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("AuthToken# SessionToken:%s", i.SessionToken))
}

func (i *AuthTokenResponse) Marshal() []byte {
	return nil
}

func (i *RegisterRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("Register# Nickname:%s,Bio:%s", i.Nickname, i.Bio))
}

func (i *RegisterResponse) Marshal() []byte {
	return nil
}

func (i *HandleCodeRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("HandleCode# Code: %d, Number: %s", i.Code, i.Number))
}

func (i *HandleCodeResponse) Marshal() []byte {
	return nil
}

func (i *SendCodeRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("SendCode# Number:%s", i.Number))
}

func (i *SendCodeResponse) Marshal() []byte {
	return []byte("SendCode#")
}

// unmarshal

func (i *RegisterRequest) Unmarshal(params map[string][]byte) error {
	i.Bio = string(params["Bio"])
	i.Nickname = string(params["Nickname"])
	i.Number = string(params["Number"])
	i.Code = string(params["Code"])
	i.PassHash = string(params["PassHash"])

	return nil
}

func (i *RegisterResponse) Unmarshal(params map[string][]byte) error {
	i.SessionToken = string(params["SessionToken"])

	return nil
}

func (i *SendCodeRequest) Unmarshal(params map[string][]byte) error {
	if len(params["PassHash"]) == 0 || len(params["Number"]) == 0 {
		return errors.New("incorrect number or pass hash")
	}

	i.Number = string(params["Number"])

	return nil
}

func (i *AuthCredentialsRequest) Unmarshal(params map[string][]byte) error {
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

func (i *AuthTokenRequest) Unmarshal(params map[string][]byte) error {
	i.SessionToken = string(params["SessionToken"])

	return nil
}

func (i *AuthTokenResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

func (i *HandleCodeRequest) Unmarshal(params map[string][]byte) error {
	i.Code = uint8(binary.BigEndian.Uint16(params["Code"]))
	i.Number = string(params["Number"])

	return nil
}

func (i *HandleCodeResponse) Unmarshal(m map[string][]byte) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthCredentialsResponse) Unmarshal(m map[string][]byte) error {
	//TODO implement me
	panic("implement me")
}
