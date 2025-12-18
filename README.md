# Introduction

An easy to use, easy to maintenance template.

This template layout reference: https://go.dev/doc/modules/layout

# Use

```
// install
go install golang.org/x/tools/cmd/gonew@latest

gonew github.com/liushuangls/go-server-template your.domain/myprog

// run
cd myprog

cp configs/example.yaml configs/local.yaml

go mod download

make generate

make run
```

# Features

- wire - injects dependencies
- ent - database orm
- gin - router
- viper - config manager
- and so on...
    - jwt
    - ...
