gen-lib:
	@echo Generating library
	@rm -rfd 1.0/_gen/*
	@python scripts/jsonnet-lib-gen.py --output-dir 1.0/ ../docs/gen/policies/gen.yaml
	@git add 1.0/_gen/*

.PHONY: gen-lib
