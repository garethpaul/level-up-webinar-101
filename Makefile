.PHONY: build check fmt lint test

check: lint test build
	./scripts/check-baseline.sh

lint:
	test -z "$$(gofmt -l *.go)"
	go vet ./...

test:
	go mod verify
	go test ./...

build:
	go build ./...

fmt:
	gofmt -w *.go
