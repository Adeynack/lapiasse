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

.PHONY: gen
gen:
	go generate ./...

# TESTS
.PHONY: test
test:
	go test -json -tags="test" ./... | go tool gotestfmt
