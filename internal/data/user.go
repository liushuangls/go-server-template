package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/predicate"
	userSchema "github.com/liushuangls/go-server-template/internal/data/ent/user"
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
)

type UserRepo struct {
	*Data
	oauth *UserOAuthRepo
}

func NewUserRepo(data *Data, oauth *UserOAuthRepo) *UserRepo {
	return &UserRepo{Data: data, oauth: oauth}
}

func (repo *UserRepo) findOne(ctx context.Context, ps ...predicate.User) (*ent.User, error) {
	user, err := repo.db.User.Query().
		Where(ps...).
		Where(userSchema.DeleteTimeIsNil()).
		First(ctx)
	if err != nil {
		return nil, repo.warpError(err)
	}
	return user, nil
}

func (repo *UserRepo) create(ctx context.Context, db *ent.Client, user *ent.User) (*ent.User, error) {
	u, err := db.User.Create().
		SetNickname(user.Nickname).
		SetRegisterType(user.RegisterType).
		SetRegisterIP(user.RegisterIP).
		SetRegisterRegion(user.RegisterRegion).
		SetEmail(user.Email).
		SetEmailVerified(user.EmailVerified).
		SetPassword(user.Password).
		SetAvatar(user.Avatar).
		SetProfile(user.Profile).
		Save(ctx)
	if err != nil {
		return nil, repo.warpError(err)
	}
	return u, err
}

func (repo *UserRepo) count(ctx context.Context, ps ...predicate.User) (int, error) {
	n, err := repo.db.User.Query().
		Where(ps...).
		Where(userSchema.DeleteTimeIsNil()).
		Count(ctx)
	if err != nil {
		return 0, repo.warpError(err)
	}
	return n, nil
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return repo.findOne(ctx, userSchema.Email(email))
}

func (repo *UserRepo) FindByID(ctx context.Context, uid int) (*ent.User, error) {
	return repo.findOne(ctx, userSchema.ID(uid))
}

func (repo *UserRepo) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	return repo.create(ctx, repo.db, u)
}

func (repo *UserRepo) CreateWithOAuth(ctx context.Context, user *ent.User, oauth *ent.UserOAuth) (*ent.User, error) {
	var (
		u   *ent.User
		err error
	)
	err = withTx(ctx, repo.db, func(tx *ent.Tx) error {
		u, err = repo.create(ctx, tx.Client(), user)
		if err != nil {
			return err
		}
		_, err = repo.oauth.create(ctx, tx.Client(), u.ID, oauth)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, repo.warpError(err)
	}
	return u, nil
}

func (repo *UserRepo) DeleteAccount(ctx context.Context, userID int) error {
	err := withTx(ctx, repo.db, func(tx *ent.Tx) error {
		_, err := tx.UserOAuth.
			Update().
			Where(useroauth.UserID(userID)).
			SetDeleteTime(time.Now()).
			Save(ctx)
		if err != nil {
			return err
		}
		_, err = tx.User.
			Update().
			Where(userSchema.IDEQ(userID)).
			SetDeleteTime(time.Now()).
			Save(ctx)
		return err
	})
	if err != nil {
		return errors.Join(err, fmt.Errorf("UserRepo.DeleteAccount userID:%d", userID))
	}
	return nil
}
