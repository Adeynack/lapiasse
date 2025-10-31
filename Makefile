# Ensuring all tasks are using the new JSONv2 Go library.
export GOEXPERIMENT := $(GOEXPERIMENT),jsonv2

# CODE
.PHONY: lint
lint:
	go vet ./...
	go tool staticcheck ./...
	go tool golangci-lint run

.PHONY: clean
clean:
	go clean -cache -testcache

.PHONY: check
check: clean build lint test

.PHONY: build
build:
	go build -o tmp/bin/lapiasse ./main.go

# TESTS
.PHONY: test
test:
	go test -json -tags="test" ./... | go tool gotestfmt
