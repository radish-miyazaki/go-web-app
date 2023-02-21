package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// シグナルを受け取ったら、グレースシャットダウンを開始する
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 別ゴルーチンでサーバを起動
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを示すので異常ではない
		if err := s.srv.Serve(s.l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を受け取る
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	return eg.Wait()
}
