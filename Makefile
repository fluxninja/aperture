aperture_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

generate-api:
	@echo Generating API
	@cd api && $(MAKE) generate

go-generate:
	@echo Generating go code
	@./scripts/go_generate.sh

go-mod-tidy:
	@echo Download go.mod dependencies
	@GOPRIVATE=github.com/fluxninja,github.com/aperture-control
	@go mod tidy

go-test:
	@echo Running go tests
	gotestsum --format=testname
	@cd ./pkg/watchdog && go test -v ./...

go-lint:
	@echo Linting go code
	@./scripts/golangci_lint.sh --fix

go-build:
	@echo Building go code
	@./scripts/go_build.sh

go-build-plugins:
	@echo Building go plugins
	@./scripts/go_build_plugins.sh

install-go-tools:
	@echo Installing tools from tools.go
	@./scripts/install_go_tools.sh

go-generate-swagger:
	@echo Generating swagger code
	@echo Generating swagger specs from go code
	@./scripts/go_generate_swagger.sh

generate-docs: generate-config-markdown generate-mermaid
	@echo Generating docs

generate-config-markdown: go-generate-swagger
	@cd ./docs && $(MAKE) generate-config-markdown

generate-mermaid:
	@cd ./docs && $(MAKE) generate-mermaid

coverage_profile:
	gotestsum --format=testname -- -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage_in_browser: profile.coverprofile
	go tool cover -html profile.coverprofile

all: install-go-tools generate-api go-generate go-mod-tidy go-lint go-build go-build-plugins go-test generate-docs
	@echo "Done"

.PHONY: install-go-tools generate-api go-generate go-generate-swagger go-mod-tidy generate-config-markdown generate-mermaid generate-docs go-test go-lint go-build go-build-plugins coverage_profile show_coverage_in_browser
