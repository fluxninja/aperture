BUFPRESENT := $(shell command -v buf 2> /dev/null)

buf-mod-update:
	@echo Updating buf dependencies
	@buf mod update

generate-buf:
	@echo Generating code from proto with buf
	@rm -rfd gen/proto gen/openapiv2
	@buf format -w
	@buf lint
	@buf generate
	@find . -name \*.pb.go -exec protoc-go-inject-tag -input={} \;

generate-ci: generate-buf
	@#echo generating HTTP clients with go-swagger
	@#rm -rfd gen/http
	@#mkdir -p gen/http
	@#swagger generate client --spec=./gen/openapiv2/aperture.swagger.yaml --template=stratoscale --target=./gen/http
	@#../../../../scripts/go_generate.sh
	git add ./gen/*

generate: buf-mod-update generate-ci
	@echo Generating code from proto with buf

.PHONY: buf-mod-update generate-buf generate generate-ci
