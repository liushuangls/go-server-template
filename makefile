install:
	go install github.com/google/wire/cmd/wire@latest
	go install entgo.io/ent/cmd/ent@latest

run:
	APP_MODE=debug go run -tags=jsoniter cmd/app/main.go cmd/app/wire_gen.go

generate:
	go generate ./internal/data/ent
	wire gen cmd/app/wire.go

build: generate
	go mod tidy -v
	go build -tags=jsoniter -o=output/server cmd/app/main.go cmd/app/wire_gen.go
