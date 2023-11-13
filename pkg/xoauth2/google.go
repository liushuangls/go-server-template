package xoauth2

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleIdTokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	// https://tools.ietf.org/html/bcp47
	Locale  string `json:"locale"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Profile string `json:"profile"`
}

// https://developers.google.com/identity/openid-connect/openid-connect?hl=zh-cn#validatinganidtoken
type googleClient struct {
	oauth2.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewGoogleClient(ctx context.Context, config Config) (Client, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	return &googleClient{
		Config: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.Secret,
			Endpoint:     google.Endpoint,
			RedirectURL:  config.RedirectURL,
			Scopes: []string{
				oidc.ScopeOpenID,
				"profile",
				"email",
			},
		},
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: config.ClientID}),
	}, nil
}

func (g *googleClient) ParseIDToken(ctx context.Context, rawIDToken string) (*IdToken, error) {
	idToken, err := g.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	var claims googleIdTokenClaims
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
		OpenID:          idToken.Subject,
		UnionID:         idToken.Subject,
		Email:           claims.Email,
		EmailVerified:   claims.EmailVerified,
		Locale:          claims.Locale,
		Name:            claims.Name,
		Picture:         claims.Picture,
		Profile:         claims.Profile,
	}, nil
}
