package peer

import (
	"errors"
	"net"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/gen/hash"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/version"
)

var (
	ErrBytesPrefix = [2]byte{0x2F, 0x20}
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

	// check if error prefix is present
	if msgRaw[offset] == ErrBytesPrefix[0] && msgRaw[offset+1] == ErrBytesPrefix[1] {
		return errors.New(string(msgRaw[offset+2:]))
	}

	// check if message has a size
	if len(msgRaw) < offset+conv.ObjectBytesPrefixLength {
		return errors.New("not enough for size and message")
	}

	err = msg.Unmarshal(conv.NewBinaryIterator(msgRaw[offset+conv.ObjectBytesPrefixLength:]))
	if err != nil {
		return err
	}

	return err
}

func (p *Peer) WriteError(handlerKey hash.HandlerHash, err error) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, handlerKey[:]...)
	resp = append(resp, ErrBytesPrefix[:]...)
	resp = append(resp, []byte(err.Error())...)

	_, writeErr := p.Conn.Write(resp)

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
