package serverconf

import (
	"context"
	"log/slog"
	"sync"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/xjson"
	"github.com/liushuangls/go-server-template/pkg/xslog"
)

type ServerConf struct {
	repo *data.AppConfigRepo
	m    sync.Map
}

func NewServerConf(ctx context.Context, repo *data.AppConfigRepo) (*ServerConf, error) {
	s := &ServerConf{repo: repo, m: sync.Map{}}
	err := s.LoadConf(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *ServerConf) LoadConf(ctx context.Context) error {
	configs, err := s.repo.GetAllServerConfig(ctx)
	if err != nil {
		return err
	}
	for _, conf := range configs {
		s.m.Store(conf.Key, conf)
	}
	return nil
}

func (s *ServerConf) Get(key string, payload any) error {
	v, ok := s.m.Load(key)
	if !ok {
		slog.Error("serverConf.Get", slog.String("key", key), xslog.Error(ecode.NotFound))
		return ecode.NotFound
	}
	val := v.(*ent.AppConfig)
	err := xjson.UnmarshalString(val.Value, payload)
	if err != nil {
		slog.Error("serverConf.Get", slog.String("key", key), xslog.Error(err))
		return err
	}
	return nil
}

func (s *ServerConf) MustGet(key string, payload any) {
	v, ok := s.m.Load(key)
	if !ok {
		slog.Error("serverConf.Get", slog.String("key", key), xslog.Error(ecode.NotFound))
		panic(ecode.NotFound)
	}
	val := v.(*ent.AppConfig)
	err := xjson.UnmarshalString(val.Value, payload)
	if err != nil {
		slog.Error("serverConf.Get", slog.String("key", key), xslog.Error(err))
		panic(err)
	}
}
