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

func (e Environment) String() string {
	switch e {
	case EnvDevelopment:
		return "development"
	case EnvProduction:
		return "production"
	case EnvTest:
		return "test"
	case EnvUndefined:
		return "undefined"
	default:
		return fmt.Sprintf("unexpected app.Environment: %#v", e)
	}
}

var runEnv Environment = EnvDevelopment

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
