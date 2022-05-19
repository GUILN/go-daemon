package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultTick = 60 * time.Second
const defaultConfigFileName = "config"

type Config struct {
	contentType string
	server      string
	statusCode  int
	tick        time.Duration
	url         string
	userAgent   string
}

func (c *Config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.String(defaultConfigFileName, "", "Path to Config file")

	var (
		statusCode  = flags.Int("status", 200, "Http Response Status Code")
		tick        = flags.Duration("tick", defaultTick, "Ticking interval")
		server      = flags.String("server", "", "Server http header value")
		contentType = flags.String("content_type", "", "Content-Type Http header value")
		userAgent   = flags.String("user_agent", "", "User-Agent Http header value")
		url         = flags.String("url", "", "Request url")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.statusCode = *statusCode
	c.tick = *tick
	c.server = *server
	c.contentType = *contentType
	c.userAgent = *userAgent
	c.url = *url

	return nil
}

func (c *Config) Reload() error {
	if err := c.init(os.Args); err != nil {
		return err
	}

	return nil
}

func Run(ctx context.Context, c *Config, out io.Writer) error {
	if err := c.Reload(); err != nil {
		return err
	}
	// 1 - Log to stadandard output
	// in this case we are injecting the output writer in order to make the test easier
	log.SetOutput(out)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(c.tick):
			resp, err := http.Get(c.url)
			if err != nil {
				return err
			}

			if resp.StatusCode != c.statusCode {
				log.Printf("Status code mismatch, got %d\n", resp.StatusCode)
			}

			if s := resp.Header.Get("server"); s != c.server {
				log.Printf("Server mismatch, got %s\n", s)
			}

			if ct := resp.Header.Get("content-type"); ct != c.contentType {
				log.Printf("Content-Type mismatch, got %s\n", ct)
			}

			if ua := resp.Header.Get("user-agent"); ua != c.userAgent {
				log.Printf("User-Agent mismatch, got %s\n", ua)
			}
		}
	}
}
