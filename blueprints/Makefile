generate-blueprints: generate-libsonnet
	@echo Generating blueprints readme and examples
	@scripts/generate_all.sh

generate-libsonnet:
	@echo Generating library
	@{ \
		rm -rf libsonnet/1.0/_gen/* || true; \
		rm -rf libsonnet/1.0/*.libsonnet || true; \
		python libsonnet/scripts/jsonnet-lib-gen.py --output-dir libsonnet/1.0/ ../docs/gen/policies/gen.yaml; \
		tk fmt libsonnet/1.0/; \
		git add libsonnet/1.0/*; \
	}


.PHONY: generate-blueprints generate-libsonnet
