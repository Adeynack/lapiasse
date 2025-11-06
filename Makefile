# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

# CODE
.PHONY: upgrade-deps
upgrade-deps:
	go get -u ./...
	go mod tidy
	go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

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

.PHONY: run-watch
run-watch:
	go tool air -c air_run_watch.toml

.PHONY: gen
gen:
	go generate ./...

# TESTS
.PHONY: test
test:
	go test -json -tags="test" ./... | go tool gotestfmt
