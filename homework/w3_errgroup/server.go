package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello"))
	})

	stop := make(chan struct{})
	serveMux.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		stop <- struct{}{}
	})

	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	group.Go(func() error {
		return server.ListenAndServe()
	})

	group.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("errgroup ctx done")
		case <-stop:
			fmt.Println("server stop")
		}
		timeoutCtx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		return server.Shutdown(timeoutCtx)
	})

	group.Go(func() error {
		quitSignal := make(chan os.Signal, 0)
		signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			fmt.Println("errgroup ctx done")
			return ctx.Err()
		case <-quitSignal:
			fmt.Println("signal quit")
			return errors.New("signal quit")
		}
	})

	_ = group.Wait()
}
