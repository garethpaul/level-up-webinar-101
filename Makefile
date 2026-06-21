ifneq ($(origin MAKEFILE_LIST),file)
$(error MAKEFILE_LIST must not be overridden)
endif
override ROOT := $(shell path='$(subst ','"'"',$(MAKEFILE_LIST))'; path=$${path\# }; dirname -- "$$path")

.PHONY: build check fmt lint test vuln

check: lint test build vuln
	cd "$(ROOT)" && ./scripts/check-baseline.sh

lint:
	cd "$(ROOT)" && test -z "$$(gofmt -l *.go)"
	cd "$(ROOT)" && go vet ./...

test:
	cd "$(ROOT)" && go mod verify
	cd "$(ROOT)" && go test ./...

build:
	cd "$(ROOT)" && go build ./...

vuln:
	cd "$(ROOT)" && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...

fmt:
	cd "$(ROOT)" && gofmt -w *.go
