package crypto_chateau

import (
	"context"
	"errors"
	"github.com/Oringik/crypto-chateau/dh"
	"github.com/Oringik/crypto-chateau/generated"
	"github.com/Oringik/crypto-chateau/transport"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Config   *Config
	Handlers map[string]*Handler
	KeyStore *dh.KeyStore
	// key: ip address  value: client peer
	Clients    map[string]*Peer
	shutdownCh chan error
	logger     *zap.Logger
}

type Config struct {
	IP   string
	Port int
}

func NewServer(cfg *Config, logger *zap.Logger) *Server {
	keyStore := &dh.KeyStore{}

	keyStore.GeneratePrivateKey()
	keyStore.GeneratePublicKey()

	return &Server{
		Config:     cfg,
		KeyStore:   keyStore,
		Handlers:   make(map[string]*Handler),
		Clients:    make(map[string]*Peer),
		shutdownCh: make(chan error),
		logger:     logger,
	}
}

func (s *Server) Run(ctx context.Context, endpoint generated.Endpoint) error {
	_, err := net.ResolveTCPAddr("tcp", s.Config.IP+":"+strconv.Itoa(s.Config.Port))
	if err != nil {
		return err
	}

	initHandlers(endpoint, s.Handlers)

	wg := sync.WaitGroup{}

	wg.Add(1)

	clientCh := make(chan *Peer)

	go func() {
		err := s.listenClients(clientCh)
		if err != nil {
			s.shutdownCh <- err
		}
		wg.Done()
	}()

	s.handleRequests(ctx, clientCh)

	wg.Wait()

	return nil
}

func (s *Server) handleRequests(ctx context.Context, clientChan <-chan *Peer) {
	for {
		//select {
		//case <-ctx.Done():
		//	return
		//case client := <-clientChan:
		//	go s.handleRequest(ctx, client)
		//default:
		//	continue
		//}
		client := <-clientChan
		go s.handleRequest(ctx, client)
	}
}

func (s *Server) handleRequest(ctx context.Context, peer *Peer) {
	defer peer.Close()

	securedConnect, err := transport.ClientHandshake(peer.conn, s.KeyStore)
	if err != nil {
		s.logger.Info("error establishing secured connect",
			zap.String("connIP", peer.conn.RemoteAddr().String()),
			zap.Error(err),
		)
		return
	}

	peer.conn = securedConnect

	err = s.handleMethod(ctx, peer)
	if err != nil {
		s.logger.Info("error handling method for peer",
			zap.String("connIP", peer.conn.RemoteAddr().String()),
			zap.Error(err),
		)
		return
	}
}

func (s *Server) handleMethod(ctx context.Context, peer *Peer) error {
	msg := make([]byte, 1024)
	n, err := peer.Read(msg)
	if err != nil {
		return err
	}

	msg = msg[:n]

	handlerName, n, err := GetHandlerName(msg)
	if err != nil {
		return err
	}

	handler, ok := s.Handlers[string(handlerName)]
	if !ok {
		return errors.New("unknown handler " + string(handlerName))
	}

	if n >= len(msg) {
		return errors.New("incorrect message")
	}

	requestMsg, err := ParseMessage(msg[n:], handler.requestMsgType)
	if err != nil {
		return err
	}

	switch handler.HandlerType {
	case HandlerT:
		fnc, err := callFuncToHandlerFunc(handler.callFunc)
		if err != nil {
			return err
		}

		responseMsg, err := fnc(ctx, requestMsg)
		if err != nil {
			writeErr := peer.WriteError(err)
			return writeErr
		}

		err = peer.WriteResponse(responseMsg)
		if err != nil {
			return err
		}
	case StreamT:
		fnc := handler.callFunc.(func(context.Context, *Stream) error)
		stream := &Stream{
			peer: peer,
		}
		err := fnc(ctx, stream)
		if err != nil {
			writeErr := peer.WriteError(err)
			return writeErr
		}
	default:
		return errors.New("incorrect handler format: InternalError")
	}

	return nil
}

func (s *Server) listenClients(clientChan chan<- *Peer) error {
	listener, err := net.Listen("tcp", s.Config.IP+":"+strconv.Itoa(s.Config.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
		}

		peer := NewPeer(conn)

		clientChan <- peer
	}
}
