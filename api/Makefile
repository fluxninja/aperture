BUFPRESENT := $(shell command -v buf 2> /dev/null)

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

.PHONY: buf-update buf-generate
