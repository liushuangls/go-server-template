package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	server, cleanup, err := app(ctx)
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()
	if err != nil {
		panic(err)
	}

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
