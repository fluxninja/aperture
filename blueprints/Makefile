generate-blueprints: generate-libsonnet
	@echo Generating blueprints readme and examples
	@scripts/generate_all.sh

generate-libsonnet:
	@echo Generating library
	@{ \
		python lib/scripts/jsonnet-lib-gen.py --output-dir lib/1.0/ ../docs/gen/policies/gen.yaml; \
		tk fmt lib/1.0/; \
		git add lib/1.0/*; \
	}


.PHONY: generate-blueprints generate-libsonnet
