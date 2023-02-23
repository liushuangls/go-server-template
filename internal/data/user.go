package data

import (
	"context"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/predicate"
	userSchema "github.com/liushuangls/go-server-template/internal/data/ent/user"
	"github.com/liushuangls/go-server-template/pkg/ecode"
)

type UserRepo struct {
	*Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{data}
}

func (u *UserRepo) findOne(ctx context.Context, wheres ...predicate.User) (*ent.User, error) {
	user, err := u.db.User.Query().
		Where(wheres...).
		Where(userSchema.DeleteTimeIsNil()).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ecode.NotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return u.findOne(ctx, userSchema.Email(email))
}

func (u *UserRepo) FindByID(ctx context.Context, uid int) (*ent.User, error) {
	return u.findOne(ctx, userSchema.ID(uid))
}
