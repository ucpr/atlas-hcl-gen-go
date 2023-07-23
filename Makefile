.PHONY: build
build: VERSION := $(shell git describe --tags --always --dirty)
build: REVISION := $(shell git rev-parse --short HEAD)
build: TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
build: BINARY_NAME := "atlas-hcl-gen-go"
build:
	go build -ldflags "-X main.BuildVersion=$(VERSION) -X main.BuildRevision=$(REVISION) -X main.BuildTimestamp=$(TIMESTAMP)" -o $(BINARY_NAME)
