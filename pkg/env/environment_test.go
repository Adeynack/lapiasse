package env_test

import (
	"testing"

	"adeynack.net/lapiasse/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestRunEnvIsTest(t *testing.T) {
	require.Equal(t, env.EnvTest, env.RunEnv)
}

func TestBuildEnvIsDevelopment(t *testing.T) {
	require.Equal(t, env.EnvDevelopment, env.BuildEnv)
}
