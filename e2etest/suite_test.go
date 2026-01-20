package e2etest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite
	ctx context.Context
	c   *Client
	api *API
}

func (s *APISuite) SetupSuite() {
	s.ctx = context.Background()
	s.c = NewClientFromEnv()
	s.api = NewAPI(s.c)

	status, env, data, err := s.api.Health.Get(s.ctx, "ping")
	if err != nil {
		if s.c.explicitBaseURL {
			require.NoError(s.T(), err, "server not reachable; check E2E_BASE_URL or start server")
		}
		s.T().Skipf("server not reachable (%v); start server or set E2E_BASE_URL", err)
	}
	require.Equal(s.T(), 200, status)
	require.Equal(s.T(), 0, env.Code)
	require.NotNil(s.T(), data)
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(APISuite))
}
