# Introduction

An easy to use, easy to maintenance template.

# Use

```
// install
go install golang.org/x/tools/cmd/gonew@latest

gonew github.com/liushuangls/go-server-template your.domain/myprog

// run
cd myprog

cp configs/test.yaml configs/local.yaml

make install

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
