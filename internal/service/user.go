package service

import (
	"context"
	"time"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/dto/response"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/jwt"
	"go.uber.org/zap"
)

type UserService struct {
	log *zap.SugaredLogger
	jwt *jwt.JWT

	userRepo *data.UserRepo
}

func NewUserService(log *zap.SugaredLogger, jwt *jwt.JWT, userRepo *data.UserRepo) *UserService {
	return &UserService{log: log, jwt: jwt, userRepo: userRepo}
}

func (u *UserService) getJwtToken(user *ent.User) (*jwt.Token, error) {
	return u.jwt.GenerateToken(jwt.ClaimsParam{UserID: user.ID}, time.Hour*24*90)
}

func (u *UserService) LoginWithEmail(ctx context.Context, req *request.EmailLoginReq) (*response.UserLoginInfo, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == ecode.NotFound {
			return nil, ecode.EmailOrPasswordErr
		}
		return nil, err
	}

	if user.Password != req.Password {
		return nil, ecode.EmailOrPasswordErr
	}
	token, err := u.getJwtToken(user)
	if err != nil {
		return nil, err
	}
	return &response.UserLoginInfo{
		UserInfo: response.UserInfo{
			ID:    user.ID,
			Email: user.Email,
		},
		AccessToken: token,
	}, nil
}
