package pubsub

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	defaultAddress = ":10"
)

type Server struct {
	address string
	port    int

	ErrCh chan error

	ctx                context.Context
	shutdownFn         context.CancelFunc
	shutdownCompleteCh chan struct{}
}

func NewServer() *Server {
	ctx, shutdownFn := context.WithCancel(context.Background())

	return &Server{
		address: defaultAddress,

		ErrCh: make(chan error),

		ctx:                ctx,
		shutdownFn:         shutdownFn,
		shutdownCompleteCh: make(chan struct{}),
	}
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		s.ErrCh <- err
		return
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		s.ErrCh <- err
		return
	}
	defer ln.Close()

	log.Println("started server")

	for {
		select {
		case <-s.ctx.Done():
			log.Println("stop doing work")
			break
		default:
		}

		// do not wait forever for a new connection
		if err := ln.SetDeadline(time.Now().Add(time.Second * 3)); err != nil {
			s.ErrCh <- err
			return
		}

		time.Sleep(time.Second)
		conn, err := ln.Accept()
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			fmt.Println(err)
			continue
		}
		if err != nil {
			fmt.Println("cannot accept new connection: ", err)
			continue
		}

		fmt.Println("accepted a connection")

		go s.handle(conn)

	}
}

func (s *Server) Shutdown() {
	defer log.Println("shutdown complete")

	s.shutdownFn()
	s.cleanup()
	<-s.shutdownCompleteCh
}

func (s *Server) cleanup() {
	defer time.Sleep(time.Second * 2)

	fmt.Println("server cleanup complete")
	close(s.shutdownCompleteCh)
}

func (s *Server) handle(conn net.Conn) {
	log.Println("handling new connection...")
}
