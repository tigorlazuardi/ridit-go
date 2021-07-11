package pkg

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func ContextWithCtrlC(ctx context.Context) context.Context {
	ctx, stop := context.WithCancel(ctx)

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-sig
		stop()
	}()

	return ctx
}
