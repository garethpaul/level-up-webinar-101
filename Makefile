.PHONY: check fmt

check:
	test -z "$$(gofmt -l *.go)"
	go mod verify
	go test ./...
	go build ./...

fmt:
	gofmt -w *.go
