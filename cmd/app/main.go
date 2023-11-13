package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	server, cleanup, err := app(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
