package xoauth2

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// https://learn.microsoft.com/zh-cn/azure/active-directory/develop/id-tokens
type microsoftIdTokenClaims struct {
	OID   string `json:"oid"` // 不同应用相同
	SubID string `json:"sub"` // 不同应用不同
	TID   string `json:"tid"` // 表示用户登录到的租户
	Ver   string `json:"ver"` // 指示 id_token 的版本
	Email string `json:"email"`
	Name  string `json:"name"`
}

type microsoftClient struct {
	oauth2.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewMicrosoftClient(ctx context.Context, config Config) (Client, error) {
	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", config.TenantID))
	if err != nil {
		return nil, err
	}
	return &microsoftClient{
		Config: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.Secret,
			Endpoint:     microsoft.AzureADEndpoint(""),
			RedirectURL:  config.RedirectURL,
			Scopes: []string{
				oidc.ScopeOpenID,
				"email",
				"profile",
			},
		},
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: config.ClientID}),
	}, nil
}

func (m *microsoftClient) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	opts = append(opts, oauth2.SetAuthURLParam("include_granted_scopes", "true"))
	return m.Config.AuthCodeURL(state, opts...)
}

func (m *microsoftClient) ParseIDToken(ctx context.Context, rawIDToken string) (*IdToken, error) {
	idToken, err := m.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	var claims microsoftIdTokenClaims
	err = idToken.Claims(&claims)
	if err != nil {
		return nil, err
	}
	return &IdToken{
		Issuer:          idToken.Issuer,
		Audience:        idToken.Audience,
		Subject:         idToken.Subject,
		Expiry:          idToken.Expiry,
		IssuedAt:        idToken.IssuedAt,
		Nonce:           idToken.Nonce,
		AccessTokenHash: idToken.AccessTokenHash,
		OpenID:          claims.OID,
		UnionID:         claims.SubID,
		Email:           claims.Email,
		Name:            claims.Name,
	}, nil
}
