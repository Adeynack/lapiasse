package app_test

import (
	"testing"

	"adeynack.net/lapiasse/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestRunEnvIsTest(t *testing.T) {
	require.Equal(t, app.EnvTest, app.RunEnv)
}

func TestBuildEnvIsDevelopment(t *testing.T) {
	require.Equal(t, app.EnvDevelopment, app.BuildEnv)
}
