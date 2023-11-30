package clients

import (
	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/pkg/hashidv2"
)

type HashID struct {
	User *hashidv2.HashID
}

func NewHashID(conf *configs.Config) *HashID {
	return &HashID{
		User: hashidv2.New(&conf.HashID.User),
	}
}

func (h *HashID) GetHashID(t string) *hashidv2.HashID {
	switch t {
	case "User":
		return h.User
	}
	return nil
}
