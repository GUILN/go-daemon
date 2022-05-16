package main

import (
	"context"
	"fmt"
	"os"

	"github.com/guiln/go-daemon/app"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := &app.Config{}

	defer func() {
		cancel()
	}()

	if err := app.Run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
