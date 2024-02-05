NAME := atlas-hcl-gen-go
BUILD_DIR := ./build
GO ?= go

.PHONY: build
build: VERSION := $(shell git describe --tags --always --dirty)
build: REVISION := $(shell git rev-parse --short HEAD)
build: TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
build:
	go build -ldflags "-X main.BuildVersion=$(VERSION) -X main.BuildRevision=$(REVISION) -X main.BuildTimestamp=$(TIMESTAMP)" -o $(NAME)

.PHONY: test
test: PKG ?= ./...
test: FLAGS ?= -race
test:
	$(GO) test $(FLAGS) $(PKG)
