package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type CloserFunc func() error

func Start(router http.Handler, port string) CloserFunc {
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	return func() error {
		ctx, release := context.WithTimeout(context.Background(), time.Second*5)
		defer release()
		return server.Shutdown(ctx)
	}
}
