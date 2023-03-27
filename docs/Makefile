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

generate-policy-graphs:
	@echo generating policy graphs
	@go run tools/aperturectl/generate-policy-graphs/generate-policy-graphs.go

# Note: using $(MAKE) -B instead of regular dependency to ensure rerun on tools update
generate-mermaid:
	@echo generating mermaid diagrams
	@./tools/mermaid/generate_mermaid.sh

generate-aperturectl-docs:
	@echo generating aperturectl
	@go run tools/aperturectl/generate-docs/generate-docs.go
	@npx prettier --prose-wrap="preserve" ./content/reference/aperturectl/ --write
	@git add ./content/reference/aperturectl/

.PHONY: generate-config-markdown generate-mermaid generate-jsonnet generate-aperturectl-docs generate-policy-graphs
