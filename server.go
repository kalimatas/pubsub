package pubsub

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	port = 1042
)

type Server struct {
	Port int

	Wg sync.WaitGroup
}

func NewServer() *Server {
	s := &Server{
		Port: port,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 2)
	s.Wg.Add(1)

	go func() {
		for {
			select {
			case <-ticker.C:
				println("tick")
			case <-ctx.Done():
				s.cleanup()
				return
			}
		}
	}()

	return nil
}

func (s *Server) cleanup() {
	defer s.Wg.Done()

	fmt.Println("server cleanup")
}
