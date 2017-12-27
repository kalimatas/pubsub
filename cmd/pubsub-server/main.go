package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kalimatas/pubsub"
)

func handleSignals(server *pubsub.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println("got signal")
	server.Shutdown()
}

func main() {
	server := pubsub.NewServer()
	go server.Start()

	go func() {
		err := <-server.ErrCh
		fmt.Println("cannot start server: ", err)
		os.Exit(1)
	}()

	handleSignals(server)

	println("finish him")
}
