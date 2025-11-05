# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

# CODE
.PHONY: lint
lint:
	go vet ./...
	go tool golangci-lint run

.PHONY: clean
clean:
	go clean -cache -testcache

.PHONY: check
check: clean build lint test

.PHONY: build
build: gen
	go build -o tmp/bin/lapiasse ./cmd/lapiasse

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
