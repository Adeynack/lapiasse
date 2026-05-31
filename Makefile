# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

actual_go_version ?= $(shell go version | awk '{print $$3}' | sed 's/^go//')
export WORKSPACE_ROOT := $(shell git rev-parse --show-toplevel)

# CODE
.PHONY: install-deps
install-deps: install-deps-go

.PHONY: install-deps-go
install-deps-go:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	go mod download -modfile=golangci-lint.mod

.PHONY: upgrade-deps
upgrade-deps: upgrade-deps-go

.PHONY: upgrade-deps-go
upgrade-deps-go:
	@echo "Setting Go version to $(actual_go_version) and upgrading dependencies..."
	@echo "... in main go.mod"
	go mod edit -go=$(actual_go_version)
	go get -u ./...
	go mod tidy
	@echo "... in golangci-lint.mod"
	go mod edit -go=$(actual_go_version) golangci-lint.mod
	go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

.PHONY: golangci-lint
golangci-lint:
	@echo "Running golangci-lint..."
# golangci-lint run with a separate mod file to avoid dependency issues
# See https://golangci-lint.run/docs/welcome/install/#install-from-sources
	go tool -modfile=golangci-lint.mod golangci-lint run

.PHONY: lint-ci
lint-ci: vet
# golangci-lint is executed as a GitHub Action

.PHONY: lint
lint: lint-ci golangci-lint

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf tmp/bin
	go clean -cache -testcache

.PHONY: check
check: clean install-deps gen build lint test

.PHONY: build-only
build-only:
	@echo "Building binaries..."
	mkdir -p tmp/bin
	go build -o tmp/bin/lapiasse ./cmd/lapiasse
	go build -o tmp/bin/import_md_to_lapiasse ./cmd/import_md_to_lapiasse

.PHONY: build
build: gen build-only

.PHONY: build-debug
build-debug: gen
	@echo "Building debug binaries..."
	go build -gcflags="all=-N -l" -o tmp/bin/lapiasse ./cmd/lapiasse

.PHONY: dev
dev:
	@echo "Starting development server with live reload..."
	go tool air -c air_run_watch.toml

.PHONY: gen
gen: gen-api
	@echo "Running go generate..."
	go generate ./...

.PHONY: gen-api
gen-api:
	@echo "Generating Go code from the OpenAPI specification..."
	go tool oapi-codegen --config=pkg/api/lapiasse.oapi-codegen.yaml pkg/api/lapiasse.oas.yaml

# TESTS
.PHONY: test
test:
	@echo "Running tests..."
	go test -json -tags="test" ./... | go tool gotestfmt

# RUN

.PHONE: run
run: build
	@echo "Running lapiasse server..."
	go run ./cmd/lapiasse

.PHONE: run-serve
run-serve: build
	@echo "Running lapiasse server..."
	go run ./cmd/lapiasse serve
