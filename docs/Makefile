SHELL := /bin/bash

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
	}
	@# generate service config swagger
	@{ \
		set -e; \
		dirs=$$(find ./gen/config -name 'config-swagger.yaml' -exec dirname {} \;); \
		for dir in $$dirs; do \
			echo generating markdown for $$dir/config-swagger.yaml; \
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
				--allow-template-override \
				--template-dir ./tools/swagger/swagger-templates \
				--config-file ./tools/swagger/markdown-config.yaml; \
			npx prettier --write $$dir/*.md; \
			rm $$dir/gen.yaml; \
		done; \
	}
	@rm merged-plugin-swagger.yaml
	$(MAKE) -B $(POLICY_MARKDOWNS)

POLICY_MARKDOWNS = gen/policies/policy.md

# Note: dependencies are underspecified, run with make -B
$(POLICY_MARKDOWNS): %/policy.md: %/config-swagger.yaml
	@echo generating markdown $@
	@cp ../api/gen/openapiv2/aperture.swagger.yaml $(@D)/
	@yq -i eval 'del(.paths)' $(@D)/aperture.swagger.yaml
	@# 'mixin' is mostly used for --keep-spec-order
	@swagger mixin $< $(@D)/aperture.swagger.yaml --keep-spec-order --format=yaml -o $(@D)/gen.yaml
	@# Fixup .info, which is altered by 'mixin'
	@yq -i eval-all 'select(fileIndex == 0).info = select(fileIndex == 1).info' \
		$(@D)/gen.yaml $(@D)/config-swagger.yaml
	@swagger flatten --with-flatten=remove-unused $(@D)/gen.yaml --format=yaml --output $(@D)/gen.yaml
	@swagger generate markdown --spec $(@D)/gen.yaml --target $(@D) \
		--skip-validation \
		--quiet \
		--with-flatten=remove-unused \
		--tags=$(shell basename $(@D))-configuration \
		--allow-template-override --template-dir tools/swagger/swagger-templates \
		--config-file tools/swagger/markdown-config.yaml
	@npx prettier --write $(@D)/*.md
	@rm $(@D)/aperture.swagger.yaml

# Note: using $(MAKE) -B instead of regular dependency to ensure rerun on tools update
generate-mermaid:
	@echo generating mermaid diagrams
	@find . -name '*.mmd' -exec $(MAKE) -B {}.svg \;

%.mmd.svg: %.mmd ./tools/mermaid/mermaid-theme.json
	@echo MMDC $@
	@for fmt in svg; \
		do npx -p @mermaid-js/mermaid-cli mmdc \
			--quiet --input $< --configFile ./tools/mermaid/mermaid-theme.json --output $<.$$fmt --backgroundColor transparent; \
	done

.PHONY: generate-config-markdown generate-mermaid
