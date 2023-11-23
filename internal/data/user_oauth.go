package data

import (
	"context"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/predicate"
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
)

type UserOAuthRepo struct {
	*Data
}

func NewUserOAuthRepo(data *Data) *UserOAuthRepo {
	return &UserOAuthRepo{Data: data}
}

func (repo *UserOAuthRepo) create(ctx context.Context, db *ent.Client, uid int, oauth *ent.UserOAuth) (*ent.UserOAuth, error) {
	o, err := db.UserOAuth.Create().
		SetPlatform(oauth.Platform).SetOpenID(oauth.OpenID).
		SetUnionID(oauth.UnionID).SetUserID(uid).
		Save(ctx)
	if err != nil {
		return nil, repo.warpError(err)
	}
	return o, nil
}

func (repo *UserOAuthRepo) findOne(ctx context.Context, ps ...predicate.UserOAuth) (*ent.UserOAuth, error) {
	o, err := repo.db.UserOAuth.Query().
		Where(ps...).
		Where(useroauth.DeleteTimeIsNil()).
		First(ctx)
	if err != nil {
		return nil, repo.warpError(err)
	}
	return o, nil
}

func (repo *UserOAuthRepo) FindByOpenID(ctx context.Context, platform useroauth.Platform,
	openID string) (*ent.UserOAuth, error) {
	return repo.findOne(ctx, useroauth.PlatformEQ(platform), useroauth.OpenID(openID))
}

func (repo *UserOAuthRepo) CreateByUser(ctx context.Context, uid int, oauth *ent.UserOAuth) (*ent.UserOAuth, error) {
	return repo.create(ctx, repo.db, uid, oauth)
}
