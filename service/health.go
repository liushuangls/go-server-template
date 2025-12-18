package service

import (
	"context"

	"github.com/liushuangls/go-server-template/dto/request"
	"github.com/liushuangls/go-server-template/dto/response"
)

type HealthService struct {
	Options
}

func NewHealthService(opt Options) *HealthService {
	return &HealthService{opt}
}

func (u *HealthService) Health(ctx context.Context, req *request.HealthReq) (*response.HealthResp, error) {
	return &response.HealthResp{
		Reply: req.Message,
	}, nil
}
