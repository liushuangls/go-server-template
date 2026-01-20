package e2etest

import (
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *APISuite) TestHealth_OK() {
	status, env, data, err := s.api.Health.Get(s.ctx, "ping")
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, status)
	require.Equal(s.T(), 0, env.Code)
	require.Equal(s.T(), "", env.Message)
	require.NotNil(s.T(), data)
	require.Equal(s.T(), "ping", data.Reply)
}

func (s *APISuite) TestHealth_InvalidParams_MissingMessage() {
	status, env, data, err := s.api.Health.Get(s.ctx, "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusBadRequest, status)
	require.Equal(s.T(), 1000, env.Code)
	require.NotEmpty(s.T(), env.Message)
	require.Contains(s.T(), env.Message, "required")
	require.Nil(s.T(), data)
}
