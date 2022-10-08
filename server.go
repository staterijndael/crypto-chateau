package crypto_chateau

import (
	"context"
	"errors"
	"github.com/Oringik/crypto-chateau/generated"
	"github.com/Oringik/crypto-chateau/message"
	"github.com/Oringik/crypto-chateau/peer"
	"github.com/Oringik/crypto-chateau/transport"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Config   *Config
	Handlers map[string]*generated.Handler
	// key: ip address  value: client peer
	Clients    map[string]*peer.Peer
	shutdownCh chan error
	logger     *zap.Logger
}

type Config struct {
	IP   string
	Port int
}

func NewServer(cfg *Config, logger *zap.Logger) *Server {
	return &Server{
		Config:     cfg,
		Handlers:   make(map[string]*generated.Handler),
		Clients:    make(map[string]*peer.Peer),
		shutdownCh: make(chan error),
		logger:     logger,
	}
}

func (s *Server) Run(ctx context.Context, endpoint generated.Endpoint) error {
	_, err := net.ResolveTCPAddr("tcp", s.Config.IP+":"+strconv.Itoa(s.Config.Port))
	if err != nil {
		return err
	}

	generated.InitHandlers(endpoint, s.Handlers)

	wg := sync.WaitGroup{}

	wg.Add(1)

	clientCh := make(chan *peer.Peer)

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

func (s *Server) handleRequests(ctx context.Context, clientChan <-chan *peer.Peer) {
	for {
		//select {
		//case <-ctx.Done():
		//	return
		//case client := <-clientChan:
		//	go s.handleRequest(ctx, client)
		client := <-clientChan
		go s.handleRequest(ctx, client)
	}
}

func (s *Server) handleRequest(ctx context.Context, peer *peer.Peer) {
	defer peer.Close()

	securedConnect, err := transport.ClientHandshake(peer.Conn)
	if err != nil {
		s.logger.Info("error establishing secured connect",
			zap.String("connIP", peer.Conn.RemoteAddr().String()),
			zap.Error(err),
		)
		return
	}

	peer.Conn = securedConnect

	err = s.handleMethod(ctx, peer)
	if err != nil {
		s.logger.Info("error handling method for peer",
			zap.String("connIP", peer.Conn.RemoteAddr().String()),
			zap.Error(err),
		)
		return
	}
}

func (s *Server) handleMethod(ctx context.Context, peer *peer.Peer) error {
	msg := make([]byte, 1024)
	n, err := peer.Read(msg)
	if err != nil {
		return err
	}

	msg = msg[:n]

	handlerName, n, err := message.GetHandlerName(msg)
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

	reqMsgParams, err := message.GetParams(msg[n:])
	if err != nil {
		return err
	}

	requestMsg := handler.RequestMsgType

	err = requestMsg.Unmarshal(reqMsgParams)
	if err != nil {
		return err
	}

	switch handler.HandlerType {
	case generated.HandlerT:
		responseMessage, err := handler.CallFuncHandler(ctx, requestMsg)
		if err != nil {
			writeErr := peer.WriteError(string(handlerName), err)
			return writeErr
		}

		err = peer.WriteResponse(responseMessage)
		if err != nil {
			return err
		}
	case generated.StreamT:
		if _, ok := handler.RequestMsgType.(generated.StreamReq); ok {
			streamReq := handler.RequestMsgType.(generated.StreamReq)
			err = streamReq.Init(peer, requestMsg)
			if err != nil {
				writeErr := peer.WriteError(string(handlerName), err)
				return writeErr
			}

			err = handler.CallFuncStream(ctx, streamReq)
			if err != nil {
				writeErr := peer.WriteError(string(handlerName), err)
				return writeErr
			}
		} else {
			return errors.New("unknown type of callFunc")
		}
	default:
		return errors.New("incorrect handler format: InternalError")
	}

	return nil
}

func (s *Server) listenClients(clientChan chan<- *peer.Peer) error {
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

		peer := peer.NewPeer(conn)

		clientChan <- peer
	}
}
