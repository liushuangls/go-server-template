install:
	go install github.com/google/wire/cmd/wire@latest
	go install entgo.io/ent/cmd/ent@latest
	go mod download

run:
	APP_MODE=local go run main.go wire_gen.go

generate:
	go generate ./data
	wire gen wire.go

build: generate
	go mod tidy -v
	go build -o=output/server cmd/app/main.go cmd/app/wire_gen.go
