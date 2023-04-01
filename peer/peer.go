package peer

import (
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/gen/hash"
	"github.com/oringik/crypto-chateau/transport/handshake"
	"github.com/oringik/crypto-chateau/transport/pipe"
	"net"
	"time"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/version"
)

var (
	ErrByte byte = 0x2F
	OkByte  byte = 0x20
)

type Peer struct {
	Pipe    *pipe.Pipe
	CloseCh chan bool
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Pipe:    pipe.NewPipe(conn),
		CloseCh: make(chan bool, 1),
	}
}

func (p *Peer) EstablishSecureConn() error {
	securedConnect, err := handshake.ClientHandshake(p.Pipe.GetConn())
	if err != nil {
		return err
	}

	p.Pipe.SetConn(securedConnect)

	return nil
}

func (p *Peer) SendRequestClient(handlerHash hash.HandlerHash, msg message.Message) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, handlerHash[:]...)
	resp = append(resp, msg.Marshal()...)

	_, err := p.Write(resp)
	return err
}

func (p *Peer) WriteResponse(msg message.Message) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, OkByte)
	resp = append(resp, msg.Marshal()...)

	_, err := p.Write(resp)
	return err
}

func (p *Peer) WriteRawResponse(data []byte) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, OkByte)
	resp = append(resp, data...)

	_, err := p.Write(resp)
	return err
}

func (p *Peer) ReadMessage(msg message.Message) error {
	msgRaw, err := p.Read(2048)
	if err != nil {
		return fmt.Errorf("failed to read from connection: %w", err)
	}

	_, offset, err := conv.GetServerRespMetaInfo(msgRaw)
	if err != nil {
		return err
	}

	// check if error prefix is present
	if msgRaw[offset] == ErrByte {
		return fmt.Errorf("chateau rpc: status = error, description = %s", string(msgRaw[2:]))
	}

	// check if message has a size
	if len(msgRaw) < offset+1+conv.ObjectBytesPrefixLength {
		return errors.New("not enough for size and message")
	}

	err = msg.Unmarshal(conv.NewBinaryIterator(msgRaw[offset+1+conv.ObjectBytesPrefixLength:]))
	if err != nil {
		return err
	}

	return err
}

func (p *Peer) WriteError(err error) error {
	var resp []byte

	resp = append(resp, version.NewProtocolByte())
	resp = append(resp, ErrByte)
	resp = append(resp, []byte(err.Error())...)

	_, writeErr := p.Write(resp)

	return writeErr
}

func (p *Peer) Write(data []byte) (int, error) {
	n, err := p.Pipe.Write(data)
	if err != nil {
		p.CloseCh <- true
		p.Close()
	}

	return n, err
}

func (p *Peer) Read(bufSize int) ([]byte, error) {
	msg, err := p.Pipe.Read(pipe.PipeReadCfg{BufSize: bufSize})
	if err != nil {
		p.CloseCh <- true
		p.Close()
	}

	return msg, err
}

func (p *Peer) Close() error {
	err := p.Pipe.CloseConn()

	return err
}

func (p *Peer) SetReadDeadline(t time.Time) error {
	return p.Pipe.SetReadDeadline(t)
}

func (p *Peer) LocalAddr() net.Addr {
	return p.Pipe.LocalAddr()
}

func (p *Peer) RemoteAddr() net.Addr {
	return p.Pipe.RemoteAddr()
}
