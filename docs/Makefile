SHELL := /bin/bash

# use GOMAXPROCS from environment if set, otherwise default to 4
export GOMAXPROCS ?= 4

generate-config-markdown:
	@# merge plugin swaggers
	@{ \
		set -e; \
		cp ./gen/config/aperture-plugin-fluxninja/plugin-swagger.yaml merged-plugin-swagger.yaml; \
		dirs=$$(find . -name 'plugin-swagger.yaml' -exec dirname {} \;); \
		for dir in $$dirs; do \
			echo "Merging $$dir/plugin-swagger.yaml"; \
			yq eval-all --inplace "select(fileIndex==0).definitions *= select(fileIndex==1).definitions | select(fileIndex==0)" merged-plugin-swagger.yaml $$dir/plugin-swagger.yaml; \
			yq eval-all --inplace "select(fileIndex==0).paths *= select(fileIndex==1).paths | select(fileIndex==0)" merged-plugin-swagger.yaml $$dir/plugin-swagger.yaml; \
		done; \
		dirs=$$(find ./gen/config -name 'config-swagger.yaml' -exec dirname {} \;); \
		for dir in $$dirs; do \
			echo generating markdown for $$dir/config-swagger.yaml; \
			basename=$$(basename $$dir); \
			cp $$dir/config-swagger.yaml $$dir/gen.yaml; \
			yq eval-all --inplace "select(fileIndex==0).definitions *= select(fileIndex==1).definitions | select(fileIndex==0)" $$dir/gen.yaml merged-plugin-swagger.yaml; \
			yq eval-all --inplace "select(fileIndex==0).paths *= select(fileIndex==1).paths | select(fileIndex==0)" $$dir/gen.yaml merged-plugin-swagger.yaml; \
			swagger flatten \
				--with-flatten=remove-unused $$dir/gen.yaml \
				--format=yaml --output $$dir/gen.yaml; \
			swagger generate markdown \
				--spec $$dir/gen.yaml \
				--target $$dir \
				--skip-validation \
				--quiet \
				--with-flatten=remove-unused \
				--tags=common-configuration \
				--tags=plugin-configuration \
				--tags=agent-configuration \
				--tags=controller-configuration \
				--allow-template-override \
				--template-dir ./tools/swagger/swagger-templates \
				--config-file ./tools/swagger/markdown-config.yaml \
				--output $$basename.md; \
			rm $$dir/gen.yaml; \
			cat $$dir/metadata $$dir/$$basename.md > $$dir/$$basename.md.tmp; \
			mv $$dir/$$basename.md.tmp $$dir/$$basename.md; \
			npx prettier --prose-wrap="preserve" --write $$dir/$$basename.md; \
			mv $$dir/$$basename.md ./content/reference/configuration; \
		done; \
		rm merged-plugin-swagger.yaml; \
	}
	@# policy markdown
	@echo generating policy markdown
	@cp content/assets/openapiv2/aperture.swagger.yaml gen/policy/
	@yq -i eval 'del(.paths)' gen/policy/aperture.swagger.yaml
	@yq -i eval 'del(.tags)' gen/policy/aperture.swagger.yaml
	@# 'mixin' is mostly used for --keep-spec-order
	@swagger mixin gen/policy/config-swagger.yaml gen/policy/aperture.swagger.yaml --keep-spec-order --format=yaml -o gen/policy/policy.yaml
	@# Fixup .info, which is altered by 'mixin'
	@yq -i eval-all 'select(fileIndex == 0).info = select(fileIndex == 1).info' \
		gen/policy/policy.yaml gen/policy/config-swagger.yaml
	@swagger flatten --with-flatten=remove-unused gen/policy/policy.yaml --format=yaml --output gen/policy/policy.yaml
	@swagger generate markdown --spec gen/policy/policy.yaml --target gen/policy \
		--skip-validation \
		--quiet \
		--with-flatten=remove-unused \
		--tags=policy-configuration \
		--allow-template-override --template-dir tools/swagger/swagger-templates \
		--config-file tools/swagger/markdown-config.yaml
	@rm gen/policy/aperture.swagger.yaml
	@# append gen/policy/metadata on top of gen/policy/policy.md
	@cat gen/policy/metadata gen/policy/policy.md > gen/policy/policy.md.tmp
	@mv gen/policy/policy.md.tmp gen/policy/policy.md
	@npx prettier --prose-wrap="preserve" --write gen/policy/policy.md
	@mv gen/policy/policy.md content/reference/policies/spec.md

generate-jsonnet:
	@echo generate jsonnet
	@./tools/jsonnet/generate_jsonnet.sh

# Note: using $(MAKE) -B instead of regular dependency to ensure rerun on tools update
generate-mermaid:
	@echo generating mermaid diagrams
	@./tools/mermaid/generate_mermaid.sh

generate-aperturectl-docs:
	@echo generating aperturectl
	@go run tools/aperturectl/generate-docs.go
	@npx prettier --prose-wrap="preserve" ./content/reference/aperturectl/ --write
	@git add ./content/reference/aperturectl/

.PHONY: generate-config-markdown generate-mermaid generate-jsonnet generate-aperturectl
