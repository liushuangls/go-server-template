package data

import (
	"context"

	entSchema "github.com/liushuangls/go-server-template/internal/data/ent"
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

func (u *UserRepo) findOne(ctx context.Context, wheres ...predicate.User) (*entSchema.User, error) {
	user, err := u.db.User.Query().
		Where(wheres...).
		Where(userSchema.DeleteTimeIsNil()).
		First(ctx)
	if err != nil {
		if entSchema.IsNotFound(err) {
			return nil, ecode.NotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*entSchema.User, error) {
	return u.findOne(ctx, userSchema.Email(email))
}

func (u *UserRepo) FindByID(ctx context.Context, uid int) (*entSchema.User, error) {
	return u.findOne(ctx, userSchema.ID(uid))
}
