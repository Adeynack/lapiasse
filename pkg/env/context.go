package env

import (
	"context"

	"adeynack.net/lapiasse/pkg/platform/ctxval"
)

func AutoRegisterEnvironments(ctx context.Context) context.Context {
	ctx = ctxval.RegisterNamed(ctx, "build", buildEnv)
	ctx = ctxval.RegisterNamed(ctx, "run", runEnv)

	return ctx
}

func GetRunEnv(ctx context.Context) Environment {
	return ctxval.MustResolveNamed[Environment](ctx, "run")
}

func GetBuildEnv(ctx context.Context) Environment {
	return ctxval.MustResolveNamed[Environment](ctx, "build")
}
