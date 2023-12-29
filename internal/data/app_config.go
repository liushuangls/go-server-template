package data

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/spf13/cast"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/appconfig"
	"github.com/liushuangls/go-server-template/internal/data/ent/predicate"
)

type AppConfigRepo struct {
	*Data
}

func NewAppConfigRepo(data *Data) *AppConfigRepo {
	return &AppConfigRepo{Data: data}
}

func (repo *AppConfigRepo) getAll(ctx context.Context, ps ...predicate.AppConfig) ([]*ent.AppConfig, error) {
	configs, err := repo.db.AppConfig.Query().
		Where(ps...).
		Where(appconfig.DeleteTimeIsNil()).
		All(ctx)
	return configs, repo.warpError(err)
}

func (repo *AppConfigRepo) ConvValue(valueType appconfig.ValueType, value string) (any, error) {
	switch valueType {
	case appconfig.ValueTypeString:
		return value, nil
	case appconfig.ValueTypeInt:
		return cast.ToInt(value), nil
	case appconfig.ValueTypeObject:
		m := map[string]any{}
		if err := json.Unmarshal([]byte(value), &m); err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, errors.New("unSupport type")
}

func (repo *AppConfigRepo) GetAllClientConfig(ctx context.Context, appName string) ([]*ent.AppConfig, error) {
	return repo.getAll(ctx, appconfig.TypeEQ(appconfig.TypeClient), appconfig.AppName(appName))
}

func (repo *AppConfigRepo) GetAllServerConfig(ctx context.Context) ([]*ent.AppConfig, error) {
	return repo.getAll(ctx, appconfig.TypeEQ(appconfig.TypeServer))
}
