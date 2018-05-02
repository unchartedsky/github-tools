.DEFAULT_GOAL = all

GO := go
VERSION  := $(shell git rev-list --count HEAD).$(shell git rev-parse --short HEAD)

NAME     := github-tools
PACKAGE  := github.com/corpix/$(NAME)
PACKAGES := $(shell go list ./... | grep -v /vendor/)

BIN := $(NAME)

.PHONY: all
all:: dependencies
all:: build

dependencies::
	dep ensure

build: dependencies
	$(GO) build -a -installsuffix cgo -o bin/$(BIN) .

test: build
	go test -v $(PACKAGES)

.PHONY: bench
bench::
	go test  -race -coverprofile=coverage.txt -covermode=atomic -bench=. -v $(PACKAGES)

.PHONY: lint
lint::
	go vet -v $(PACKAGES)

.PHONY: check
check:: lint test

.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.txt

cleanall: clean
	rm -rf vendor/
	# git clean -xddff
