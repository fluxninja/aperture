SHELL := /bin/bash

# use GOMAXPROCS from environment if set, otherwise default to 4
export GOMAXPROCS ?= 4

BUFPRESENT := $(shell command -v buf 2> /dev/null)
GIT_ROOT := $(shell git rev-parse --show-toplevel)

buf-update:
	@echo Updating buf dependencies
	@buf mod update

buf-generate:
	@echo Generating code from proto with buf
	@rm -rfd gen/proto gen/openapiv2
	@buf format -w
	@buf lint
	@buf generate
	@find . -name \*.pb.go -exec protoc-go-inject-tag -input={} \;
	@rm -rfd $(GIT_ROOT)/docs/content/assets/openapiv2
	@mv gen/openapiv2 $(GIT_ROOT)/docs/content/assets/openapiv2
	@git add $(GIT_ROOT)/docs/content/assets/openapiv2
	@git add ./gen/*

generate: buf-update buf-generate

.PHONY: generate buf-update buf-generate
