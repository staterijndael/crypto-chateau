package handshake

import (
	"crypto/rand"
	"errors"
	"github.com/oringik/crypto-chateau/dh"
	conn2 "github.com/oringik/crypto-chateau/transport/conn"
	pipe2 "github.com/oringik/crypto-chateau/transport/pipe"
	"golang.org/x/crypto/curve25519"
	"io"
	"net"
	"time"
)

type publicKeyInitMsg struct {
	publicKey [32]byte
}

const (
	maxReadTime  = 5 * time.Second
	maxWriteTime = 3 * time.Second
)

func ClientHandshake(tcpConn net.Conn) (net.Conn, error) {
	conn := conn2.NewConn(tcpConn, conn2.ConnCfg{ReadDeadline: maxReadTime, WriteDeadline: maxWriteTime})

	pipe := pipe2.NewPipe(conn)

	msg, err := pipe.Read(pipe2.PipeReadCfg{BufSize: 9})
	if err != nil {
		return nil, err
	}
	if string(msg) != "handshake" {
		return nil, errors.New("incorrect init message")
	}

	var priv [32]byte
	if _, err := io.ReadFull(rand.Reader, priv[:]); err != nil {
		panic(err)
	}

	priv[0] &= 248
	priv[31] &= 63
	priv[31] |= 64

	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, &priv)

	publicKeyMsg := publicKeyInitMsg{publicKey: pub}

	_, err = pipe.Write(publicKeyMsg.publicKey[:])
	if err != nil {
		return nil, err
	}

	connPublicKey, err := readConnPubKey(pipe)
	if err != nil {
		return nil, err
	}

	_, err = pipe.Write([]byte{'1'})
	if err != nil {
		return nil, err
	}

	sharedKey, err := dh.DH(priv, connPublicKey)
	if err != nil {
		return nil, err
	}

	err = conn.EnableEncryption(sharedKey)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ServerHandshake(tcpConn net.Conn) (net.Conn, error) {
	conn := conn2.NewConn(tcpConn, conn2.ConnCfg{ReadDeadline: maxReadTime, WriteDeadline: maxWriteTime})

	pipe := pipe2.NewPipe(conn)

	_, err := pipe.Write([]byte("handshake"))
	if err != nil {
		return nil, err
	}

	connPublicKey, err := readConnPubKey(pipe)
	if err != nil {
		return nil, err
	}

	var priv [32]byte
	if _, err := io.ReadFull(rand.Reader, priv[:]); err != nil {
		panic(err)
	}

	priv[0] &= 248
	priv[31] &= 63
	priv[31] |= 64

	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, &priv)

	publicKeyMsg := publicKeyInitMsg{publicKey: pub}

	_, err = pipe.Write(publicKeyMsg.publicKey[:])
	if err != nil {
		return nil, err
	}

	msg, err := pipe.Read(pipe2.PipeReadCfg{BufSize: 1})
	if err != nil {
		return nil, err
	}

	if msg[0] != '1' {
		return nil, errors.New("incorrect result byte")
	}

	sharedKey, err := dh.DH(priv, connPublicKey)
	if err != nil {
		return nil, err
	}

	err = conn.EnableEncryption(sharedKey)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func readConnPubKey(pipe *pipe2.Pipe) ([32]byte, error) {
	msg, err := pipe.Read(pipe2.PipeReadCfg{BufSize: 32})
	if err != nil {
		return [32]byte{}, err
	}

	var s [32]byte
	copy(s[:], msg)
	return s, nil
}
