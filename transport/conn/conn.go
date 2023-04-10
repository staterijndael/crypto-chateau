package conn

import (
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/aes-256"
	"github.com/oringik/crypto-chateau/transport"
	"github.com/oringik/crypto-chateau/transport/message"
	"net"
	"strconv"
	"sync"
	"time"
)

type Conn struct {
	tcpConn    net.Conn
	cfg        ConnCfg
	encryption encryption

	connReservedDataMx sync.Mutex
	connReservedData   []byte

	msgController *message.MessageController
}

type ConnCfg struct {
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
}

type encryption struct {
	enabled   bool
	sharedKey []byte
}

func NewConn(tcpConn net.Conn, cfg ConnCfg) *Conn {
	return &Conn{
		tcpConn:            tcpConn,
		cfg:                cfg,
		msgController:      &message.MessageController{},
		connReservedDataMx: sync.Mutex{},
	}
}

func (cn *Conn) EnableEncryption(sharedKey [32]byte) error {
	if cn.encryption.enabled {
		return errors.New("encryption already enabled")
	}

	sharedKeyHash, err := transport.GetSha256FromBytes(sharedKey)
	if err != nil {
		return err
	}

	cn.encryption.enabled = true
	cn.encryption.sharedKey = sharedKeyHash

	return nil
}

func (cn *Conn) Write(p []byte) (int, error) {
	before := p
	if cn.encryption.enabled {
		encryptedData, err := aes_256.Encrypt(p, cn.encryption.sharedKey)
		if err != nil {
			return 0, err
		}

		p = encryptedData
	}

	dataWithLength := make([]byte, 0, len(p)+2)
	convertedLength := uint16(len(p))
	dataWithLength = append(dataWithLength, byte(convertedLength), byte(convertedLength>>8))
	dataWithLength = append(dataWithLength, p...)

	n, err := cn.tcpConn.Write(dataWithLength)
	if err != nil {
		return 0, err
	}

	fmt.Println(cn.tcpConn.RemoteAddr(), " - ", before)

	return n, nil
}

func (cn *Conn) Read(b []byte) (int, error) {
	cn.connReservedDataMx.Lock()
	if len(cn.connReservedData) > 0 {
		defer cn.connReservedDataMx.Unlock()

		if len(cn.connReservedData) > len(b) {
			copy(b, cn.connReservedData)
			cn.connReservedData = cn.connReservedData[len(b):]
			return len(b), nil
		}

		copy(b, cn.connReservedData)
		return len(cn.connReservedData), nil
	}
	cn.connReservedDataMx.Unlock()

	fullMsg, err := cn.msgController.GetFullMessage(cn.tcpConn, len(b), 4096)
	if err != nil {
		return 0, err
	}

	if cn.encryption.enabled {
		decryptedData, err := aes_256.Decrypt(fullMsg, cn.encryption.sharedKey)
		if err != nil {
			return 0, err
		}

		fullMsg = decryptedData
	}

	if len(fullMsg) > 4*4096 {
		return 0, errors.New("buffer overflow: bufSize - " + strconv.Itoa(len(b)))
	}

	returnLength := len(fullMsg)

	cn.connReservedDataMx.Lock()
	if len(fullMsg) > len(b) {
		cn.connReservedData = fullMsg[len(b):]
		returnLength = len(b)
	}
	cn.connReservedDataMx.Unlock()

	copy(b, fullMsg)
	return returnLength, nil
}

func (cn *Conn) Close() error {
	return cn.tcpConn.Close()
}

func (cn *Conn) LocalAddr() net.Addr {
	return cn.tcpConn.LocalAddr()
}

func (cn *Conn) RemoteAddr() net.Addr {
	return cn.tcpConn.RemoteAddr()
}

func (cn *Conn) SetDeadline(t time.Time) error {
	return cn.tcpConn.SetDeadline(t)
}

func (cn *Conn) SetReadDeadline(t time.Time) error {
	return cn.tcpConn.SetReadDeadline(t)
}

func (cn *Conn) SetWriteDeadline(t time.Time) error {
	return cn.tcpConn.SetWriteDeadline(t)
}
