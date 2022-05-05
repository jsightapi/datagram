.PHONY: all
all: fmt lint

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: lint
lint:
	golangci-lint run
