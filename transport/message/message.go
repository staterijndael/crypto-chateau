package message

import (
	"io"
	"net"
)

type MessageController struct {
	reservedData       []byte
	futurePacketLength int
}

func (m *MessageController) GetFullMessage(tcpConn net.Conn, bufSize int, bufSizeRead int) ([]byte, error) {
	if bufSize == 0 {
		bufSize = 1024
	}

	buf := make([]byte, 0, bufSize+2)

	for {
		if len(m.reservedData) > 0 {
			if m.futurePacketLength == 0 {
				packetLength := uint16(m.reservedData[0]) | uint16(m.reservedData[1])<<8
				m.futurePacketLength = int(packetLength)
				m.reservedData = m.reservedData[2:]
			}

			buf = append(buf, m.reservedData...)
			m.reservedData = []byte{}

			if len(buf) >= m.futurePacketLength {
				oldFuturePacketLength := m.futurePacketLength
				m.futurePacketLength = 0
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

		if m.futurePacketLength == 0 {
			m.futurePacketLength = int(uint16(buf[0]) | uint16(buf[1])<<8)
			buf = buf[2:]
		}

		if len(buf) >= m.futurePacketLength {
			oldFuturePacketLength := m.futurePacketLength
			m.futurePacketLength = 0
			if len(buf) != oldFuturePacketLength {
				m.reservedData = buf[oldFuturePacketLength:]
			}
			return buf[:oldFuturePacketLength], nil
		}
	}
}
