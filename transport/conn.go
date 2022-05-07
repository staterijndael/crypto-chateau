package transport

import (
	"context"
	aes_256 "crypto-chateau/aes-256"
	"crypto-chateau/dh"
	"errors"
	"github.com/xelaj/go-dry/ioutil"
	"math/big"
	"net"
	"time"
)

type Conn struct {
	tcpConn    net.Conn
	reader     *ioutil.CancelableReader
	cfg        connCfg
	encryption encryption
}

type connCfg struct {
	writeDeadline time.Duration
	readDeadline  time.Duration
}

type encryption struct {
	enabled   bool
	sharedKey []byte
}

func newConn(ctx context.Context, tcpConn net.Conn, cfg connCfg) *Conn {
	reader := ioutil.NewCancelableReader(ctx, tcpConn)

	return &Conn{
		tcpConn: tcpConn,
		reader:  reader,
		cfg:     cfg,
	}
}

func (cn *Conn) enableEncryption(sharedKey *big.Int) error {
	if cn.encryption.enabled {
		return errors.New("encryption already enabled")
	}

	if !dh.IsKeyValid(sharedKey) {
		return errors.New("invalid shared key")
	}

	sharedKeyBytes, err := getSha256FromBigInt(sharedKey)
	if err != nil {
		return err
	}

	cn.encryption.enabled = true
	cn.encryption.sharedKey = sharedKeyBytes

	return nil
}

func (cn *Conn) Write(p []byte) (int, error) {
	if cn.cfg.writeDeadline > 0 {
		err := cn.SetWriteDeadline(time.Now().Add(cn.cfg.writeDeadline))
		if err != nil {
			return 0, err
		}
	}

	var data []byte

	if cn.encryption.enabled {
		encryptedData, err := aes_256.Encrypt(p, cn.encryption.sharedKey)
		if err != nil {
			return 0, err
		}

		data = encryptedData
	} else {
		data = p
	}

	n, err := cn.tcpConn.Write(data)
	return n, err
}

func (cn *Conn) Read(p []byte) (int, error) {
	if cn.cfg.readDeadline > 0 {
		err := cn.SetReadDeadline(time.Now().Add(cn.cfg.readDeadline))
		if err != nil {
			return 0, err
		}
	}

	if cn.encryption.enabled {
		buf := make([]byte, len(p))

		n, err := cn.reader.Read(buf)
		if err != nil {
			return 0, err
		}

		decryptedData, err := aes_256.Decrypt(buf, cn.encryption.sharedKey)
		if err != nil {
			return 0, err
		}

		copy(p, decryptedData)

		return n, nil
	} else {
		n, err := cn.reader.Read(p)

		return n, err
	}
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
