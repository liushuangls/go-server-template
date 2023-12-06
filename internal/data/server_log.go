package data

import (
	"context"

	"github.com/samber/lo"

	"github.com/liushuangls/go-server-template/internal/data/ent"
)

type ServerLogRepo struct {
	*Data
}

func NewServerLogRepo(data *Data) *ServerLogRepo {
	return &ServerLogRepo{
		Data: data,
	}
}

func (repo *ServerLogRepo) create(log *ent.ServerLog) *ent.ServerLogCreate {
	log.ErrMsg = lo.Substring(log.ErrMsg, 0, 4000)
	log.RespErrMsg = lo.Substring(log.RespErrMsg, 0, 4000)
	creator := repo.db.ServerLog.Create().
		SetUserID(log.UserID).SetIP(log.IP).SetMethod(log.Method).
		SetPath(log.Path).SetQuery(log.Query).SetErrMsg(log.ErrMsg).
		SetRespErrMsg(log.RespErrMsg).SetBody(log.Body).SetExtra(log.Extra).
		SetCode(log.Code)
	if log.Level != "" {
		creator.SetLevel(log.Level)
	}
	if log.From != "" {
		creator.SetFrom(log.From)
	}
	return creator
}

func (repo *ServerLogRepo) Create(ctx context.Context, log *ent.ServerLog) (*ent.ServerLog, error) {
	return repo.create(log).Save(ctx)
}

func (repo *ServerLogRepo) Creates(ctx context.Context, logs ...*ent.ServerLog) error {
	creators := make([]*ent.ServerLogCreate, 0, len(logs))
	for _, log := range logs {
		creators = append(creators, repo.create(log))
	}
	return repo.db.ServerLog.CreateBulk(creators...).Exec(ctx)
}
