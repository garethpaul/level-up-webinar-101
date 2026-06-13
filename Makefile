ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.PHONY: build check fmt lint test

check: lint test build
	cd "$(ROOT)" && ./scripts/check-baseline.sh

lint:
	cd "$(ROOT)" && test -z "$$(gofmt -l *.go)"
	cd "$(ROOT)" && go vet ./...

test:
	cd "$(ROOT)" && go mod verify
	cd "$(ROOT)" && go test ./...

build:
	cd "$(ROOT)" && go build ./...

fmt:
	cd "$(ROOT)" && gofmt -w *.go
