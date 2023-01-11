package peer

import (
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/gen/hash"
	"github.com/oringik/crypto-chateau/transport"
	"net"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/version"
)

var (
	ErrByte byte = 0x2F
	OkByte  byte = 0x20
)

type Peer struct {
	Pipe *transport.Pipe
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Pipe: transport.NewPipe(conn),
	}
}

func (p *Peer) EstablishSecureConn() error {
	securedConnect, err := transport.ClientHandshake(p.Pipe.GetConn())
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

func (p *Peer) ReadMessage(msg message.Message) error {
	msgRaw, err := p.Read()
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
	if offset+1+len(msgRaw) < conv.ObjectBytesPrefixLength {
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

	return n, err
}

func (p *Peer) Read() ([]byte, error) {
	msg, err := p.Pipe.Read(transport.PipeReadCfg{})

	return msg, err
}

func (p *Peer) Close() error {
	err := p.Pipe.CloseConn()

	return err
}
