package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"tmp/pragmatic-cases/opentelemetry/fib"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := fib.NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
