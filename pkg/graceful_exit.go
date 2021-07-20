package pkg

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func RegisterInterrupt() <-chan os.Signal {
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return sig
}

func ContextWithInterrupt(ctx context.Context, sig <-chan os.Signal) context.Context {
	ctx, release := context.WithCancel(ctx)

	go func() {
		<-sig
		release()
	}()

	return ctx
}
