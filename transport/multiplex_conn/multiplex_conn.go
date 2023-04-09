package multiplex_conn

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type MultiplexConnPool struct {
	tcpConn          net.Conn
	currentRequestID atomic.Uint32

	toWriteQueue  chan ToWriteMsg
	listenClients chan *MultiplexConn

	multiplexConnByRequestID   map[uint32]*MultiplexConn
	multiplexConnByRequestIDMx sync.RWMutex

	closeConnsCh chan uint32
	terminateCh  chan bool

	isClient bool
}

func NewMultiplexConnPool(tcpConn net.Conn, isClient bool) *MultiplexConnPool {
	return &MultiplexConnPool{
		tcpConn:          tcpConn,
		currentRequestID: atomic.Uint32{},

		toWriteQueue: make(chan ToWriteMsg, 4096),

		listenClients: make(chan *MultiplexConn),

		multiplexConnByRequestID:   make(map[uint32]*MultiplexConn),
		multiplexConnByRequestIDMx: sync.RWMutex{},

		closeConnsCh: make(chan uint32),
		terminateCh:  make(chan bool),

		isClient: isClient,
	}
}

func (p *MultiplexConnPool) NewMultiplexConn() *MultiplexConn {
	requestID := p.currentRequestID.Add(1)

	newMultiplexConn := &MultiplexConn{
		requestID:          requestID,
		localAddr:          p.tcpConn.LocalAddr(),
		remoteAddr:         p.tcpConn.RemoteAddr(),
		readQueue:          make(chan []byte, 5),
		writeQueue:         p.toWriteQueue,
		errChan:            make(chan error),
		connReservedDataMx: sync.Mutex{},
		connReservedData:   []byte{},
		closeCh:            p.closeConnsCh,
		readDeadline:       2 * time.Minute,
		isClosedMx:         sync.RWMutex{},
		closedNotifyCh:     make(chan bool, 1),
	}

	p.multiplexConnByRequestIDMx.Lock()
	p.multiplexConnByRequestID[requestID] = newMultiplexConn
	p.multiplexConnByRequestIDMx.Unlock()

	if !p.isClient {
		p.listenClients <- newMultiplexConn
	}

	return newMultiplexConn
}

func (p *MultiplexConnPool) SetRawTCPDeadline(t time.Time) error {
	return p.tcpConn.SetDeadline(t)
}

func (p *MultiplexConnPool) Close() {
	p.terminateCh <- true
}

func (p *MultiplexConnPool) ListenClients() chan *MultiplexConn {
	return p.listenClients
}

func (p *MultiplexConnPool) Run() {
	go func() {
		for {
			select {
			case toWriteMsg := <-p.toWriteQueue:
				dataWithRequestID := make([]byte, 0, len(toWriteMsg.Data)+2)
				dataWithRequestID = append(dataWithRequestID, byte(uint16(toWriteMsg.RequestID)>>8), byte(uint16(toWriteMsg.RequestID)))
				dataWithRequestID = append(dataWithRequestID, toWriteMsg.Data...)

				_, err := p.tcpConn.Write(dataWithRequestID)
				p.multiplexConnByRequestIDMx.RLock()
				multiplexConn, ok := p.multiplexConnByRequestID[toWriteMsg.RequestID]
				if !ok {
					fmt.Println("multiplex conn not found: requestID: " + strconv.Itoa(int(toWriteMsg.RequestID)))
					p.multiplexConnByRequestIDMx.RUnlock()
					continue
				}
				p.multiplexConnByRequestIDMx.RUnlock()

				multiplexConn.errChan <- err
			case requestID := <-p.closeConnsCh:
				p.multiplexConnByRequestIDMx.Lock()
				multiplexConn, ok := p.multiplexConnByRequestID[requestID]
				if ok {
					multiplexConn.isClosedMx.Lock()
					multiplexConn.isClosed = true
					multiplexConn.isClosedMx.Unlock()
				}

				multiplexConn.closedNotifyCh <- true

				delete(p.multiplexConnByRequestID, requestID)
				p.multiplexConnByRequestIDMx.Unlock()
			case <-p.terminateCh:
				p.multiplexConnByRequestIDMx.Lock()
				for _, multiplexConn := range p.multiplexConnByRequestID {
					multiplexConn.isClosedMx.Lock()
					multiplexConn.isClosed = true
					multiplexConn.isClosedMx.Unlock()

					multiplexConn.closedNotifyCh <- true
				}
				p.multiplexConnByRequestIDMx.Unlock()
				p.tcpConn.Close()
				close(p.toWriteQueue)
				close(p.listenClients)
				close(p.closeConnsCh)
				return
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := p.tcpConn.Read(buf)
			if err != nil {
				p.terminateCh <- true
				return
			}

			if n == 0 {
				break
			}

			buf = buf[:n]

			requestID := uint16(buf[0])<<8 | uint16(buf[1])
			p.multiplexConnByRequestIDMx.Lock()
			if multiplexConn, ok := p.multiplexConnByRequestID[uint32(requestID)]; ok {
				select {
				case multiplexConn.readQueue <- buf[2:]:
				}
			} else {
				newMultiplexConn := &MultiplexConn{
					requestID:          uint32(requestID),
					localAddr:          p.tcpConn.LocalAddr(),
					remoteAddr:         p.tcpConn.RemoteAddr(),
					readQueue:          make(chan []byte, 500),
					writeQueue:         p.toWriteQueue,
					errChan:            make(chan error),
					connReservedDataMx: sync.Mutex{},
					connReservedData:   []byte{},
					closeCh:            p.closeConnsCh,
					isClosedMx:         sync.RWMutex{},
					readDeadline:       2 * time.Minute,
					closedNotifyCh:     make(chan bool, 1),
				}

				newMultiplexConn.readQueue <- buf[2:]
				p.multiplexConnByRequestID[uint32(requestID)] = newMultiplexConn

				p.listenClients <- newMultiplexConn
			}
			p.multiplexConnByRequestIDMx.Unlock()
		}
	}()
}

type MultiplexConn struct {
	localAddr  net.Addr
	remoteAddr net.Addr

	requestID uint32

	readQueue  chan []byte
	writeQueue chan<- ToWriteMsg

	errChan chan error

	connReservedDataMx sync.Mutex
	connReservedData   []byte

	readDeadline time.Duration

	closeCh chan<- uint32

	isClosed   bool
	isClosedMx sync.RWMutex

	closedNotifyCh chan bool
}

func (cn *MultiplexConn) Write(p []byte) (int, error) {
	cn.isClosedMx.RLock()
	if cn.isClosed {
		cn.isClosedMx.RUnlock()
		return 0, errors.New("conn is closed")
	}
	cn.isClosedMx.RUnlock()

	select {
	case cn.writeQueue <- ToWriteMsg{
		RequestID: cn.requestID,
		Data:      p,
	}:
	default:
		return 0, errors.New("writing to multiplex conn was blocked")
	}

	err, ok := <-cn.errChan
	if !ok {
		return 0, errors.New("reading from errChan in during write in multiplexConn was blocked")
	}

	return 0, err
}

func (cn *MultiplexConn) Read(b []byte) (int, error) {
	cn.isClosedMx.RLock()
	if cn.isClosed {
		cn.isClosedMx.RUnlock()
		return 0, errors.New("conn is closed")
	}
	cn.isClosedMx.RUnlock()

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

	select {
	case data := <-cn.readQueue:
		returnLength := len(data)

		cn.connReservedDataMx.Lock()
		if len(data) > len(b) {
			cn.connReservedData = data[len(b):]
			returnLength = len(b)
		}
		cn.connReservedDataMx.Unlock()

		copy(b, data)

		return returnLength, nil
	case <-time.After(cn.readDeadline):
		return 0, io.EOF
	}
}

func (cn *MultiplexConn) Close() error {
	cn.closeCh <- cn.requestID
	return nil
}

func (cn *MultiplexConn) LocalAddr() net.Addr {
	return cn.localAddr
}

func (cn *MultiplexConn) RemoteAddr() net.Addr {
	return cn.remoteAddr
}

func (cn *MultiplexConn) SetDeadline(t time.Time) error {
	return nil
}

func (cn *MultiplexConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (cn *MultiplexConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (cn *MultiplexConn) ClosedNotifyChannel() <-chan bool {
	return cn.closedNotifyCh
}
