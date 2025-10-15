GOEXPERIMENT ?= jsonv2

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
	go test -tags="test" ./... | go tool gotestfmt

# DATABASE
.PHONY: db_new_migration
db_new_migration:
	go tool migrate

foo:
	echo ${GOEXPERIMENT}