aperture_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

go-mod-tidy:
	@echo Download go.mod dependencies
	@GOPRIVATE=github.com/FluxNinja,github.com/aperture-control
	@go mod tidy

install-go-tools: go-mod-tidy
	@echo Installing tools from tools.go
	@./scripts/install_go_tools.sh

go-generate-swagger:
	@echo Generating swagger code
	@echo Generating swagger specs from go code
	@./scripts/go_generate_swagger.sh

generate-config-markdown: go-generate-swagger
	@cd ./docs && $(MAKE) generate-config-markdown

generate-mermaid:
	@cd ./docs && $(MAKE) generate-mermaid

coverage_profile:
	go test -v -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage_in_browser: profile.coverprofile
	go tool cover -html profile.coverprofile

.PHONY: download install-go-tools coverage_profile show_coverage_in_browser
