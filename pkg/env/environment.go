package env

import (
	"fmt"
	"testing"
)

type Environment uint

const (
	EnvUndefined Environment = iota
	EnvDevelopment
	EnvTest
	EnvProduction
)

var runEnv Environment = EnvDevelopment

func (e Environment) String() string {
	switch runEnv {
	case EnvDevelopment:
		return "development"
	case EnvProduction:
		return "production"
	case EnvTest:
		return "test"
	case EnvUndefined:
		return "undefined"
	default:
		return fmt.Sprintf("unexpected app.Environment: %#v", runEnv)
	}
}

func init() {
	runEnv = determineRuntimeEnvironment()
}

func determineRuntimeEnvironment() Environment {
	if testing.Testing() {
		return EnvTest
	}

	if buildEnv != EnvUndefined {
		return buildEnv
	}

	return EnvDevelopment
}
