package crypto_chateau

import (
	"context"
	"crypto-chateau/transport"
	"log"
	"net"
	"sync"
)

const (
	msgDelim = '\n'
)

type Server struct {
	Config   *Config
	Handlers map[string]*Handler
	// key: ip address  value: client peer
	Clients    map[string]*Peer
	shutdownCh chan struct{}
}

type Config struct {
	IP   string
	Port string
}

func (s *Server) Run(ctx context.Context, config Config) error {
	_, err := net.ResolveTCPAddr("tcp", config.IP+":"+config.Port)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	clientCh := make(chan *Peer)

	go func() {
		s.listenClients(ctx, clientCh)
		wg.Done()
	}()

	s.handleRequests(ctx, clientCh)

	wg.Wait()

	return nil
}

func (s *Server) handleRequests(ctx context.Context, clientChan <-chan *Peer) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

func (s *Server) handleRequest(ctx context.Context, peer *Peer) error {
	securedConnect, err := transport.ClientHandshake(ctx, peer.conn)
	if err != nil {
		return err
	}

	peer.conn = securedConnect

	return nil
}

func (s *Server) handleMethod(ctx context.Context, peer *Peer) error {
	//bytesMsg, err := bufio.NewReader(peer).ReadBytes(msgDelim)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Server) listenClients(ctx context.Context, clientChan chan<- *Peer) {
	listener, err := net.Listen("tcp", s.Config.IP+":"+s.Config.Port)
	if err != nil {
		log.Println(err)
		s.shutdownCh <- struct{}{}
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				log.Println("Failed to accept connection:", err.Error())
			}

			peer := NewPeer(conn)

			clientChan <- peer
		}
	}
}
