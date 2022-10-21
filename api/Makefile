SHELL := /bin/bash

# use GOMAXPROCS from environment if set, otherwise default to 4
export GOMAXPROCS ?= 4

BUFPRESENT := $(shell command -v buf 2> /dev/null)
GIT_ROOT := $(shell git rev-parse --show-toplevel)
DOCS_OPENAPI := $(GIT_ROOT)/docs/content/assets/openapiv2

buf-update:
	@echo Updating buf dependencies
	@buf mod update

buf-generate:
	@echo Generating code from proto with buf
	@rm -rfd gen/proto gen/openapiv2
	@buf format -w
	@buf lint
	@buf generate
	@#inject go annotations
	@find . -name \*.pb.go -exec protoc-go-inject-tag -input={} \;
	@#generate sdk flowcontrol stubs and copy them over
	@#golang
	@rm -rf ../sdks/aperture-go/gen/proto/flowcontrol
	@cp -R gen/proto/go/aperture/flowcontrol ../sdks/aperture-go/gen/proto/flowcontrol
	@#java
	@rm -rf ../sdks/aperture-java/src/main/java/com/fluxninja/aperture/flowcontrol
	@cp -R gen/proto/java/com/fluxninja/aperture/flowcontrol ../sdks/aperture-java/src/main/java/com/fluxninja/aperture/
	@rm -rf gen/proto/java
	@#generate api docs
	@rm -rfd $(DOCS_OPENAPI)
	@mv gen/openapiv2 $(DOCS_OPENAPI)
	@{ \
		set -e; \
		for process in aperture-agent aperture-controller; do \
			yq eval "del(.paths | .[] | select(.*.tags | contains([\"$$process\"]) | not))" $(DOCS_OPENAPI)/aperture.swagger.yaml > $(DOCS_OPENAPI)/$$process.swagger.yaml; \
			yq eval -i "del(.tags)" $(DOCS_OPENAPI)/$$process.swagger.yaml; \
			yq eval -i ".host = \"$$process\"" $(DOCS_OPENAPI)/$$process.swagger.yaml; \
			yq eval -i '.schemes = ["https"]' $(DOCS_OPENAPI)/$$process.swagger.yaml; \
			swagger flatten --with-flatten=remove-unused $(DOCS_OPENAPI)/$$process.swagger.yaml --format=yaml --output=$(DOCS_OPENAPI)/$$process.swagger.yaml; \
		done; \
	}
	@git add $(DOCS_OPENAPI)
	@git add ./gen/*

generate: buf-update buf-generate

.PHONY: generate buf-update buf-generate
