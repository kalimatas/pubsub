package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kalimatas/pubsub"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	server := pubsub.NewServer(ctx)

	if err := server.Start(); err != nil {
		fmt.Println("cannot start server: ", err)
		os.Exit(1)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	cancelFn()
	<-server.Closed
}
