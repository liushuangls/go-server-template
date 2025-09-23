package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

const (
	bearer = "Bearer"
	space  = " "
)

type ClaimsParam struct {
	UserID int `json:"user_id"`
}

type Config struct {
	Secret string `koanf:"Secret"`
	Issuer string `koanf:"Issuer"`
}

type Token struct {
	Token    string
	ExpireAt int64
}

type customClaims struct {
	ClaimsParam
	jwt.RegisteredClaims
}

type JWT struct {
	secret  []byte
	issuer  string
	subject string
}

func NewJWT(conf *Config) (*JWT, error) {
	if conf == nil || conf.Secret == "" || conf.Issuer == "" {
		return nil, errors.New("jwt config error")
	}
	return &JWT{
		issuer: conf.Issuer,
		secret: []byte(conf.Secret),
	}, nil
}

func (j *JWT) GenerateToken(param ClaimsParam, duration time.Duration) (*Token, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)
	claims := customClaims{
		param,
		jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   "",
			Audience:  []string{""},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}
	withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := withClaims.SignedString(j.secret)
	if err != nil {
		return nil, err
	}
	return &Token{
		Token:    fmt.Sprintf("%s%s%s", bearer, space, signedString),
		ExpireAt: expiresAt.Unix(),
	}, nil
}

func (j *JWT) ParseToken(input string) (*ClaimsParam, error) {
	split := strings.Split(input, space)
	if len(split) != 2 || split[0] != bearer {
		return nil, ecode.InvalidToken.WithCause(fmt.Errorf("format error for split"))
	}
	signedString := split[1]

	token, err := jwt.ParseWithClaims(signedString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, ecode.InvalidToken.WithCause(fmt.Errorf("jwt parse err: %s", err))
	}

	claims, ok := token.Claims.(*customClaims)
	if !(ok && token.Valid) {
		return nil, ecode.InvalidToken.WithCause(fmt.Errorf("jwt: token invalid"))
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, ecode.ExpiresToken.WithCause(fmt.Errorf("token expires"))
	}

	return &claims.ClaimsParam, nil
}
