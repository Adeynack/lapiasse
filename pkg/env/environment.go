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

var RunEnv Environment = EnvDevelopment

func (e Environment) String() string {
	switch RunEnv {
	case EnvDevelopment:
		return "development"
	case EnvProduction:
		return "production"
	case EnvTest:
		return "test"
	case EnvUndefined:
		return "undefined"
	default:
		return fmt.Sprintf("unexpected app.Environment: %#v", RunEnv)
	}
}

func init() {
	RunEnv = determineRuntimeEnvironment()
}

func determineRuntimeEnvironment() Environment {
	if testing.Testing() {
		return EnvTest
	}

	if BuildEnv != EnvUndefined {
		return BuildEnv
	}

	return EnvDevelopment
}
