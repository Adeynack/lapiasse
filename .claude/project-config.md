# Project Configuration

This file contains project-specific configuration and requirements for Claude Code.

## Go Environment

**IMPORTANT:** All `go` commands in this project MUST be run with the `GOEXPERIMENT=jsonv2` environment variable set.

### Examples

```bash
# Running tests
GOEXPERIMENT=jsonv2 go test ./...

# Running go vet
GOEXPERIMENT=jsonv2 go vet ./...

# Building
GOEXPERIMENT=jsonv2 go build ./...

# Running the application
GOEXPERIMENT=jsonv2 go run ./cmd/app

# Installing tools
GOEXPERIMENT=jsonv2 go install example.com/tool@latest
```

### Why?

This project uses the `jsonv2` experiment which provides an improved JSON API for Go.

### Linting Tools

When running linting tools that invoke `go`, ensure they also respect this environment variable:

```bash
# golangci-lint
GOEXPERIMENT=jsonv2 go tool golangci-lint run

# staticcheck may also need this if it invokes go
GOEXPERIMENT=jsonv2 staticcheck ./...
```

## Makefile

The project includes a Makefile that already sets `GOEXPERIMENT=jsonv2`. Prefer using `make` targets when available.
