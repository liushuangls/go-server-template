package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	userSchema "github.com/liushuangls/go-server-template/internal/data/ent/user"
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/dto/response"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/jwt"
	"github.com/liushuangls/go-server-template/pkg/xoauth2"
)

type UserService struct {
	*Options
}

func NewUserService(opt *Options) *UserService {
	return &UserService{opt}
}

func (u *UserService) UserInfo(ctx context.Context, user *ent.User) (*response.UserInfo, error) {
	return u.getUserInfo(user), nil
}

func (u *UserService) getAndSaveOAuthState(ctx context.Context) string {
	nonce := uuid.NewString()
	key := fmt.Sprintf("oauth2:%s", nonce)
	u.Redis.Set(ctx, key, 1, time.Minute*5)
	return nonce
}

func (u *UserService) GetOAuthCodeURL(ctx context.Context, req *request.OAuthCodeURLReq) (*response.OAuthCodeURLResp, error) {
	oauth, ok := u.OAuthClients.GetClient(useroauth.Platform(req.Platform))
	if !ok {
		return nil, ecode.InvalidParams
	}
	nonce := u.getAndSaveOAuthState(ctx)
	return &response.OAuthCodeURLResp{Url: oauth.AuthCodeURL(nonce)}, nil
}

func (u *UserService) verifyAndDelOAuthState(ctx context.Context, nonce string) (bool, error) {
	key := fmt.Sprintf("oauth2:%s", nonce)
	_, err := u.Redis.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	u.Redis.Del(ctx, key)
	return true, nil
}

func (u *UserService) OAuthCallback(ctx context.Context, ipInfo *request.IPInfo, req *request.OAuthCallbackReq) (*response.UserLoginInfo, error) {
	oauth, ok := u.OAuthClients.GetClient(useroauth.Platform(req.Platform))
	if !ok {
		return nil, ecode.InvalidParams
	}

	// 防止重放攻击
	pass, err := u.verifyAndDelOAuthState(ctx, req.State)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, ecode.InvalidOAuthState
	}

	token, err := oauth.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	idToken, err := oauth.ParseIDToken(ctx, cast.ToString(token.Extra("id_token")))
	if err != nil {
		return nil, err
	}

	loginInfo, err := u.oauthLoginOrRegister(ctx, idToken, ipInfo, req.Platform)
	if err != nil {
		return nil, err
	}
	return loginInfo, nil
}

func (u *UserService) oauthLoginOrRegister(ctx context.Context, idToken *xoauth2.IdToken, ipInfo *request.IPInfo,
	platform string) (*response.UserLoginInfo, error) {
	var (
		user         *ent.User
		userPlatform = useroauth.Platform(platform)
		openID       = idToken.OpenID
	)
	oauthInfo, err := u.UserOauthRepo.FindByOpenID(ctx, userPlatform, openID)
	if err != nil {
		// 注册
		if errors.Is(err, ecode.NotFound) {
			user, err = u.registerByOAuth(ctx, idToken, ipInfo, platform)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// 登录
	if user == nil {
		user, err = u.UserRepo.FindByID(ctx, oauthInfo.UserID)
		if err != nil {
			return nil, err
		}
	}
	return u.getLoginInfo(user)
}

func (u *UserService) generateToken(user *ent.User) (*jwt.Token, error) {
	claims := jwt.ClaimsParam{
		UserID: user.ID,
	}
	token, err := u.Jwt.GenerateToken(claims, time.Hour*24*180)
	if err != nil {
		return nil, err
	}
	return &jwt.Token{
		Token:    token.Token,
		ExpireAt: token.ExpireAt,
	}, nil
}

func (u *UserService) getLoginInfo(user *ent.User) (*response.UserLoginInfo, error) {
	token, err := u.generateToken(user)
	if err != nil {
		return nil, err
	}
	userInfo := u.getUserInfo(user)
	return &response.UserLoginInfo{
		UserInfo:    *userInfo,
		AccessToken: token,
	}, nil
}

func (u *UserService) getUserInfo(user *ent.User) *response.UserInfo {
	return &response.UserInfo{
		ID:         user.ID,
		Email:      user.Email,
		Avatar:     user.Avatar,
		NickName:   user.Nickname,
		RegisterAt: user.CreateTime.Unix(),
	}
}

func (u *UserService) registerByOAuth(ctx context.Context, idToken *xoauth2.IdToken,
	ipInfo *request.IPInfo, platform string) (user *ent.User, err error) {
	userOAuth := &ent.UserOAuth{
		Platform: useroauth.Platform(platform),
		OpenID:   idToken.OpenID,
		UnionID:  idToken.UnionID,
	}

	// 如果不同oauth平台的email相同，则认为是同一个用户
	if idToken.Email != "" {
		user, err := u.UserRepo.FindByEmail(ctx, idToken.Email)
		if err != nil && !errors.Is(err, ecode.NotFound) {
			return nil, err
		}
		if user != nil {
			_, err := u.UserOauthRepo.CreateByUser(ctx, user.ID, userOAuth)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
	}

	idToken.Picture = u.Avatar.Handle(idToken.Picture)
	// 创建新账号
	user = &ent.User{
		Nickname:       idToken.Name,
		RegisterType:   userSchema.RegisterTypeOauth2,
		RegisterIP:     ipInfo.IP,
		RegisterRegion: ipInfo.CountryShort,
		Email:          idToken.Email,
		EmailVerified:  idToken.EmailVerified,
		Avatar:         idToken.Picture,
		Profile:        idToken.Profile,
	}
	user, err = u.UserRepo.CreateWithOAuth(ctx, user, userOAuth)
	if err != nil {
		return nil, err
	}
	return user, nil
}
