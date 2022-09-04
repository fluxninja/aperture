SHELL := /bin/bash

generate-helm-readme:
	@{ \
		dirs=$$(find ./ -name 'values.yaml' -exec dirname {} \;); \
		for dir in $${dirs}; do \
			chart=$$(basename $${dir}); \
			echo "Generating README for $${chart}"; \
			npx readme-generator-for-helm -v $$dir/values.yaml -r $$dir/README.md; \
		done; \
		\
	}

.PHONY: generate-helm-readme
