package xoauth2

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type Client interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	ParseIDToken(ctx context.Context, rawIDToken string) (*IdToken, error)
}

type Config struct {
	ClientID    string `yaml:"ClientID"`
	Secret      string `yaml:"Secret"`
	RedirectURL string `yaml:"RedirectURL"`
	UseSecret   bool   `yaml:"UseSecret"`
	ClientType  string `yaml:"ClientType"`
	TenantID    string `yaml:"TenantID"`
	TeamID      string `yaml:"TeamID"`
	KeyID       string `yaml:"KeyID"`
}

// IdToken is an OpenID Connect extension that provides a predictable representation
// of an authorization event.
//
// The ClientID Token only holds fields OpenID Connect requires. To access additional
// claims returned by the server, use the Claims method.
type IdToken struct {
	// The URL of the server which issued this token. OpenID Connect
	// requires this value always be identical to the URL used for
	// initial discovery.
	//
	// Note: Because of a known issue with Google Accounts' implementation
	// this value may differ when using Google.
	//
	// See: https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo
	Issuer string

	// The client ClientID, or set of client IDs, that this token is issued for. For
	// common uses, this is the client that initialized the auth flow.
	//
	// This package ensures the audience contains an expected value.
	Audience []string

	// A unique string which identifies the end user.
	Subject string

	// Expiry of the token. Ths package will not process tokens that have
	// expired unless that validation is explicitly turned off.
	Expiry time.Time
	// When the token was issued by the provider.
	IssuedAt time.Time

	// Initial nonce provided during the authentication redirect.
	//
	// This package does NOT provided verification on the value of this field
	// and it's the user's responsibility to ensure it contains a valid value.
	Nonce string

	// at_hash claim, if set in the ClientID token. Callers can verify an access token
	// that corresponds to the ClientID token using the VerifyAccessToken method.
	AccessTokenHash string

	OpenID         string `json:"open_id"`
	UnionID        string `json:"union_id"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	IsPrivateEmail bool   `json:"is_private_email"`
	// https://tools.ietf.org/html/bcp47
	Locale  string `json:"locale"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Profile string `json:"profile"`
}
