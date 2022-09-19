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
				--allow-template-override \
				--template-dir ./tools/swagger/swagger-templates \
				--config-file ./tools/swagger/markdown-config.yaml \
				--output $$basename.md; \
			rm $$dir/gen.yaml; \
			cat $$dir/metadata $$dir/$$basename.md > $$dir/$$basename.md.tmp; \
			mv $$dir/$$basename.md.tmp $$dir/$$basename.md; \
			npx prettier --write $$dir/$$basename.md; \
			mv $$dir/$$basename.md ./content/references/configuration; \
		done; \
		rm merged-plugin-swagger.yaml; \
	}
	@# policies markdown
	@echo generating policies markdown
	@cp ../api/gen/openapiv2/aperture.swagger.yaml gen/policies/
	@yq -i eval 'del(.paths)' gen/policies/aperture.swagger.yaml
	@# 'mixin' is mostly used for --keep-spec-order
	@swagger mixin gen/policies/config-swagger.yaml gen/policies/aperture.swagger.yaml --keep-spec-order --format=yaml -o gen/policies/gen.yaml
	@# Fixup .info, which is altered by 'mixin'
	@yq -i eval-all 'select(fileIndex == 0).info = select(fileIndex == 1).info' \
		gen/policies/gen.yaml gen/policies/config-swagger.yaml
	@swagger flatten --with-flatten=remove-unused gen/policies/gen.yaml --format=yaml --output gen/policies/gen.yaml
	@swagger generate markdown --spec gen/policies/gen.yaml --target gen/policies \
		--skip-validation \
		--quiet \
		--with-flatten=remove-unused \
		--tags=$(name)-configuration \
		--allow-template-override --template-dir tools/swagger/swagger-templates \
		--config-file tools/swagger/markdown-config.yaml
	@rm gen/policies/aperture.swagger.yaml
	@# append gen/policies/metadata on top of gen/policies/policies.md
	@cat gen/policies/metadata gen/policies/policies.md > gen/policies/policies.md.tmp
	@mv gen/policies/policies.md.tmp gen/policies/policies.md
	@npx prettier --write gen/policies/policies.md
	@mv gen/policies/policies.md content/references/configuration

generate-mmd:
	@echo generating mermaid files from markdown files
	@./tools/mermaid/md_to_mmd.sh

# Note: using $(MAKE) -B instead of regular dependency to ensure rerun on tools update
generate-mermaid: generate-mmd
	@echo generating mermaid diagrams
	@find . -name '*.mmd' -exec $(MAKE) -B {}.svg \;

%.mmd.svg: %.mmd ./tools/mermaid/mermaid-theme.json
	@echo MMDC $@
	@for fmt in svg; \
		do npx -p @mermaid-js/mermaid-cli mmdc \
			--quiet --input $< --configFile ./tools/mermaid/mermaid-theme.json --output $<.$$fmt --backgroundColor transparent \
			&& git add $<.$$fmt; \
	done

.PHONY: generate-config-markdown generate-mermaid generate-mmd
