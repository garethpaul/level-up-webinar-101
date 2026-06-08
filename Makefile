.PHONY: check fmt

check:
	test -z "$$(gofmt -l main.go main_test.go)"
	go test ./...

fmt:
	gofmt -w main.go main_test.go
