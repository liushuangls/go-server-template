package xoauth2

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

// https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api/authenticating_users_with_sign_in_with_apple#3383773
type appleIdTokenClaims struct {
	Email          string `json:"email"`
	EmailVerified  any    `json:"email_verified"`
	IsPrivateEmail any    `json:"is_private_email"`
	//RealUserStatus string `json:"real_user_status"`
	TransferSub string `json:"transfer_sub"`
}

// https://developer.apple.com/documentation/sign_in_with_apple/implementing_user_authentication_with_sign_in_with_apple
type appleClient struct {
	oauth2.Config
	provider        *oidc.Provider
	verifier        *oidc.IDTokenVerifier
	secretExpiredAt time.Time
	conf            Config
}

func NewAppleClient(ctx context.Context, config Config) (Client, error) {
	provider, err := oidc.NewProvider(ctx, "https://appleid.apple.com")
	if err != nil {
		return nil, err
	}
	secretExpiredAt := time.Now().AddDate(0, 0, 180)
	secret, err := GenerateAppleClientSecret(config.Secret, config.TeamID, config.ClientID, config.KeyID, secretExpiredAt)
	if err != nil {
		return nil, err
	}
	return &appleClient{
		Config: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://appleid.apple.com/auth/authorize",
				TokenURL: "https://appleid.apple.com/auth/token",
			},
			RedirectURL: config.RedirectURL,
			Scopes: []string{
				"name",
				"email",
			},
		},
		provider:        provider,
		verifier:        provider.Verifier(&oidc.Config{ClientID: config.ClientID}),
		secretExpiredAt: secretExpiredAt,
		conf:            config,
	}, nil
}

func (g *appleClient) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	opts = append(opts,
		oauth2.SetAuthURLParam("response_mode", "form_post"),
	)
	return g.Config.AuthCodeURL(state, opts...)
}

func (g *appleClient) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	// 签名过期,重新生成
	if g.secretExpiredAt.Before(time.Now()) {
		expiredAt := time.Now().AddDate(0, 0, 180)
		secret, err := GenerateAppleClientSecret(g.conf.Secret, g.conf.TeamID, g.conf.ClientID, g.conf.KeyID, expiredAt)
		if err != nil {
			return nil, err
		}
		g.secretExpiredAt = expiredAt
		g.Config.ClientSecret = secret
	}
	return g.Config.Exchange(ctx, code, opts...)
}

func (g *appleClient) ParseIDToken(ctx context.Context, rawIDToken string) (*IdToken, error) {
	idToken, err := g.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	var claims appleIdTokenClaims
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
		EmailVerified:   cast.ToBool(claims.EmailVerified),
		IsPrivateEmail:  cast.ToBool(claims.IsPrivateEmail),
	}, nil
}

func GenerateAppleClientSecret(signingKey, teamID, clientID, keyID string, expiredAt time.Time) (string, error) {
	block, _ := pem.Decode([]byte(signingKey))
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Create the Claims
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    teamID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		Subject:   clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = keyID

	return token.SignedString(privateKey)
}
