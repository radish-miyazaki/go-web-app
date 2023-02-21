package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radish-miyazaki/go-web-app/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	// シグナルを受け取ったら、グレースシャットダウンを開始する
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// FIXME: コマンドラインで一時的に確認するため
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	// 別ゴルーチンでサーバを起動
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを示すので異常ではない
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を受け取る
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %s", err)
		os.Exit(1)
	}
}
