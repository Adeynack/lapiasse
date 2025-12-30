# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

actual_go_version ?= $(shell go version | awk '{print $$3}' | sed 's/^go//')
export WORKSPACE_ROOT := $(shell git rev-parse --show-toplevel)

# CODE
.PHONY: install-deps
install-deps: install-deps-go install-deps-npm

.PHONY: install-deps-go
install-deps-go:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	go mod download -modfile=golangci-lint.mod

.PHONY: install-deps-npm
install-deps-npm:
	@echo "Installing NPM dependencies..."
	npm install

.PHONY: upgrade-deps
upgrade-deps: upgrade-deps-go upgrade-deps-npm

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

.PHONY: upgrade-deps-npm
upgrade-deps-npm:
	@echo "Upgrading NPM dependencies..."
	npm upgrade

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

.PHONY: lint
lint: vet golangci-lint

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf tmp/bin
	go clean -cache -testcache

.PHONY: check
check: clean install-deps build lint test

.PHONY: build
build: gen
	@echo "Building binaries..."
	mkdir -p tmp/bin
	go build -o tmp/bin/lapiasse ./cmd/lapiasse
	go build -o tmp/bin/import_md_to_lapiasse ./cmd/import_md_to_lapiasse

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
	@echo "Generating API code..."
	@echo "... TypeSpec -to- OpenAPI Specification ..."
	npm exec -- tsp compile --config ./pkg/api/lapiasse.tsp.yaml --output-dir ./pkg/api ./pkg/api/lapiasse.tsp
	@echo "... OpenAPI Specification to Go code ..."
	go tool oapi-codegen --config=pkg/api/lapiasse.oapi-codegen.yaml pkg/api/lapiasse.oas.yaml

# TESTS
.PHONY: test
test:
	@echo "Running tests..."
	go test -json -tags="test" ./... | go tool gotestfmt
