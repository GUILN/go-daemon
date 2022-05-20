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

	// 2 - shutdown on SIGTERM/SIGINT | 3 - Reload config on SIGHUP
	// sinalChan and go routine below are used to gracefully shutdown the daemon in case
	// of signals: SIGTERM, SIGINT
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Printf("Got SIGINT/SIGTERM, exiting...")
					cancel()
					os.Exit(1)
					return
				case os.Interrupt:
					log.Printf("Got Interrupt, exiting...")
					cancel()
					os.Exit(1)
					return
				case syscall.SIGHUP:
					log.Printf("Got SIGHUP, reloading...")
					if err := c.Reload(); err != nil {
						fmt.Fprintf(os.Stderr, "%s\n", err)
						cancel()
						os.Exit(1)
						return
					}
				}
			case <-ctx.Done():
				log.Printf("Done")
				os.Exit(1)
				return
			}

		}
	}()

	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
		cancel()
	}()

	// 1 - Log to stadandard output
	// out parameter (last parameter is logging into standard output)
	// the dependency inversion is to make easier to test
	if err := app.Run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
