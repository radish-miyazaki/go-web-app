package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/radish-miyazaki/go-web-app/config"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Printf("failed to listen port %d: %v", cfg.Port, err)
		return err
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux, cleanup, err := NewMux(ctx, cfg)
	defer cleanup()
	if err != nil {
		return err
	}

	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %s", err)
		os.Exit(1)
	}
}
