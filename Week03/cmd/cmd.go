package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schwarzeni/Go-000/Week03/pkg/errgroup"
)

const serverShutdownDuration = time.Second

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// register to system signal
	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigs:
			log.Println("catch system term signal")
			g.StopAll()
		case <-ctx.Done():
		}
		return nil
	})

	g.Go(func() error { return newServer(ctx, ":9000", g.StopAll) })
	g.Go(func() error { return newServer(ctx, ":9001", g.StopAll) })
	g.Go(func() error { return newServer(ctx, ":9002", g.StopAll) })

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// newServer 启动一个新的服务
func newServer(ctx context.Context, addr string, afterShutdownFn func()) error {
	s := &http.Server{Addr: addr}
	s.RegisterOnShutdown(afterShutdownFn)
	log.Println(addr + " server is starting")

	go func() {
		<-ctx.Done()
		log.Println(addr + " server is shutting down")
		shutdownCtx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(serverShutdownDuration))
		defer func() {
			log.Println(addr + " server shuts down")
			cancelFunc()
		}()
		_ = s.Shutdown(shutdownCtx)
	}()

	err := s.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			err = nil
		} else {
			log.Println(addr+" server started failed", err)
		}
	}
	return err
}
