package clients

import (
	"context"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/xoauth2"
)

type OauthClients struct {
	Google    xoauth2.Client
	Microsoft xoauth2.Client
	Apple     xoauth2.Client
}

func NewOauthClients(ctx context.Context, conf *configs.Config) (*OauthClients, error) {
	var (
		err     error
		clients = &OauthClients{}
	)
	if conf.OAuth2.Google.ClientID != "" {
		clients.Google, err = xoauth2.NewGoogleClient(ctx, conf.OAuth2.Google)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
	}
	if conf.OAuth2.Microsoft.ClientID != "" {
		clients.Microsoft, err = xoauth2.NewGoogleClient(ctx, conf.OAuth2.Microsoft)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
	}
	if conf.OAuth2.Apple.ClientID != "" {
		clients.Apple, err = xoauth2.NewGoogleClient(ctx, conf.OAuth2.Apple)
		if err != nil {
			return nil, ecode.WithCaller(err)
		}
	}
	return clients, nil
}

func (clients *OauthClients) GetClient(platform useroauth.Platform) (xoauth2.Client, bool) {
	switch platform {
	case useroauth.PlatformGoogle:
		return clients.Google, true
	case useroauth.PlatformApple:
		return clients.Microsoft, true
	case useroauth.PlatformMicrosoft:
		return clients.Microsoft, true
	default:
		return nil, false
	}
}
