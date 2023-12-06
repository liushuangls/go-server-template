package clients

import (
	"github.com/google/wire"
	"github.com/ip2location/ip2location-go/v9"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/clients/logdbsync"
	"github.com/liushuangls/go-server-template/internal/clients/publicoss"
)

var ProviderSet = wire.NewSet(
	NewOauthClients,
	publicoss.NewAvatar,
	NewIp2Location,
	NewHashID,
	logdbsync.NewClient,
)

func NewIp2Location(conf *configs.Config) (*ip2location.DB, error) {
	return ip2location.OpenDB(conf.App.IP2Location)
}
