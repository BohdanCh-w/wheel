package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func OSInterruptContext(c context.Context) context.Context {
	ctx, cancel := context.WithCancel(c)

	go func() {
		defer cancel()

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		defer signal.Stop(shutdown)

		select {
		case <-shutdown:
		case <-ctx.Done():
		}
	}()

	return ctx
}
