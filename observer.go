package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guiln/go-daemon/app"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := &app.Config{}

	// sinalChan and go routine below are used to gracefully shutdown the daemon in case
	// of signals: SIGTERM, SIGINT
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			log.Printf("Got SIGINT/SIGTERM, exiting...")
			cancel()
			os.Exit(1)
		case <-ctx.Done():
			log.Printf("Done")
			os.Exit(1)
		}
	}()

	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
		cancel()
	}()

	if err := app.Run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
