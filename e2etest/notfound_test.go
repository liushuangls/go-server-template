package e2etest

import (
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *APISuite) TestNotFound_UnknownRoute() {
	status, body, err := s.c.Get(s.ctx, "/__not_found__", nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusNotFound, status)

	env, data, err := decodeEnvelope[map[string]any](body)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1000, env.Code)
	require.NotEmpty(s.T(), env.Message)
	require.Nil(s.T(), data)
}
