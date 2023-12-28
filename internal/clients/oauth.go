package clients

import (
	"context"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/xoauth2"
)

type OauthClients struct {
	Google    map[string]xoauth2.Client
	Microsoft map[string]xoauth2.Client
	Apple     map[string]xoauth2.Client
}

func NewOauthClients(ctx context.Context, conf *configs.Config) (*OauthClients, error) {
	clients := &OauthClients{
		Google:    make(map[string]xoauth2.Client),
		Microsoft: make(map[string]xoauth2.Client),
		Apple:     make(map[string]xoauth2.Client),
	}

	for _, c := range conf.OAuth2.Google {
		cli, err := xoauth2.NewGoogleClient(ctx, c)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
		clients.Google[c.ClientType] = cli
	}

	for _, c := range conf.OAuth2.Microsoft {
		cli, err := xoauth2.NewMicrosoftClient(ctx, c)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
		clients.Microsoft[c.ClientType] = cli
	}

	for _, c := range conf.OAuth2.Apple {
		cli, err := xoauth2.NewAppleClient(ctx, c)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
		clients.Apple[c.ClientType] = cli
	}

	return clients, nil
}

func (clients *OauthClients) GetClient(platform useroauth.Platform, clientType string) (c xoauth2.Client,
	ok bool) {
	switch platform {
	case useroauth.PlatformGoogle:
		c, ok = clients.Google[clientType]
	case useroauth.PlatformApple:
		c, ok = clients.Apple[clientType]
	case useroauth.PlatformMicrosoft:
		c, ok = clients.Microsoft[clientType], true
	default:
		return nil, false
	}
	return
}
