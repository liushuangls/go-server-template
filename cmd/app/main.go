package main

func main() {
	server, cleanup, err := app()
	if err != nil {
		panic(err)
	}
	err = server.Run()
	if err != nil {
		panic(err)
	}
	defer cleanup()
}
