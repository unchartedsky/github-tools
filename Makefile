NAME	:= github-tools

BIN		:= $(NAME)

SHELL := /bin/bash


# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all: test

.PHONY: build
build: fmt
	@go build -o bin/$(BIN) .

.PHONY: test
test: build
	@go test .

.PHONY: fmt
fmt:
	@ if ! which goimports > /dev/null; then \
		go get -u -v golang.org/x/tools/cmd/goimports; \
	fi

	go mod tidy
	goimports -l -w $(SRC)
	gofmt -l -w -s $(SRC)

.PHONY: run
run: build
	bin/$(BIN)

.PHONY: snapshot
snapshot:
	goreleaser --snapshot --rm-dist
