package main

import (
	"OnlineMusic/pkg/app"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

// @title OnlineMusic
// @version 0.0.1
// @description API Server for online music library
func main() {
	a := app.New()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	go func() {
		if err := a.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	if err := a.Stop(context.TODO()); err != nil {
		log.Fatalf("app shutdown returned an err: %v\n", err)
	}
}
