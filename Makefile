# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

actual_go_version ?= $(shell go version | awk '{print $$3}' | sed 's/^go//')
export WORKSPACE_ROOT := $(shell git rev-parse --show-toplevel)

# CODE
.PHONY: upgrade-deps
upgrade-deps:
	@echo "Setting Go version to $(actual_go_version) and upgrading dependencies..."
	@echo "... in main go.mod"
	go mod edit -go=$(actual_go_version)
	go get -u ./...
	go mod tidy
	@echo "... in golangci-lint.mod"
	go mod edit -go=$(actual_go_version) golangci-lint.mod
	go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "... in GitHub Actions"
	sed -i '' "s/go-version: \".*\"/go-version: \"1.25.4\"/" .github/workflows/pr-check.yml

.PHONY: lint
lint:
	go vet ./...
# golangci-lint run with a separate mod file to avoid dependency issues
# See https://golangci-lint.run/docs/welcome/install/#install-from-sources
	go tool -modfile=golangci-lint.mod golangci-lint run

.PHONY: clean
clean:
	go clean -cache -testcache

.PHONY: check
check: clean build lint test

.PHONY: build
build: gen
	go build -o tmp/bin/lapiasse ./cmd/lapiasse
	go build -o tmp/bin/import_md_to_lapiasse ./cmd/import_md_to_lapiasse

.PHONY: build-debug
build-debug:
	go build -gcflags="all=-N -l" -o tmp/bin/lapiasse ./cmd/lapiasse

.PHONY: dev
dev:
	go tool air -c air_run_watch.toml

.PHONY: gen
gen:
	go generate ./...

# TESTS
.PHONY: test
test:
	go test -json -tags="test" ./... | go tool gotestfmt
