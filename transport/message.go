package transport

import (
	"io"
	"net"
)

type FullMessage struct {
	msg                   []byte
	gotReservedData       []byte
	gotFuturePacketLength uint16
}

func GetFullMessage(tcpConn net.Conn, bufSize int, reservedData []byte, futurePacketLength uint16) (FullMessage, error) {
	if bufSize == 0 {
		bufSize = 1024
	}

	buf := make([]byte, 0, bufSize+4)

	for {
		if len(reservedData) > 0 {
			if futurePacketLength == 0 {
				packetLength := uint16(reservedData[0]) | uint16(reservedData[1])<<8
				futurePacketLength = packetLength
				reservedData = reservedData[2:]
			}

			buf = append(buf, reservedData...)
			reservedData = make([]byte, 0, bufSize)

			if uint16(len(buf)) >= futurePacketLength {
				if futurePacketLength != uint16(len(buf)) {
					reservedData = buf[futurePacketLength:]
				}
				return FullMessage{
					msg:                   buf[:futurePacketLength],
					gotFuturePacketLength: 0,
					gotReservedData:       reservedData,
				}, nil
			}
		}

		localBuf := make([]byte, bufSize)

		n, err := tcpConn.Read(localBuf)
		if err != nil {
			return FullMessage{}, nil
		}

		localBuf = localBuf[:n]

		buf = append(buf, localBuf...)

		if len(buf) == 0 {
			return FullMessage{}, io.EOF
		}

		if futurePacketLength == 0 {
			futurePacketLength = uint16(buf[0]) | uint16(buf[1])<<8
			buf = buf[2:]
		}

		if uint16(len(buf)) >= futurePacketLength {
			if futurePacketLength != uint16(len(buf)) {
				reservedData = buf[futurePacketLength:]
			}
			return FullMessage{
				msg:                   buf[:futurePacketLength],
				gotFuturePacketLength: 0,
				gotReservedData:       reservedData,
			}, nil
		}
	}
}
