package peer

import (
	"fmt"
	"net"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/gen/hash"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/version"
)

type Peer struct {
	Conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) WriteResponse(handlerName hash.HandlerHash, msg message.Message) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, handlerName[:]...)
	resp = append(resp, msg.Marshal()...)

	_, err := p.Conn.Write(resp)
	return err
}

func (p *Peer) ReadMessage(msg message.Message) error {
	var msgRaw []byte

	_, err := p.Conn.Read(msgRaw)
	_, _, offset, err := conv.GetHandler(msgRaw)
	if err != nil {
		return err
	}

	err = msg.Unmarshal(conv.NewBinaryIterator(msgRaw[offset:]))
	if err != nil {
		return err
	}

	return err
}

func (p *Peer) WriteError(handlerName string, err error) error {
	msg := fmt.Sprintf("%s# error: %s", handlerName, err.Error())

	_, writeErr := p.Conn.Write([]byte(msg))

	return writeErr
}

func (p *Peer) Write(data []byte) (int, error) {
	n, err := p.Conn.Write(data)

	return n, err
}

func (p *Peer) Read(b []byte) (int, error) {
	n, err := p.Conn.Read(b)

	return n, err
}

func (p *Peer) Close() error {
	err := p.Conn.Close()

	return err
}
