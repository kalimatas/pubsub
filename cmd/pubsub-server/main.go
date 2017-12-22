package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/kalimatas/pubsub"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	server := pubsub.NewServer()
	if err := server.Start(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "cannot start Server: ", err)
		os.Exit(1)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		cancel()
		server.Wg.Wait()
	}

	println("finish him")
}
