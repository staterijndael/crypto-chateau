package transport

import (
	"errors"
	"github.com/Oringik/crypto-chateau/dh"
	"math/big"
	"net"
	"time"
)

const (
	maxReadTime  = 5 * time.Second
	maxWriteTime = 3 * time.Second
)

func ClientHandshake(tcpConn net.Conn, keyStore *dh.KeyStore) (net.Conn, error) {
	conn := newConn(tcpConn, connCfg{readDeadline: maxReadTime, writeDeadline: maxWriteTime})

	if !dh.IsKeyValid(keyStore.PrivateKey) {
		return nil, errors.New("incorrect private key")
	}
	if !dh.IsKeyValid(keyStore.PublicKey) {
		return nil, errors.New("incorrect public key")
	}

	msg := make([]byte, 9)
	_, err := conn.Read(msg)
	if err != nil {
		return nil, err
	}
	if string(msg) != "handshake" {
		return nil, errors.New("incorrect init message")
	}

	dhParams := dhParamsInitMsg{g: dh.Generator, pHash: dh.PrimeHash}
	_, err = conn.Write(formatMsg(dhParams.g.Bytes(), dhParams.pHash))
	if err != nil {
		return nil, err
	}

	publicKeyMsg := publicKeyInitMsg{publicKey: keyStore.PublicKey}

	_, err = conn.Write(formatMsg(publicKeyMsg.publicKey.Bytes()))
	if err != nil {
		return nil, err
	}

	connPublicKey, err := readConnBigInt(conn)
	if err != nil {
		return nil, err
	}

	if !dh.IsKeyValid(connPublicKey) {
		return nil, errors.New("invalid public key")
	}

	_, err = conn.Write([]byte{'1'})
	if err != nil {
		return nil, err
	}

	err = keyStore.GenerateSharedKey(connPublicKey)
	if err != nil {
		return nil, err
	}

	err = conn.enableEncryption(keyStore.SharedKey)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func readConnBigInt(conn *Conn) (*big.Int, error) {
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	convertedBigIntBytes := new(big.Int)
	convertedBigIntBytes.SetBytes(buf)

	return convertedBigIntBytes, nil
}
