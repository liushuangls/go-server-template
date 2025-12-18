run:
	APP_MODE=local go run main.go wire_gen.go

generate:
	go generate ./data
	go tool wire gen wire.go

build: generate
	go mod tidy -v
	go build -o=output/server main.go wire_gen.go
