package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/peer"
	"github.com/oringik/crypto-chateau/transport"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

type HandlerFunc func(context.Context, message.Message) (message.Message, error)
type StreamFunc func(ctx context.Context, peer interface{}, message message.Message) error

type HandlerType int

var HandlerT HandlerType = 0
var StreamT HandlerType = 1

type Handler struct {
	CallFuncHandler HandlerFunc
	CallFuncStream  StreamFunc
	HandlerType
	RequestMsgType message.Message
}

type Server struct {
	Config   *Config
	Handlers map[string]*Handler
	// key: ip address  value: client peer
	Clients    map[string]*peer.Peer
	shutdownCh chan error
	logger     *zap.Logger
}

type Config struct {
	IP   string
	Port int
}

func NewServer(cfg *Config, logger *zap.Logger, handlers map[string]*Handler) *Server {
	return &Server{
		Config:     cfg,
		Handlers:   handlers,
		Clients:    make(map[string]*peer.Peer),
		shutdownCh: make(chan error),
		logger:     logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	_, err := net.ResolveTCPAddr("tcp", s.Config.IP+":"+strconv.Itoa(s.Config.Port))
	if err != nil {
		return err
	}

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

	handlerName, n, err := conv.GetHandlerName(msg)
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

	_, reqMsgParams, err := conv.GetParams(msg[n:])
	if err != nil {
		return err
	}

	requestMsg := handler.RequestMsgType

	err = requestMsg.Unmarshal(reqMsgParams)
	if err != nil {
		return err
	}

	switch handler.HandlerType {
	case HandlerT:
		responseMessage, err := handler.CallFuncHandler(ctx, requestMsg)
		if err != nil {
			writeErr := peer.WriteError(string(handlerName), err)
			return writeErr
		}

		err = peer.WriteResponse(string(handlerName), responseMessage)
		if err != nil {
			return err
		}
	case StreamT:
		go func() {
			err = handler.CallFuncStream(ctx, peer, requestMsg)
			if err != nil {
				writeErr := peer.WriteError(string(handlerName), err)
				if writeErr != nil {
					fmt.Println(writeErr)
				}
				return
			}
		}()
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