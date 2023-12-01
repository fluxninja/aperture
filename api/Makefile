SHELL := /bin/bash

# use GOMAXPROCS from environment if set, otherwise default to 4
export GOMAXPROCS ?= 4

BUFPRESENT := $(shell command -v buf 2> /dev/null)
GIT_ROOT := $(shell git rev-parse --show-toplevel)
DOCS_OPENAPI := $(GIT_ROOT)/docs/content/assets/openapiv2
DOCS_ROOT := $(GIT_ROOT)/docs

buf-update:
	@echo Updating buf dependencies
	@buf mod update

buf-generate:
	@echo Generating code from proto with buf
	@rm -rfd gen/proto gen/openapiv2
	@buf format -w
	@buf lint
	@buf generate --path aperture --template buf.gen.aperture-langs2.yaml
	@buf generate --path aperture --template buf.gen.aperture-langs.yaml
	@buf generate --path aperture --template buf.gen.aperture-grpc-ecosystem.yaml

	@#inject go annotations
	@find . -name \*.pb.go -exec protoc-go-inject-tag -input={} \;

	@#generate sdk flowcontrol stubs and copy them over

	@#golang
	@rm -rf ../sdks/aperture-go/gen/proto/flowcontrol/*
	@cp -R gen/proto/go/aperture/flowcontrol/check ../sdks/aperture-go/gen/proto/flowcontrol/
	@cp -R gen/proto/go/aperture/flowcontrol/checkhttp ../sdks/aperture-go/gen/proto/flowcontrol/

	@#csharp
	@rm -rf ../sdks/aperture-csharp/Generated/Aperture/Flowcontrol/
	@cp -R gen/proto/csharp/Aperture/Flowcontrol/Check ../sdks/aperture-csharp/Generated/Aperture/Flowcontrol/
	@rm -rf gen/proto/csharp

	@#java
	@rm -rf ../sdks/aperture-java/lib/core/src/main/java/com/fluxninja/generated/aperture/flowcontrol/*
	@cp -R gen/proto/java/com/fluxninja/generated/aperture/flowcontrol/check ../sdks/aperture-java/lib/core/src/main/java/com/fluxninja/generated/aperture/flowcontrol/
	@cp -R gen/proto/java/com/fluxninja/generated/aperture/flowcontrol/checkhttp ../sdks/aperture-java/lib/core/src/main/java/com/fluxninja/generated/aperture/flowcontrol/
	@rm -rf gen/proto/java

	@#javascript
	@rm -rf ../sdks/aperture-js/src/generated/aperture/flowcontrol/*
	@cp -R aperture/flowcontrol/check ../sdks/aperture-js/proto/flowcontrol/
	@{ \
		set -e; \
		pushd ../sdks/aperture-js; \
		npm install; \
		node ./node_modules/@grpc/proto-loader/build/bin/proto-loader-gen-types.js --defaults=true --longs=String --outDir=sdk/gen --grpcLib=@grpc/grpc-js proto/flowcontrol/check/v1/check.proto; \
		git add ./sdk/gen; \
		popd; \
	}

	@#python
	@rm -rf ../sdks/aperture-py/aperture_sdk/_gen/aperture/flowcontrol/*
	@cp -R gen/proto/python/aperture/flowcontrol/check ../sdks/aperture-py/aperture_sdk/_gen/aperture/flowcontrol/
	@# Rewrite the generated structure to add __init__.py files and fix imports
	@{ \
	  set -euo pipefail; \
		readonly python_gen="../sdks/aperture-py/aperture_sdk/_gen"; \
		SED="sed"; \
		FIND="find"; \
		if [[ "$$OSTYPE" == "darwin"* ]]; then \
			SED="gsed"; \
			FIND="gfind"; \
		fi; \
		echo "Fixing up Python gRPC imports"; \
		generated_py_files_str="$$($$FIND "$$python_gen" -type f -name '*.py')"; \
		while IFS= read -r generated_py_file; do \
			generated_py_files+=("$$generated_py_file"); \
		done <<<"$$generated_py_files_str"; \
		$$SED -i "s/^from aperture\..* import \([^ ]*_pb2\) as \([^ ]*\)$$/from . import \1 as \2/" "$${generated_py_files[@]}"; \
		generated_pyi_files_str="$$($$FIND "$$python_gen" -type f -name '*.pyi')"; \
		while IFS= read -r generated_pyi_file; do \
			generated_pyi_files+=("$$generated_pyi_file"); \
		done <<<"$$generated_pyi_files_str"; \
		$$SED -i "s/^import aperture\.[^ ]*\.\([^ ]*_pb2\)$$/from . import \1/" "$${generated_pyi_files[@]}"; \
		$$SED -i "s/aperture\.[^ ]*\.\([^ ]*_pb2\)/\1/" "$${generated_pyi_files[@]}"; \
		$$SED -i "s/^import grpc$$/import grpc\nimport grpc.aio\nimport typing/" "$${generated_pyi_files[@]}"; \
		$$SED -i "s/: grpc\.Server/: typing.Union[grpc.Server, grpc.aio.Server]/" "$${generated_pyi_files[@]}"; \
		echo "Generating __init__.py files"; \
		$$FIND "$$python_gen" -type d -exec touch "{}/__init__.py" \;; \
	}
	@rm -rf gen/proto/python

generate: buf-update buf-generate

.PHONY: generate buf-update buf-generate
