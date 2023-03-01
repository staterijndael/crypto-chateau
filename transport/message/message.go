package message

import (
	"io"
	"net"
)

type MessageController struct {
	reservedData []byte
}

func (m *MessageController) GetFullMessage(tcpConn net.Conn, bufSize int, bufSizeRead int) ([]byte, error) {
	futurePacketLength := 0

	buf := make([]byte, 0, bufSize+2)

	for {
		if len(m.reservedData) > 0 {
			if futurePacketLength == 0 {
				packetLength := uint16(m.reservedData[0]) | uint16(m.reservedData[1])<<8
				futurePacketLength = int(packetLength)
				m.reservedData = m.reservedData[2:]
			}

			buf = append(buf, m.reservedData...)
			m.reservedData = []byte{}

			if len(buf) >= futurePacketLength {
				oldFuturePacketLength := futurePacketLength
				futurePacketLength = 0
				if len(buf) != oldFuturePacketLength {
					m.reservedData = buf[oldFuturePacketLength:]
				}
				return buf[:oldFuturePacketLength], nil
			}
		}

		localBuf := make([]byte, bufSizeRead)

		n, err := tcpConn.Read(localBuf)
		if err != nil {
			return nil, err
		}

		localBuf = localBuf[:n]

		buf = append(buf, localBuf...)

		if len(buf) == 0 {
			return nil, io.EOF
		}

		if futurePacketLength == 0 {
			futurePacketLength = int(uint16(buf[0]) | uint16(buf[1])<<8)
			buf = buf[2:]
		}

		if len(buf) >= futurePacketLength {
			oldFuturePacketLength := futurePacketLength
			futurePacketLength = 0
			if len(buf) != oldFuturePacketLength {
				m.reservedData = buf[oldFuturePacketLength:]
			}
			return buf[:oldFuturePacketLength], nil
		}
	}
}
