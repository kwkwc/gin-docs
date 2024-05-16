SHELL=/bin/bash

.PHONY: install
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install golang.org/x/tools/cmd/goimports@v0.20.0
	go mod tidy

.PHONY: format
format:
	gofmt -l -s -w .

.PHONY: format-check
format-check:
	diff -u <(echo -n) <(gofmt -d .)

.PHONY: isort
isort:
	find . -type f -name '*.go' -not -name '*.pb.go' | xargs goimports -l -w -local github.com/kwkwc/gin-docs

.PHONY: isort-check
isort-check:
	diff -u <(echo -n) <(find . -type f -name '*.go' -not -name '*.pb.go' | xargs goimports -d -local github.com/kwkwc/gin-docs)

.PHONY: lint
lint:
	golangci-lint run --timeout=5m

.PHONY: test
test:
	go test \
		-timeout 120s \
		-covermode=set \
		-coverprofile=coverage.out \
		. \
		-v
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: check-all
check-all: format-check isort-check lint test
