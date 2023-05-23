SHELL := /bin/bash

# use GOMAXPROCS from environment if set, otherwise default to 4
export GOMAXPROCS ?= 4

generate-config-markdown:
	@echo generate config markdown
	@./tools/swagger/generate_config_markdown.sh

generate-jsonnet:
	@echo generate jsonnet
	@./tools/jsonnet/generate_jsonnet.sh
	@./tools/aperturectl/validate_all.sh

# Note: using $(MAKE) -B instead of regular dependency to ensure rerun on tools update
generate-mermaid:
	@echo generating mermaid diagrams
	@./tools/mermaid/generate_all.sh

generate-aperturectl-docs:
	@echo generating aperturectl
	@go run tools/aperturectl/generate-docs/generate-docs.go
	@prettier --prose-wrap="preserve" ./content/reference/aperturectl/ --write
	@git add ./content/reference/aperturectl/

generate-otel-docs:
	@echo generating metrics docs
	@./tools/otel_collectors/generate_otel_docs.sh
	@git add ./content/integrations/metrics

.PHONY: generate-config-markdown generate-mermaid generate-jsonnet generate-aperturectl-docs generate-otel-docs
