package generated

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oringik/crypto-chateau/message"
	crypto_chateau "github.com/Oringik/crypto-chateau/peer"
)

type Endpoint struct {
	UserEndpoint UserEndpoint
}

type HandlerFunc func(context.Context, message.Message) (message.Message, error)
type StreamFunc func(ctx context.Context, req StreamReq) error

type UserEndpoint interface {
	SendCode(context.Context, *SendCodeRequest) (*SendCodeResponse, error)
	HandleCode(context.Context, *HandleCodeRequest) (*HandleCodeResponse, error)
	RequiredOPK(context.Context, *RequiredOPKRequest) (*RequiredOPKResponse, error)
	LoadOPK(context.Context, *LoadOPKRequest) (*LoadOPKResponse, error)
	FindUsersByPartNickname(context.Context, *FindUsersByPartNicknameRequest) (*FindUsersByPartNicknameResponse, error)
	GetInitMsgKeys(context.Context, *GetInitMsgKeysRequest) (*GetInitMsgKeysResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	AuthToken(context.Context, *AuthTokenRequest) (*AuthTokenResponse, error)
	AuthCredentials(context.Context, *AuthCredentialsRequest) (*AuthCredentialsResponse, error)
	GetEvents(context.Context, *EventStream) error
}

// squeezes

func SendCodeSqueeze(fnc func(context.Context, *SendCodeRequest) (*SendCodeResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*SendCodeRequest); ok {
			return fnc(ctx, msg.(*SendCodeRequest))
		} else {
			return nil, errors.New("unknown message type: expected SendCodeRequest")
		}
	}
}

func HandleCodeSqueeze(fnc func(context.Context, *HandleCodeRequest) (*HandleCodeResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*HandleCodeRequest); ok {
			return fnc(ctx, msg.(*HandleCodeRequest))
		} else {
			return nil, errors.New("unknown message type: expected HandleCodeRequest")
		}
	}
}

func RequiredOPKSqueeze(fnc func(context.Context, *RequiredOPKRequest) (*RequiredOPKResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*RequiredOPKRequest); ok {
			return fnc(ctx, msg.(*RequiredOPKRequest))
		} else {
			return nil, errors.New("unknown message type: expected RequiredOPKRequest")
		}
	}
}

func FindUsersByPartNicknameSqueeze(fnc func(context.Context, *FindUsersByPartNicknameRequest) (*FindUsersByPartNicknameResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*FindUsersByPartNicknameRequest); ok {
			return fnc(ctx, msg.(*FindUsersByPartNicknameRequest))
		} else {
			return nil, errors.New("unknown message type: expected FindUsersByPartNickname")
		}
	}
}

func GetInitMsgKeysSqueeze(fnc func(context.Context, *GetInitMsgKeysRequest) (*GetInitMsgKeysResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*GetInitMsgKeysRequest); ok {
			return fnc(ctx, msg.(*GetInitMsgKeysRequest))
		} else {
			return nil, errors.New("unknown message type: expected GetInitMsgKeysRequest")
		}
	}
}

func LoadOPKSqueeze(fnc func(context.Context, *LoadOPKRequest) (*LoadOPKResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*LoadOPKRequest); ok {
			return fnc(ctx, msg.(*LoadOPKRequest))
		} else {
			return nil, errors.New("unknown message type: expected LoadOPKRequest")
		}
	}
}

func RegisterSqueeze(fnc func(context.Context, *RegisterRequest) (*RegisterResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*HandleCodeRequest); ok {
			return fnc(ctx, msg.(*RegisterRequest))
		} else {
			return nil, errors.New("unknown message type: expected RegisterRequest")
		}
	}
}

func AuthTokenSqueeze(fnc func(context.Context, *AuthTokenRequest) (*AuthTokenResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*HandleCodeRequest); ok {
			return fnc(ctx, msg.(*AuthTokenRequest))
		} else {
			return nil, errors.New("unknown message type: expected AuthTokenRequest")
		}
	}
}

func AuthCredentialsSqueeze(fnc func(context.Context, *AuthCredentialsRequest) (*AuthCredentialsResponse, error)) HandlerFunc {
	return func(ctx context.Context, msg message.Message) (message.Message, error) {
		if _, ok := msg.(*AuthCredentialsRequest); ok {
			return fnc(ctx, msg.(*AuthCredentialsRequest))
		} else {
			return nil, errors.New("unknown message type: expected AuthCredentialsRequest")
		}
	}
}

func GetEventsSqueeze(fnc func(context.Context, *EventStream) error) StreamFunc {
	return func(ctx context.Context, req StreamReq) error {
		if _, ok := req.(*EventStream); ok {
			return fnc(ctx, req.(*EventStream))
		} else {
			return errors.New("unknown message type: expected EventStream")
		}
	}
}

type EventStream struct {
	Peer        *crypto_chateau.Peer
	InitMessage *EventStreamInitMessage
}

func (e *EventStream) Init(peer *crypto_chateau.Peer, initMessage message.Message) error {
	e.Peer = peer
	if _, ok := initMessage.(*EventStreamInitMessage); ok {
		e.InitMessage = initMessage.(*EventStreamInitMessage)
	} else {
		return errors.New("unknown type of init message")
	}

	return nil
}

type EventStreamInitMessage struct {
	SessionToken string
	LastEventID  uint64
}

func (e *EventStreamInitMessage) Marshal() []byte {
	return nil
}

func (r *GetInitMsgKeysResponse) Marshal() []byte {
	var buf []byte
	buf = append(buf, []byte("GetInitMsgKeys# OPK:")...)
	buf = append(buf, r.OPK[:]...)
	buf = append(buf, ",SignedLTPK:"...)
	buf = append(buf, r.SignedLTPK[:]...)
	buf = append(buf, ",Signature:"...)
	buf = append(buf, r.Signature[:]...)
	buf = append(buf, ",OPKId:"...)
	var bufOpkID []byte
	// this method mutate first 4 bytes of buffer
	binary.BigEndian.PutUint32(bufOpkID, r.OPKId)
	buf = append(buf, bufOpkID...)

	return buf
}

func (e *EventStreamInitMessage) Unmarshal(params map[string][]byte) error {
	e.SessionToken = string(params["SessionToken"])
	e.LastEventID = binary.BigEndian.Uint64(params["LastEventID"])

	return nil
}

type Event struct {
	Type string
	Info message.Message
}

func (e *Event) Marshal() []byte {
	return []byte(fmt.Sprintf("Type:%v,Info:{%v}", e.Type, e.Info.Marshal()))
}

func (e *EventStream) Write(event Event) error {
	msg := []byte(fmt.Sprintf("GetEvents# Event: {%v}", event.Marshal()))
	n, err := e.Peer.Write(msg)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("0 bytes written")
	}

	return nil
}

type FindUsersByPartNicknameRequest struct {
	PartNickname string
}

func (f *FindUsersByPartNicknameRequest) Marshal() []byte {
	return nil
}

func (f *FindUsersByPartNicknameRequest) Unmarshal(params map[string][]byte) error {
	f.PartNickname = string(params["PartNickname"])

	return nil
}

type PresentUser struct {
	IdentityKey [32]byte
	Nickname    string
	PictureID   string
	Status      string
}

func (p *PresentUser) Marshal() []byte {
	var buf []byte
	buf = append(buf, []byte("FindUserByPartNickname# IdentityKey:")...)
	buf = append(buf, p.IdentityKey[:]...)
	buf = append(buf, fmt.Sprintf(",Nickname%s,PictureID:%s,Status:%s", p.Nickname, p.PictureID, p.Status)...)

	return buf
}

type FindUsersByPartNicknameResponse struct {
	Users []*PresentUser
}

func (f *FindUsersByPartNicknameResponse) Marshal() []byte {
	var usersStr string
	for i, user := range f.Users {
		usersStr += "{"
		usersStr += string(user.Marshal())
		usersStr += "}"
		if i < len(f.Users)-1 {
			usersStr += ","
		}
	}

	return []byte(fmt.Sprintf("FindUsersByPartNickname# {%s}", usersStr))
}

func (f *FindUsersByPartNicknameResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

type GetInitMsgKeysRequest struct {
	IdentityKey [32]byte
}

type GetInitMsgKeysResponse struct {
	OPKId      uint32
	OPK        [32]byte
	SignedLTPK [32]byte
	Signature  [64]byte
}

type OPKPair struct {
	OPKId uint32
	OPK   [32]byte
}

func (o *OPKPair) Marshal() []byte {
	return nil
}

func (o *OPKPair) Unmarshal(params map[string][]byte) error {
	o.OPKId = binary.BigEndian.Uint32(params["OPKId"])
	copy(o.OPK[:], params["OPK"])

	return nil
}

type LoadOPKRequest struct {
	SessionToken [16]byte
	OPK          []OPKPair
}

type LoadOPKResponse struct{}

type RequiredOPKRequest struct {
	SessionToken [16]byte
}

type RequiredOPKResponse struct {
	Count int
}

type AuthTokenRequest struct {
	SessionToken [16]byte
}

type AuthTokenResponse struct {
	SessionToken [16]byte
}

type AuthCredentialsRequest struct {
	Number   string
	PassHash string
}

type AuthCredentialsResponse struct {
	SessionToken [16]byte
}

type RegisterRequest struct {
	Number   string
	Code     uint8
	Nickname string
	PassHash string

	DeviceID    string
	IdentityKey [32]byte
}

type RegisterResponse struct {
	SessionToken [16]byte
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

func (r *GetInitMsgKeysRequest) Marshal() []byte {
	return nil
}

func (r *RequiredOPKResponse) Marshal() []byte {
	return []byte(fmt.Sprintf("RequiredOPK# Count: %d", r.Count))
}

func (a *AuthCredentialsResponse) Marshal() []byte {
	return []byte(fmt.Sprintf("AuthCreds# SessionToken:%s", a.SessionToken))
}

func (l *LoadOPKRequest) Marshal() []byte {
	return nil
}

func (r *RequiredOPKRequest) Marshal() []byte {
	return nil
}

func (i *AuthCredentialsRequest) Marshal() []byte {
	return nil
}

func (l *LoadOPKResponse) Marshal() []byte {
	return []byte("LoadOPK#")
}

func (i *AuthTokenRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("AuthToken# SessionToken:%s", i.SessionToken))
}

func (i *AuthTokenResponse) Marshal() []byte {
	return []byte(fmt.Sprintf("AuthToken# SessionToken:%s", i.SessionToken))
}

func (i *RegisterRequest) Marshal() []byte {
	var buf []byte
	buf = append(buf, fmt.Sprintf("Register# Nickname:%s,Number:%s,Code:%d,PassHash:%s,DeviceID:%s,IdentityKey:", i.Nickname, i.Number, i.Code, i.PassHash, i.DeviceID)...)
	buf = append(buf, i.IdentityKey[:]...)

	return buf
}

func (i *RegisterResponse) Marshal() []byte {
	return []byte("Register#")
}

func (i *HandleCodeRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("HandleCode# Code: %d, Number: %s", i.Code, i.Number))
}

func (i *HandleCodeResponse) Marshal() []byte {
	return []byte("HandleCode#")
}

func (i *SendCodeRequest) Marshal() []byte {
	return []byte(fmt.Sprintf("SendCode# Number:%s", i.Number))
}

func (i *SendCodeResponse) Marshal() []byte {
	return []byte("SendCode#")
}

// unmarshal

func (r *GetInitMsgKeysResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

func (r *GetInitMsgKeysRequest) Unmarshal(params map[string][]byte) error {
	copy(r.IdentityKey[:], params["IdentityKey"])

	return nil
}

func (l *LoadOPKResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

func (l *LoadOPKRequest) Unmarshal(params map[string][]byte) error {
	copy(l.SessionToken[:], params["SessionToken"])

	OPKArr, err := message.ParseArray(params["OPK"])
	if err != nil {
		return err
	}

	l.OPK = make([]OPKPair, 0, len(OPKArr))
	for _, opk := range OPKArr {
		opkParams, err := message.GetParams(opk)
		if err != nil {
			return err
		}

		var opkPair OPKPair
		err = opkPair.Unmarshal(opkParams)
		if err != nil {
			return err
		}

		l.OPK = append(l.OPK, opkPair)
	}

	return nil
}

func (i *RegisterRequest) Unmarshal(params map[string][]byte) error {
	i.Nickname = string(params["Nickname"])
	i.Number = string(params["Number"])
	i.Code = uint8(binary.BigEndian.Uint16(params["Code"]))
	i.PassHash = string(params["PassHash"])

	i.DeviceID = string(params["DeviceID"])

	return nil
}

func (r *RequiredOPKRequest) Unmarshal(params map[string][]byte) error {
	copy(r.SessionToken[:], params["SessionToken"])

	return nil
}

func (r *RequiredOPKResponse) Unmarshal(params map[string][]byte) error {
	return nil
}

func (i *RegisterResponse) Unmarshal(params map[string][]byte) error {
	copy(i.SessionToken[:], params["SessionToken"])

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
	copy(i.SessionToken[:], params["SessionToken"])

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
