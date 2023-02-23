package routes

import (
	"github.com/google/wire"
	v1 "github.com/liushuangls/go-server-template/internal/routes/v1"
)

var ProviderSet = wire.NewSet(NewEngine, NewHttpEngine, v1.NewUserRoute)
