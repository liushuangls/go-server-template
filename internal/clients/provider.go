package clients

import (
	"github.com/google/wire"
	"github.com/oschwald/maxminddb-golang"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/clients/logdbsync"
	"github.com/liushuangls/go-server-template/internal/clients/publicoss"
	"github.com/liushuangls/go-server-template/internal/clients/serverconf"
	"github.com/liushuangls/go-server-template/pkg/ecode"
)

var ProviderSet = wire.NewSet(
	NewOauthClients,
	publicoss.NewAvatar,
	NewGeoLite2,
	NewHashID,
	logdbsync.NewClient,
	serverconf.NewServerConf,
)

func NewGeoLite2(conf *configs.Config) (*maxminddb.Reader, error) {
	r, err := maxminddb.Open(conf.App.GeoLite2DB)
	if err != nil {
		return nil, ecode.WithCaller(err)
	}
	return r, nil
}
