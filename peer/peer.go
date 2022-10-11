package peer

import (
	"fmt"
	"github.com/Oringik/crypto-chateau/message"
	"net"
)

type Peer struct {
	Conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) WriteResponse(handlerName string, msg message.Message) error {
	var resp []byte

	resp = append(resp, []byte(handlerName+"#")...)
	resp = append(resp, msg.Marshal()...)

	_, err := p.Conn.Write(resp)
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
