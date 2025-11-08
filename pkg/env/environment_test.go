package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunEnvIsTest(t *testing.T) {
	require.Equal(t, EnvTest, runEnv)
}

func TestBuildEnvIsDevelopment(t *testing.T) {
	require.Equal(t, EnvDevelopment, buildEnv)
}
