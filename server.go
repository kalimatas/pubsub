package pubsub

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	defaultAddress = ":1042"
)

type Server struct {
	address string

	ln net.Listener

	ctx    context.Context
	Closed chan struct{}
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		address: defaultAddress,

		ctx:    ctx,
		Closed: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.ln = ln

	go func() {
	loop:
		for {
			select {
			case <-s.ctx.Done():
				break loop
			default:
			}

			// do not wait forever for a new connection
			if err := ln.SetDeadline(time.Now().Add(time.Second * 3)); err != nil {
				continue
			}

			conn, err := ln.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("cannot accept new connection: ", err)
				continue
			}

			go s.handle(conn)

		}

		s.shutdown()
	}()

	return nil
}

func (s *Server) shutdown() {
	s.ln.Close()
	close(s.Closed)
}

func (s *Server) handle(conn net.Conn) {
	log.Println("handling new connection...")
}
