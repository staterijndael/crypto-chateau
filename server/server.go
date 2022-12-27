package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/oringik/crypto-chateau/gen/conv"
	"github.com/oringik/crypto-chateau/gen/hash"
	"github.com/oringik/crypto-chateau/message"
	"github.com/oringik/crypto-chateau/peer"
	"github.com/oringik/crypto-chateau/transport"
)

type HandlerFunc func(context.Context, message.Message) (message.Message, error)
type StreamFunc func(ctx context.Context, peer *peer.Peer, message message.Message) error

type HandlerType int

var HandlerT HandlerType = 0
var StreamT HandlerType = 1

type Handler struct {
	CallFuncHandler HandlerFunc
	CallFuncStream  StreamFunc
	HandlerType
	RequestMsgType  message.Message
	ResponseMsgType message.Message
	Tags            map[string]string
}

type Server struct {
	Config   *Config
	Handlers map[hash.HandlerHash]*Handler
	// key: ip address  value: client peer
	Clients    map[string]*peer.Peer
	shutdownCh chan error
	logger     *zap.Logger
}

type Config struct {
	IP   string
	Port int

	ConnReadDeadline  *time.Duration
	ConnWriteDeadline *time.Duration
}

func NewServer(cfg *Config, logger *zap.Logger, handlers map[hash.HandlerHash]*Handler) *Server {
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
	msg := make([]byte, 0, 1024)

	for {
		buf := make([]byte, 1024)
		n, err := peer.Read(buf)
		if err != nil {
			return err
		}

		if n == 0 {
			break
		}

		if n < len(buf) {
			buf = buf[:n]
			msg = append(msg, buf...)
			break
		}

		msg = append(msg, buf...)
	}

	_, handlerKey, offset, err := conv.GetHandler(msg)
	if err != nil {
		return err
	}

	handler, ok := s.Handlers[handlerKey]
	if !ok {
		return errors.New(fmt.Sprintf("handler not found for key: %v", handlerKey))
	}

	// check if message has a size
	if len(msg) < offset+conv.ObjectBytesPrefixLength {
		return errors.New("not enough bytes for size and message")
	}

	requestMsg := handler.RequestMsgType

	err = requestMsg.Unmarshal(conv.NewBinaryIterator(msg[offset+conv.ObjectBytesPrefixLength:]))
	if err != nil {
		return err
	}

	switch handler.HandlerType {
	case HandlerT:
		responseMessage, err := handler.CallFuncHandler(ctx, requestMsg)
		if err != nil {
			writeErr := peer.WriteError(handlerKey, err)
			return writeErr
		}

		err = peer.WriteResponse(handlerKey, responseMessage)
		if err != nil {
			return err
		}

		if val, ok2 := s.Handlers[handlerKey].Tags["keep_conn_alive"]; !ok2 || val != "true" {
			err = peer.Close()
			if err != nil {
				return err
			}
		}
	case StreamT:
		go func() {
			err = handler.CallFuncStream(ctx, peer, requestMsg)
			if err != nil {
				writeErr := peer.WriteError(handlerKey, err)
				if writeErr != nil {
					fmt.Println(writeErr)
				}
				return
			}

			if val, ok2 := s.Handlers[handlerKey].Tags["keep_conn_alive"]; !ok2 || val != "true" {
				err = peer.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
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

		if s.Config.ConnReadDeadline != nil {
			err = conn.SetReadDeadline(time.Now().Add(*s.Config.ConnReadDeadline))
			if err != nil {
				s.logger.Error("error setting read deadline", zap.Error(err))
				continue
			}
		}
		if s.Config.ConnWriteDeadline != nil {
			err = conn.SetWriteDeadline(time.Now().Add(*s.Config.ConnWriteDeadline))
			if err != nil {
				s.logger.Error("error setting write deadline", zap.Error(err))
				continue
			}
		}

		peer := peer.NewPeer(conn)

		clientChan <- peer
	}
}
