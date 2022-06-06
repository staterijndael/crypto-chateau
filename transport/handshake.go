package transport

import (
	"context"
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

func ClientHandshake(ctx context.Context, tcpConn net.Conn) (net.Conn, error) {
	conn := newConn(ctx, tcpConn, connCfg{readDeadline: maxReadTime, writeDeadline: maxWriteTime})

	keyStore := dh.KeyStore{}
	keyStore.GeneratePrivateKey()

	_, err := conn.Write([]byte("handshake"))
	if err != nil {
		return nil, err
	}

	dhParams := dhParamsInitMsg{g: dh.Generator, p: dh.Prime}
	_, err = conn.Write(formatMsg(dhParams.g.Bytes(), dhParams.p.Bytes()))
	if err != nil {
		return nil, err
	}

	err = keyStore.GeneratePublicKey()
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
	buf := make([]byte, 258)
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	convertedBigIntBytes := new(big.Int)
	convertedBigIntBytes.SetBytes(buf)

	return convertedBigIntBytes, nil
}
